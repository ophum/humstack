package image

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strconv"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ImageAgent struct {
	client                     *client.Clients
	logger                     *zap.Logger
	nodeName                   string
	config                     *ImageAgentConfig
	localImageDirectory        string
	localBlockStorageDirectory string
}

const (
	ImageEntityV0AnnotationType = "imageentityv0/type"
)

const (
	ImageEntityV0ImageEntityTypeLocal = "Local"
)

func NewImageAgent(client *client.Clients, config *ImageAgentConfig, logger *zap.Logger) *ImageAgent {
	nodeName, err := os.Hostname()
	if err != nil {
		logger.Panic(
			"get hostname",
			zap.String("msg", err.Error()),
			zap.Time("time", time.Now()),
		)
	}
	return &ImageAgent{
		client:                     client,
		logger:                     logger,
		nodeName:                   nodeName,
		config:                     config,
		localImageDirectory:        config.ImageDirPath,
		localBlockStorageDirectory: config.BlockStorageDirPath,
	}
}

func (a *ImageAgent) Run() {

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				a.logger.Error(
					"get group list",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
				continue
			}

			for _, group := range grList {
				imageEntityList, err := a.client.SystemV0().ImageEntity().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get imageentity list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}

				for _, imageEntity := range imageEntityList {
					oldHash := imageEntity.ResourceHash

					if imageEntity.Status.State == system.ImageEntityStateAvailable {
						continue
					}

					bs, err := a.client.SystemV0().BlockStorage().Get(
						imageEntity.Group,
						imageEntity.Spec.Source.Namespace,
						imageEntity.Spec.Source.BlockStorageID)

					if err != nil {
						a.logger.Error(
							"get blockstorage list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					nodeName := bs.Annotations[blockstorage.BlockStorageV0AnnotationNodeName]

					// 別のノードのBSの場合は何もしない
					if nodeName != a.nodeName {
						continue
					}

					// ファイルがなければPENDINGとして扱う
					if _, err := os.Stat(filepath.Join(a.localImageDirectory, imageEntity.Group, imageEntity.ID)); err != nil {
						imageEntity.Status.State = system.ImageEntityStatePending
					}

					// とりあえずPending以外になってたら何もしない
					if imageEntity.Status.State != "" && imageEntity.Status.State != system.ImageEntityStatePending && imageEntity.DeleteState != meta.DeleteStateDelete {
						continue
					}

					entityType, ok := imageEntity.Annotations[ImageEntityV0AnnotationType]
					if !ok {
						entityType = ImageEntityV0ImageEntityTypeLocal
					}

					switch entityType {
					case ImageEntityV0ImageEntityTypeLocal:
						if err := a.syncLocalImageEntity(imageEntity, bs); err != nil {
							a.logger.Error(
								"sync local imageentity",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}
					}

					if imageEntity.ResourceHash == oldHash {
						continue
					}

					if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
						a.logger.Error(
							"update imageentity",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}
				}
			}
		}
	}
}

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

func (a *ImageAgent) syncCephImageEntity(imageEntity *system.ImageEntity, bs *system.BlockStorage) error {

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

	snapshot, err := image.CreateSnapshot(imageEntity.ID)
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

func setHash(imageEntity *system.ImageEntity) error {
	imageEntity.ResourceHash = ""
	resourceJSON, err := json.Marshal(imageEntity)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	imageEntity.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}

// TODO: BlockStorageAgentにも同じ実装がある
func (a ImageAgent) newCephConn() (*rados.Conn, error) {
	cephConn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}

	if err := cephConn.ReadConfigFile(a.config.CephBackend.ConfigPath); err != nil {
		return nil, err
	}

	if err := cephConn.Connect(); err != nil {
		return nil, err
	}
	return cephConn, nil
}

func fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

const (
	UnitGigabyte = 'G'
	UnitMegabyte = 'M'
	UnitKilobyte = 'K'
)

func withUnitToWithoutUnit(numberWithUnit string) string {
	length := len(numberWithUnit)
	if numberWithUnit[length-1] >= '0' && numberWithUnit[length-1] <= '9' {
		return numberWithUnit
	}

	number, err := strconv.ParseInt(numberWithUnit[:length-1], 10, 64)
	if err != nil {
		return "0"
	}
	switch numberWithUnit[length-1] {
	case UnitGigabyte:
		return fmt.Sprintf("%d", number*1024*1024*1024)
	case UnitMegabyte:
		return fmt.Sprintf("%d", number*1024*1024)
	case UnitKilobyte:
		return fmt.Sprintf("%d", number*1024)
	}
	return "0"
}
