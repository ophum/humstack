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
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				log.Println("[BS] %s", err.Error())
				continue
			}
			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					log.Println("[BS] %s", err.Error())
					continue
				}

				for _, ns := range nsList {
					bsList, err := a.client.SystemV0().BlockStorage().List(group.ID, ns.ID)
					if err != nil {
						log.Printf("[BS] %s", err.Error())
						continue
					}

					vmList, err := a.client.SystemV0().VirtualMachine().List(group.ID, ns.ID)
					if err != nil {
						log.Printf("[BS] %s", err.Error())
						continue
					}
					usedBSIDs := []string{}
					for _, vm := range vmList {
						if vm.Status.State == system.VirtualMachineStateRunning {
							usedBSIDs = append(usedBSIDs, vm.Spec.BlockStorageIDs...)
						}
					}

					for _, bs := range bsList {
						oldHash := bs.ResourceHash

						// state check
						if bs.Status.State != system.BlockStorageStateDeleting &&
							bs.Status.State != system.BlockStorageStatePending {

							oldState := bs.Status.State
							bs.Status.State = system.BlockStorageStateActive
							log.Println("====")
							log.Println(usedBSIDs)
							for i, usedID := range usedBSIDs {
								if bs.ID == usedID {
									bs.Status.State = system.BlockStorageStateUsed
									usedBSIDs = append(usedBSIDs[:i], usedBSIDs[i+1:]...)
									break
								}
							}
							log.Println(usedBSIDs)
							log.Println("====")

							if bs.Status.State != oldState {
								bs, err = a.client.SystemV0().BlockStorage().Update(bs)
								if err != nil {
									log.Println(err)
									continue
								}
							}
						}

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
