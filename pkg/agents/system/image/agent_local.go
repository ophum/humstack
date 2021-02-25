package image

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ceph/go-ceph/rbd"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
)

// 同じノードにあるBSを元にイメージを作成する
func (a *ImageAgent) syncLocalImageEntity(imageEntity *system.ImageEntity, bs *system.BlockStorage) error {

	if imageEntity.DeleteState == meta.DeleteStateDelete {
		if imageEntity.Status.State != system.ImageEntityStateAvailable {
			return nil
		}

		imageEntity.Status.State = system.ImageEntityStateDeleting
		if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
			return err
		}

		path := filepath.Join(a.localImageDirectory, imageEntity.Group, imageEntity.ID)
		if fileIsExists(path) {
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		if err := a.client.SystemV0().ImageEntity().Delete(imageEntity.Group, imageEntity.ID); err != nil {
			return err
		}

		return nil
	}

	imageEntity.Status.State = system.ImageEntityStatePending
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	//if bs.Status.State != system.BlockStorageStateActive {
	//	return nil
	//}

	imageEntity.Status.State = system.ImageEntityStateCopying
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}
	bs.Status.State = system.BlockStorageStateCopying
	if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
		return err
	}

	var src io.Reader
	size := int64(0)
	if t, ok := bs.Annotations["blockstoragev0/type"]; ok && t == "Ceph" {
		imageName, ok := bs.Annotations["ceph-image-name"]
		if !ok {
			return fmt.Errorf("ceph-image-name not found")
		}

		conn, err := a.newCephConn()
		if err != nil {
			return errors.Wrap(err, "new ceph conn")
		}
		defer conn.Shutdown()

		ioctx, err := conn.OpenIOContext(a.config.CephBackend.PoolName)
		if err != nil {
			return errors.Wrap(err, "open io context")
		}
		defer ioctx.Destroy()

		if image, err := rbd.OpenImageReadOnly(ioctx, imageName, ""); err != nil {
			return errors.Wrapf(err, "open rbd image `%s`", imageName)
		} else {
			defer image.Close()
			src = image

			limitSize, err := image.GetSize()
			if err != nil {
				return err
			}

			sum := uint64(0)
			if err := image.DiffIterate(rbd.DiffIterateConfig{
				Offset: 0,
				Length: limitSize,
				Callback: func(o, l uint64, e int, x interface{}) int {
					sum += l
					return 0
				},
			}); err != nil {
				return errors.Wrap(err, "calc rbd size")
			}
			size = int64(sum)
		}
	} else {
		srcPath := filepath.Join(a.localBlockStorageDirectory, bs.Group, bs.Namespace, bs.ID)
		s, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer s.Close()
		src = s
		if finfo, err := s.Stat(); err != nil {
			return err
		} else {
			size = finfo.Size()
		}
	}
	destPath := filepath.Join(a.localImageDirectory, imageEntity.Group, imageEntity.ID)
	destDirPath := filepath.Dir(destPath)

	if err := os.MkdirAll(destDirPath, 0755); err != nil {
		return err
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.CopyN(dest, src, size); err != nil {
		return errors.Wrapf(err, "copy image from source bs")
	}

	hasher := sha256.New()
	if _, err := dest.Seek(0, 0); err != nil {
		return errors.Wrap(err, "seek dest file cursor")
	}

	if _, err := io.Copy(hasher, dest); err != nil {
		return err
	}

	imageEntity.Spec.Hash = fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	if imageEntity.Annotations == nil {
		imageEntity.Annotations = map[string]string{}
	}
	imageEntity.Annotations["image-entity-download-host"] = fmt.Sprintf("%s:%d", a.config.DownloadAPI.AdvertiseAddress, a.config.DownloadAPI.ListenPort)
	imageEntity.Status.State = system.ImageEntityStateAvailable
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	bs.Status.State = system.BlockStorageStateActive
	if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
		return err
	}

	return setHash(imageEntity)
}
