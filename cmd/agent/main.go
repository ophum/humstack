package main

import (
	"flag"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/agents/system/image"
	"github.com/ophum/humstack/pkg/agents/system/network"
	"github.com/ophum/humstack/pkg/agents/system/node"
	"github.com/ophum/humstack/pkg/agents/system/virtualmachine"
	"github.com/ophum/humstack/pkg/agents/system/virtualrouter"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ApiServerAddress string `yaml:"apiServerAddress"`
	ApiServerPort    int32  `yaml:"apiServerPort"`
	LimitMemory      string `yaml:"limitMemory"`
	LimitVcpus       string `yaml:"limitVcpus"`

	BlockStorageAgentConfig blockstorage.BlockStorageAgentConfig `yaml:"blockStorageAgentConfig"`

	NetworkAgentConfig network.NetworkAgentConfig `yaml:"networkAgentConfig"`

	ImageAgentConfig image.ImageAgentConfig `yaml:"imageAgentConfig"`
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

	log.Println(config)
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	client := client.NewClients(config.ApiServerAddress, config.ApiServerPort)
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
	}, client,
		logger.With(zap.Namespace("NodeAgent")),
	)

	netAgent := network.NewNetworkAgent(
		client,
		&config.NetworkAgentConfig,
		logger.With(zap.Namespace("NetworkAgent")),
	)

	bsAgent := blockstorage.NewBlockStorageAgent(
		client,
		&config.BlockStorageAgentConfig,
		logger.With(zap.Namespace("BlockStorageAgent")),
	)

	imAgent := image.NewImageAgent(
		client,
		&config.ImageAgentConfig,
		logger.With(zap.Namespace("ImageAgent")),
	)

	vmAgent := virtualmachine.NewVirtualMachineAgent(
		client,
		logger.With(zap.Namespace("VirtualMachineAgent")),
	)

	vrAgent := virtualrouter.NewVirtualRouterAgent(
		client,
		"exBr",
		"10.0.0.0/24",
		[]string{"10.0.0.1", "10.0.0.2"},
		logger.With(zap.Namespace("VirtualRouterAgent")),
	)

	log.Println(config.ImageAgentConfig.DownloadAPI)
	go nodeAgent.Run()
	go bsAgent.Run()
	go bsAgent.DownloadAPI(&config.BlockStorageAgentConfig.DownloadAPI)
	go imAgent.Run()
	go imAgent.DownloadAPI(&config.ImageAgentConfig.DownloadAPI)
	go vmAgent.Run()
	go vrAgent.Run()
	netAgent.Run()

}
