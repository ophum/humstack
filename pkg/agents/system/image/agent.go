package image

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
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
					if imageEntity.Status.State != "" && imageEntity.Status.State != system.ImageEntityStatePending {
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

	srcPath := filepath.Join(a.localBlockStorageDirectory, bs.Group, bs.Namespace, bs.ID)
	destPath := filepath.Join(a.localImageDirectory, imageEntity.Group, imageEntity.ID)
	destDirPath := filepath.Dir(destPath)

	if err := os.MkdirAll(destDirPath, 0755); err != nil {
		return err
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return err
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
