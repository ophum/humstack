package image

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type ImageAgent struct {
	client                     *client.Clients
	nodeName                   string
	localImageDirectory        string
	localBlockStorageDirectory string
}

const (
	ImageEntityV0AnnotationType = "imageentityv0/type"
)

const (
	ImageEntityV0ImageEntityTypeLocal = "Local"
)

func NewImageAgent(client *client.Clients, localImageDirectory, localBlockStorageDirectory string) *ImageAgent {
	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &ImageAgent{
		client:                     client,
		nodeName:                   nodeName,
		localImageDirectory:        localImageDirectory,
		localBlockStorageDirectory: localBlockStorageDirectory,
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
				log.Println(err)
				continue
			}

			for _, group := range grList {
				imageEntityList, err := a.client.SystemV0().ImageEntity().List(group.ID)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, imageEntity := range imageEntityList {
					log.Println("imageEntity: ", imageEntity)
					oldHash := imageEntity.ResourceHash
					bs, err := a.client.SystemV0().BlockStorage().Get(
						imageEntity.Group,
						imageEntity.Spec.Source.Namespace,
						imageEntity.Spec.Source.BlockStorageID)

					if err != nil {
						continue
					}

					nodeName := bs.Annotations[blockstorage.BlockStorageV0AnnotationNodeName]

					// 別のノードのBSの場合は何もしない
					if nodeName != a.nodeName {
						log.Println(nodeName)
						log.Println(a.nodeName)
						log.Println("ImageEntity: other node")
						continue
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
							log.Println(err)
							continue
						}
					}

					if imageEntity.ResourceHash == oldHash {
						continue
					}

					if _, err := a.client.SystemV0().ImageEntity().Update(imageEntity); err != nil {
						log.Println(err)
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
	if _, err := io.Copy(hasher, dest); err != nil {
		return err
	}

	imageEntity.Spec.Hash = fmt.Sprintf("sha256:%x", hasher.Sum(nil))

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
