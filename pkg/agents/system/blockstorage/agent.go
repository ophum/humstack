package blockstorage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type BlockStorageAgent struct {
	client                     *client.Clients
	localBlockStorageDirectory string
}

const (
	BlockStorageV0AnnotationType     = "blockstoragev0/type"
	BlockStorageV0AnnotationNodeName = "blockstoragev0/node_name"
)

const (
	BlockStorageV0BlockStorageTypeLocal = "Local"
)

func NewBlockStorageAgent(client *client.Clients, localBlockStorageDirectory string) *BlockStorageAgent {
	return &BlockStorageAgent{
		client:                     client,
		localBlockStorageDirectory: localBlockStorageDirectory,
	}
}

func (a *BlockStorageAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ticker.C:
			nsList, err := a.client.CoreV0().Namespace().List()
			if err != nil {
				continue
			}

			for _, ns := range nsList {
				bsList, err := a.client.SystemV0().BlockStorage().List(ns.ID)
				if err != nil {
					continue
				}

				for _, bs := range bsList {
					oldHash := bs.ResourceHash
					switch bs.Annotations[BlockStorageV0AnnotationType] {
					case BlockStorageV0BlockStorageTypeLocal:
						if bs.Annotations[BlockStorageV0AnnotationNodeName] != nodeName {
							continue
						}

						err = a.syncLocalBlockStorage(bs)
						if err != nil {
							log.Println(err)
							continue
						}
					}

					if bs.ResourceHash == oldHash {
						continue
					}

					_, err := a.client.SystemV0().BlockStorage().Update(bs)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		}
	}
}

func setHash(bs *system.BlockStorage) error {
	bs.ResourceHash = ""
	resourceJSON, err := json.Marshal(bs)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	bs.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
