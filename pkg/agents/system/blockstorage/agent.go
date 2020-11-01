package blockstorage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
)

type BlockStorageAgent struct {
	client                     *client.Clients
	config                     *BlockStorageAgentConfig
	localBlockStorageDirectory string
	localImageDirectory        string
	logger                     *zap.Logger
}

const (
	BlockStorageV0AnnotationType     = "blockstoragev0/type"
	BlockStorageV0AnnotationNodeName = "blockstoragev0/node_name"
)

const (
	BlockStorageV0BlockStorageTypeLocal = "Local"
)

func NewBlockStorageAgent(client *client.Clients, config *BlockStorageAgentConfig, logger *zap.Logger) *BlockStorageAgent {
	return &BlockStorageAgent{
		client:                     client,
		config:                     config,
		localBlockStorageDirectory: config.BlockStorageDirPath,
		localImageDirectory:        config.ImageDirPath,
		logger:                     logger,
	}
}

func (a *BlockStorageAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		a.logger.Panic(
			"get hostname",
			zap.String("msg", err.Error()),
			zap.Time("time", time.Now()))
	}

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
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get namespace list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}

				for _, ns := range nsList {
					bsList, err := a.client.SystemV0().BlockStorage().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get blockstorage list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					vmList, err := a.client.SystemV0().VirtualMachine().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get virtualmahine list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
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
							for i, usedID := range usedBSIDs {
								if bs.ID == usedID {
									bs.Status.State = system.BlockStorageStateUsed
									usedBSIDs = append(usedBSIDs[:i], usedBSIDs[i+1:]...)
									break
								}
							}

							if bs.Status.State != oldState {
								bs, err = a.client.SystemV0().BlockStorage().Update(bs)
								if err != nil {
									a.logger.Error(
										"update blockstorage",
										zap.String("msg", err.Error()),
										zap.Time("time", time.Now()),
									)
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
								a.logger.Error(
									"sync local blockstorage",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}
						}

						if bs.ResourceHash == oldHash {
							continue
						}

						_, err := a.client.SystemV0().BlockStorage().Update(bs)
						if err != nil {
							a.logger.Error(
								"update blockstorage",
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
