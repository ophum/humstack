package image

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ceph/go-ceph/rbd"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
)

func (a *ImageAgent) syncCephImageEntityFromImage(imageEntity *system.ImageEntity) error {
	if imageEntity.DeleteState == meta.DeleteStateDelete {
		return a.deleteCephImageEntity(imageEntity)
	}

	sourceImageEntity, err := a.getImageEntityByImage(
		imageEntity.Group,
		imageEntity.Spec.Source.ImageName,
		imageEntity.Spec.Source.ImageTag,
	)
	if err != nil {
		return err
	}

	// コピー元がAvailableではない
	if sourceImageEntity.Status.State != system.ImageEntityStateAvailable {
		return nil
	}

	imageEntity.Status.State = system.ImageEntityStateCopying
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	var image *rbd.Image
	if sourceImageEntity.Spec.Type == "Ceph" {
		return fmt.Errorf("not implements")
	} else { // Local or ""
		stream, size, err := a.client.SystemV0().Image().Download(
			sourceImageEntity.Group,
			imageEntity.Spec.Source.ImageName,
			imageEntity.Spec.Source.ImageTag,
		)
		if err != nil {
			return err
		}
		defer stream.Close()

		conn, err := a.newCephConn()
		if err != nil {
			return err
		}
		defer conn.Shutdown()

		ioctx, err := conn.OpenIOContext(a.config.CephBackend.PoolName)
		if err != nil {
			return err
		}
		defer ioctx.Destroy()

		imageName := getSnapParentImageName(imageEntity)
		imageEntity.Annotations[ImageEntityV0AnnotationCephImageName] = imageName

		image, err = rbd.Create(ioctx, imageName, uint64(size), 22)
		if err != nil {
			return err
		}

		if err := image.Open(); err != nil {
			return err
		}
		defer image.Close()

		if _, err := io.CopyN(image, stream, int64(size)); err != nil {
			return err
		}
	}

	snapName := imageEntity.ID
	snapshot, err := image.CreateSnapshot(snapName)
	if err != nil {
		return errors.Wrap(err, "Failed to create ceph snapshot from ceph image.")
	}

	if imageEntity.Annotations == nil {
		imageEntity.Annotations = map[string]string{}
	}
	imageEntity.Annotations[ImageEntityV0AnnotationCephSnapName] = snapName

	if err := snapshot.Protect(); err != nil {
		return errors.Wrap(err, "Failed to protect ceph snapshot.")
	}

	imageEntity.Annotations["image-entity-download-host"] = fmt.Sprintf("%s:%d", a.config.DownloadAPI.AdvertiseAddress, a.config.DownloadAPI.ListenPort)
	imageEntity.Status.State = system.ImageEntityStateAvailable
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	return setHash(imageEntity)
}

func (a *ImageAgent) syncCephImageEntityFromBlockStorage(imageEntity *system.ImageEntity, bs *system.BlockStorage) error {
	if imageEntity.DeleteState == meta.DeleteStateDelete {
		return a.deleteCephImageEntity(imageEntity)
	}

	imageEntity.Status.State = system.ImageEntityStatePending
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	if bs.Status.State != system.BlockStorageStateActive {
		return nil
	}

	imageEntity.Status.State = system.ImageEntityStateCopying
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}
	bs.Status.State = system.BlockStorageStateCopying
	if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
		return err
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

	var image *rbd.Image
	if t, ok := bs.Annotations["blockstoragev0/type"]; ok && t == "Ceph" {
		imageName, ok := bs.Annotations["ceph-image-name"]
		if !ok {
			return fmt.Errorf("ceph-image-name not found")
		}

		if fromImage, err := rbd.OpenImageReadOnly(ioctx, imageName, ""); err != nil {
			return errors.Wrapf(err, "open rbd image `%s`", imageName)
		} else {
			defer fromImage.Close()
			size, err := fromImage.GetSize()
			if err != nil {
				return errors.Wrapf(err, "Failed to get image size `%s`", imageName)
			}

			destImageName := getSnapParentImageName(imageEntity)
			image, err = rbd.Create(ioctx, destImageName, size, 22)
			if err != nil {
				return errors.Wrapf(err, "Failed to create image `%s`", destImageName)
			}
			if err := image.Open(); err != nil {
				return errors.Wrapf(err, "Failed to open image `%s`", destImageName)
			}
			defer image.Close()

			if err := fromImage.Copy2(image); err != nil {
				return errors.Wrapf(err, "Failed to copy image from `%s` to `%s`.", imageName, destImageName)
			}
			imageEntity.Annotations[ImageEntityV0AnnotationCephImageName] = destImageName
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

		imageName := getSnapParentImageName(imageEntity)
		imageEntity.Annotations[ImageEntityV0AnnotationCephImageName] = imageName

		image, err = rbd.Create(ioctx, imageName, size, 22)
		if err != nil {
			return err
		}

		if err := image.Open(); err != nil {
			return err
		}
		defer image.Close()

		// 一時Imageのデータをcephのimageに書き込む
		if finfo, err := s.Stat(); err == nil {
			_, err = io.CopyN(image, s, finfo.Size())
		} else {
			_, err = io.Copy(image, s)
		}
		if err != nil {
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

func (a ImageAgent) cephImageIsExists(imageEntity *system.ImageEntity) bool {
	// typeがCephでない
	if imageEntity.Spec.Type != ImageEntityV0ImageEntityTypeCeph {
		return false
	}

	imageName := imageEntity.Annotations[ImageEntityV0AnnotationCephImageName]

	conn, err := a.newCephConn()
	if err != nil {
		return false
	}
	defer conn.Shutdown()

	ioctx, err := conn.OpenIOContext(a.config.CephBackend.PoolName)
	if err != nil {
		return false
	}
	defer ioctx.Destroy()

	if image, err := rbd.OpenImageReadOnly(ioctx, imageName, ""); err != nil {
		return false
	} else {
		defer image.Close()
	}
	return true
}

func (a *ImageAgent) getImageEntityByImage(group, id, tag string) (*system.ImageEntity, error) {
	image, err := a.client.SystemV0().Image().Get(group, id)
	if err != nil {
		return nil, err
	}
	imageEntityID, ok := image.Spec.EntityMap[tag]
	if !ok {
		return nil, fmt.Errorf("image tag is not found")
	}

	return a.client.SystemV0().ImageEntity().Get(group, imageEntityID)
}

func (a *ImageAgent) deleteCephImageEntity(imageEntity *system.ImageEntity) error {
	if imageEntity.Status.State != "" &&
		imageEntity.Status.State != system.ImageEntityStateAvailable {
		return nil
	}

	imageEntity.Status.State = system.ImageEntityStateDeleting
	if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
		return err
	}

	conn, err := a.newCephConn()
	if err != nil {
		imageEntity.Status.State = ""
		if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
			return err
		}
		return errors.Wrap(err, "Failed to create ceph connection.")
	}
	defer conn.Shutdown()

	ioctx, err := conn.OpenIOContext(a.config.CephBackend.PoolName)
	if err != nil {
		imageEntity.Status.State = ""
		if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
			return err
		}
		return errors.Wrap(err, "Failed to open io context")
	}
	defer ioctx.Destroy()

	if a.cephImageIsExists(imageEntity) {
		imageName := imageEntity.Annotations[ImageEntityV0AnnotationCephImageName]
		image, err := rbd.OpenImage(ioctx, imageName, rbd.NoSnapshot)
		if err != nil {
			return errors.Wrapf(err, "Failed to open image `%s`", imageName)
		}
		defer image.Close()

		// snapshotをすべて消す
		if snaps, err := image.GetSnapshotNames(); err != nil {
			return errors.Wrapf(err, "Failed to get snapshot names by image `%s`", imageName)
		} else {
			for _, snap := range snaps {
				snapshot := image.GetSnapshot(snap.Name)
				if is, err := snapshot.IsProtected(); err != nil {
					return err
				} else if is {
					if err := snapshot.Unprotect(); err != nil {
						return errors.Wrapf(err, "Failed to unprotect snapshot `%s` (image: `%s`)", snap.Name, imageName)
					}
				}

				if err := snapshot.Remove(); err != nil {
					return errors.Wrapf(err, "Failed to remove snapshot `%s` (image: `%s`)", snap.Name, imageName)
				}
			}
		}

		// イメージを削除
		image.Close()
		if err := image.Remove(); err != nil {
			imageEntity.Status.State = ""
			if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
				return err
			}
			return errors.Wrapf(err, "Failed to remove image `%s`", imageName)
		}
	}

	if err := a.client.SystemV0().ImageEntity().Delete(imageEntity.Group, imageEntity.ID); err != nil {
		imageEntity.Status.State = ""
		if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
			return err
		}
		return err
	}
	return nil
}

func getSnapParentImageName(imageEntity *system.ImageEntity) string {
	return "ceph-tmp-" + imageEntity.ID
}
