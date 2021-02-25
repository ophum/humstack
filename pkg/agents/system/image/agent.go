package image

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ceph/go-ceph/rados"
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
	ImageEntityV0AnnotationType          = "imageentityv0/type"
	ImageEntityV0AnnotationCephSnapName  = "imageentityv0/ceph-snapname"
	ImageEntityV0AnnotationCephImageName = "imageentityv0/ceph-imagename"
)

const (
	ImageEntityV0ImageEntityTypeLocal = "Local"
	ImageEntityV0ImageEntityTypeCeph  = "Ceph"
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
						// save local as default place
						entityType = ImageEntityV0ImageEntityTypeLocal
						if imageEntity.Spec.Type == "Ceph" {
							entityType = ImageEntityV0ImageEntityTypeCeph
						} else if imageEntity.Spec.Type == "Local" {
							entityType = ImageEntityV0ImageEntityTypeLocal
						} else {
							errors.Errorf("Image type value is invalid: ", imageEntity.Spec.Type)
						}
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
					case ImageEntityV0ImageEntityTypeCeph:
						if err := a.syncCephImageEntity(imageEntity, bs); err != nil {
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
