package network

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type NetworkAgent struct {
	client *client.Clients
	logger *zap.Logger
}

func NewNetworkAgent(client *client.Clients, logger *zap.Logger) *NetworkAgent {
	return &NetworkAgent{
		client: client,
		logger: logger,
	}
}

func (a *NetworkAgent) Run() {
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
					netList, err := a.client.CoreV0().Network().List(group.ID, ns.ID)

					if err != nil {
						a.logger.Error(
							"get virtualmachine list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					for _, net := range netList {
						oldHash := net.ResourceHash
						err = a.syncNetwork(net)
						if err != nil {
							a.logger.Error(
								"sync bridge network",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}

						if net.ResourceHash == oldHash {
							continue
						}

						_, err = a.client.CoreV0().Network().Update(net)
						if err != nil {
							a.logger.Error(
								"update network",
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

func (a *NetworkAgent) syncNetwork(net *core.Network) error {
	nodeList, err := a.client.SystemV0().Node().List()
	if err != nil {
		return err
	}

	// 削除処理
	if net.DeleteState == meta.DeleteStateDelete {
		isDelete := true
		for _, node := range nodeList {
			nodeNet, err := a.client.SystemV0().NodeNetwork().Get(net.Group, net.Namespace, getNodeNetworkID(net.ID, node.ID))
			if err != nil {
				a.logger.Error(
					"set delete state",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
				continue
			}

			if nodeNet.ID != "" {
				isDelete = false
				if nodeNet.DeleteState == meta.DeleteStateDelete {
					continue
				}
				if err := a.client.SystemV0().NodeNetwork().DeleteState(net.Group, net.Namespace, getNodeNetworkID(net.ID, node.ID)); err != nil {
					a.logger.Error(
						"set delete state",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}
			}
		}

		if isDelete {
			if err := a.client.CoreV0().Network().Delete(net.Group, net.Namespace, net.ID); err != nil {
				return errors.Wrap(err, "delete network")
			}
		}
		return nil
	}
	// 各ノードに作られていなければ作成する
	for _, node := range nodeList {
		nodeNet, err := a.client.SystemV0().NodeNetwork().Get(net.Group, net.Namespace, fmt.Sprintf("%s_%s", net.ID, node.ID))
		if err != nil {
			a.logger.Error(
				"get node network",
				zap.String("msg", err.Error()),
				zap.Time("time", time.Now()),
			)
			continue
		}
		if nodeNet.ID == "" {
			nodeNet := &system.NodeNetwork{
				Meta: meta.Meta{
					ID:          fmt.Sprintf("%s_%s", net.ID, node.ID),
					Name:        fmt.Sprintf("%s_%s", net.ID, node.ID),
					Namespace:   net.Namespace,
					Group:       net.Group,
					Annotations: net.Spec.Template.Annotations,
					OwnerReferences: []meta.OwnerReference{
						{
							Meta: net.Meta,
						},
					},
				},
				Spec: net.Spec.Template.Spec,
			}
			nodeNet.Annotations["nodenetworkv0/node_name"] = node.ID
			if _, err := a.client.SystemV0().NodeNetwork().Create(nodeNet); err != nil {
				a.logger.Error(
					"create node network",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
				continue
			}
		}
	}

	return setHash(net)
}

func getNodeNetworkID(networkID, nodeID string) string {
	return fmt.Sprintf("%s_%s", networkID, nodeID)
}
func setHash(network *core.Network) error {
	network.ResourceHash = ""
	resourceJSON, err := json.Marshal(network)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	network.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
