package main

import (
	"flag"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/agents/system/network"
	"github.com/ophum/humstack/pkg/agents/system/node"
	"github.com/ophum/humstack/pkg/agents/system/virtualmachine"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

type Config struct {
	apiServerAddress    string `yaml:"apiServerAddress"`
	LimitMemory         string `yaml:"limitMemory"`
	LimitVcpus          string `yaml:"limitVcpus"`
	BlockStorageDirPath string `yaml:"blockStorageDirPath"`

	NetworkAgentConfig network.NetworkAgentConfig `yaml:"networkAgentConfig"`
}

var (
	config Config = Config{}
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()

	if configPath == "" {
		log.Fatal("unexpected --config")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal("error open config file")
	}

	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("failed decode config")
	}

}

func main() {
	client := client.NewClients(config.apiServerAddress, 8080)
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	nodeAgent := node.NewNodeAgent(system.Node{
		Meta: meta.Meta{
			ID:   hostname,
			Name: hostname,
		},
		Spec: system.NodeSpec{
			LimitMemory: config.LimitMemory,
			LimitVcpus:  config.LimitVcpus,
		},
	}, client)

	netAgent := network.NewNetworkAgent(client, &config.NetworkAgentConfig)

	bsAgent := blockstorage.NewBlockStorageAgent(client, config.BlockStorageDirPath)

	vmAgent := virtualmachine.NewVirtualMachineAgent(client)

	go nodeAgent.Run()
	go bsAgent.Run()
	go vmAgent.Run()
	netAgent.Run()

}
