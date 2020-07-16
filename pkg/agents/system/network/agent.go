package network

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type NetworkAgent struct {
	client *client.Clients
}

const (
	NetworkV0AnnotationNetworkType = "networkv0/network_type"
)

const (
	NetworkV0NetworkTypeVXLAN  = "VXLAN"
	NetworkV0NetworkTypeVLAN   = "VLAN"
	NetworkV0NetworkTypeBridge = "Bridge"
)

func NewNetworkAgent(client *client.Clients) *NetworkAgent {
	return &NetworkAgent{
		client: client,
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
				log.Printf("[NET] %s", err.Error())
				continue
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					continue
				}

				for _, ns := range nsList {
					netList, err := a.client.SystemV0().Network().List(group.ID, ns.ID)
					if err != nil {
						continue
					}

					for _, net := range netList {
						oldHash := net.ResourceHash
						switch net.Annotations[NetworkV0AnnotationNetworkType] {
						case NetworkV0NetworkTypeBridge:
							err = syncBridgeNetwork(net)
							if err != nil {
								log.Println("error sync bridge network")
								log.Println(err)
								continue
							}

						}

						if net.ResourceHash == oldHash {
							log.Println("no update")
							continue
						}
						log.Println("update store")
						_, err = a.client.SystemV0().Network().Update(net)
						if err != nil {
							log.Println("error Update store")
							log.Println(err)
							continue
						}
					}
				}
			}
		}
	}
}

func setHash(network *system.Network) error {
	network.ResourceHash = ""
	resourceJSON, err := json.Marshal(network)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	network.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
