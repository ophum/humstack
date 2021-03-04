package blockstorage

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type BlockStorageAgent struct {
	client                     *client.Clients
	config                     *BlockStorageAgentConfig
	localBlockStorageDirectory string
	localImageDirectory        string
	parallelSemaphore          *semaphore.Weighted
	logger                     *zap.Logger
}

const (
	BlockStorageV0AnnotationType     = "blockstoragev0/type"
	BlockStorageV0AnnotationNodeName = "blockstoragev0/node_name"
)

const (
	BlockStorageV0BlockStorageTypeLocal = "Local"
	BlockStorageV0BlockStorageTypeCeph  = "Ceph"
)

func NewBlockStorageAgent(client *client.Clients, config *BlockStorageAgentConfig, logger *zap.Logger) *BlockStorageAgent {
	return &BlockStorageAgent{
		client:                     client,
		config:                     config,
		localBlockStorageDirectory: config.BlockStorageDirPath,
		localImageDirectory:        config.ImageDirPath,
		parallelSemaphore:          semaphore.NewWeighted(config.ParallelLimit),
		logger:                     logger,
	}
}

func (a *BlockStorageAgent) Run(pollingDuration time.Duration) {
	ticker := time.NewTicker(pollingDuration)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		a.logger.Panic(
			"get hostname",
			zap.String("msg", err.Error()),
			zap.Time("time", time.Now()))
	}
	// init
	grList, err := a.client.CoreV0().Group().List()
	if err != nil {
		a.logger.Error(
			"get group list",
			zap.String("msg", err.Error()),
			zap.Time("time", time.Now()),
		)
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
			for _, bs := range bsList {
				switch bs.Status.State {
				case system.BlockStorageStateCopying, system.BlockStorageStateDownloading, system.BlockStorageStateDeleting, system.BlockStorageStateQueued:
					bs.Status.State = ""
					if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
						a.logger.Panic(
							"init state Copying or Downloading or Deleting or Queued => ``",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()))

					}
				}
			}
		}
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

			wg := sync.WaitGroup{}
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
						if bs.DeleteState != meta.DeleteStateDelete && bs.Status.State == system.BlockStorageStateQueued {
							continue
						}

						err := a.parallelSemaphore.Acquire(context.TODO(), 1)
						if err != nil {
							a.logger.Error(
								"acqure semaphre",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}

						//if bs.Status.State == "" {
						//	bs.Status.State = system.BlockStorageStateQueued
						//	if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
						//		a.logger.Error(
						//			"update blockstorage state",
						//			zap.String("msg", err.Error()),
						//			zap.Time("time", time.Now()),
						//		)
						//		continue
						//	}
						//}
						wg.Add(1)
						copiedBS := *bs
						go func(usedBSIDs []string, bs *system.BlockStorage) {
							defer func() {
								a.parallelSemaphore.Release(1)
								wg.Done()
							}()
							if bs.Annotations[BlockStorageV0AnnotationNodeName] != nodeName {
								return
							}
							oldHash := bs.ResourceHash

							// state check
							if bs.Status.State != system.BlockStorageStateDeleting &&
								bs.Status.State != system.BlockStorageStatePending &&
								bs.Status.State != system.BlockStorageStateQueued {

								isUsed := false
								for i, usedID := range usedBSIDs {
									if bs.ID == usedID {
										isUsed = true
										usedBSIDs = append(usedBSIDs[:i], usedBSIDs[i+1:]...)
										break
									}
								}

								if bs.Status.State != system.BlockStorageStateUsed && isUsed {
									bs.Status.State = system.BlockStorageStateUsed
									_, err := a.client.SystemV0().BlockStorage().Update(bs)
									if err != nil {
										a.logger.Error(
											"update blockstorage",
											zap.String("msg", err.Error()),
											zap.Time("time", time.Now()),
										)
										return
									}
								} else if bs.Status.State == system.BlockStorageStateUsed && !isUsed {
									bs.Status.State = system.BlockStorageStateActive
									bs, err = a.client.SystemV0().BlockStorage().Update(bs)
									if err != nil {
										a.logger.Error(
											"update blockstorage",
											zap.String("msg", err.Error()),
											zap.Time("time", time.Now()),
										)
										return
									}
								}

								if isUsed {
									a.logger.Info(
										"skip bs is used",
										zap.String("bs", bs.Namespace+"/"+bs.ID),
										zap.Time("time", time.Now()),
									)
									return
								}
							}

							switch bs.Annotations[BlockStorageV0AnnotationType] {
							case BlockStorageV0BlockStorageTypeLocal:
								if bs.Annotations[BlockStorageV0AnnotationNodeName] != nodeName {
									return
								}

								err = a.syncLocalBlockStorage(bs)
								if err != nil {
									a.logger.Error(
										"sync local blockstorage",
										zap.String("msg", err.Error()),
										zap.Time("time", time.Now()),
									)
									return
								}

							case BlockStorageV0BlockStorageTypeCeph:
								if bs.Annotations[BlockStorageV0AnnotationNodeName] != nodeName {
									return
								}

								err = a.syncCephBlockStorage(bs)
								if err != nil {
									a.logger.Error(
										"sync ceph blockstorage",
										zap.String("msg", err.Error()),
										zap.Time("time", time.Now()),
									)
									return
								}
							}

							if bs.ResourceHash == oldHash {
								return
							}

							_, err := a.client.SystemV0().BlockStorage().Update(bs)
							if err != nil {
								a.logger.Error(
									"update blockstorage",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								return
							}
							a.logger.Info("sync because different hash",
								zap.String("bs", bs.Namespace+"/"+bs.ID),
								zap.Time("time", time.Now()),
							)
						}(usedBSIDs, &copiedBS)
					}
				}
			}
			wg.Wait()
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
