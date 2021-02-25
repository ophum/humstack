package image

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ceph/go-ceph/rbd"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
)

func (a *ImageAgent) syncCephImageEntityFromBlockStorage(imageEntity *system.ImageEntity, bs *system.BlockStorage) error {

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

	var image *rbd.Image
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

		if image, err = rbd.OpenImageReadOnly(ioctx, imageName, ""); err != nil {
			return errors.Wrapf(err, "open rbd image `%s`", imageName)
		}

	} else {
		srcPath := filepath.Join(a.localBlockStorageDirectory, bs.Group, bs.Namespace, bs.ID)
		s, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer s.Close()

		var size uint64
		if finfo, err := s.Stat(); err != nil {
			return err
		} else {
			size = uint64(finfo.Size())
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

		imageName := "ceph-tmp-" + bs.Group + bs.Namespace + bs.ID
		imageEntity.Annotations[ImageEntityV0AnnotationCephImageName] = imageName

		cephImage, err := rbd.Create(ioctx, imageName, size, 22)
		if err != nil {
			return err
		}

		if err := cephImage.Open(); err != nil {
			return err
		}
		defer cephImage.Close()

		// 一時Imageのデータをcephのimageに書き込む
		if finfo, err := s.Stat(); err == nil {
			_, err = io.CopyN(cephImage, s, finfo.Size())
		} else {
			_, err = io.Copy(cephImage, s)
		}
		if err != nil {
			return err
		}

		// リサイズ
		imageNameFull := filepath.Join(a.config.CephBackend.PoolName, imageName)
		command := "qemu-img"
		args := []string{
			"resize",
			fmt.Sprintf("rbd:%s", imageNameFull),
			withUnitToWithoutUnit(bs.Spec.LimitSize),
		}
		cmd := exec.Command(command, args...)
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}

	}

	snapName := imageEntity.ID
	snapshot, err := image.CreateSnapshot(snapName)
	imageEntity.Annotations[ImageEntityV0AnnotationCephSnapName] = snapName
	if err != nil {
		return errors.Wrap(err, "Failed to create ceph snapshot from ceph image.")
	}

	err = snapshot.Protect()
	if err != nil {
		return errors.Wrap(err, "Failed to protect ceph snapshot.")
	}

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
