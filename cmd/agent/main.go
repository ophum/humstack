package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ophum/humstack/pkg/agents/core/group"
	"github.com/ophum/humstack/pkg/agents/core/namespace"
	"github.com/ophum/humstack/pkg/agents/core/network"
	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/agents/system/image"
	"github.com/ophum/humstack/pkg/agents/system/node"
	"github.com/ophum/humstack/pkg/agents/system/nodenetwork"
	"github.com/ophum/humstack/pkg/agents/system/virtualmachine"
	"github.com/ophum/humstack/pkg/agents/system/virtualrouter"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type AgentMode string

const (
	AgentModeCore   AgentMode = "Core"
	AgentModeSystem AgentMode = "System"
	AgentModeAll    AgentMode = "All"
)

type Config struct {
	AgentMode        AgentMode `yaml:"agentMode"`
	ApiServerAddress string    `yaml:"apiServerAddress"`
	ApiServerPort    int32     `yaml:"apiServerPort"`
	LimitMemory      string    `yaml:"limitMemory"`
	LimitVcpus       string    `yaml:"limitVcpus"`
	NodeAddress      string    `yaml:"nodeAddress"`
	PollingSeconds   int       `yaml:"pollingSeconds"`

	BlockStorageAgentConfig blockstorage.BlockStorageAgentConfig `yaml:"blockStorageAgentConfig"`

	NetworkAgentConfig nodenetwork.NetworkAgentConfig `yaml:"networkAgentConfig"`

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

	pollingDuration := time.Second * time.Duration(config.PollingSeconds)

	nodeAgent := node.NewNodeAgent(&system.Node{
		Meta: meta.Meta{
			ID:   hostname,
			Name: hostname,
			Annotations: map[string]string{
				"agentMode": string(config.AgentMode),
			},
		},
		Spec: system.NodeSpec{
			Address:     config.NodeAddress,
			LimitMemory: config.LimitMemory,
			LimitVcpus:  config.LimitVcpus,
		},
	}, client,
		logger.With(zap.Namespace("NodeAgent")),
	)
	go nodeAgent.Run(pollingDuration)

	if config.AgentMode == AgentModeAll || config.AgentMode == AgentModeCore {
		grAgent := group.NewGroupAgent(
			client,
			logger.With(zap.Namespace("GroupAgent")),
		)

		nsAgent := namespace.NewNamespaceAgent(
			client,
			logger.With(zap.Namespace("NamespaceAgent")),
		)

		netAgent := network.NewNetworkAgent(
			client,
			logger.With(zap.Namespace("NetworkAgent")),
		)

		go grAgent.Run(pollingDuration)
		go nsAgent.Run(pollingDuration)
		go netAgent.Run(pollingDuration)
	}

	if config.AgentMode == AgentModeAll || config.AgentMode == AgentModeSystem {
		nodeNetAgent := nodenetwork.NewNodeNetworkAgent(
			client,
			&config.NetworkAgentConfig,
			logger.With(zap.Namespace("NodeNetworkAgent")),
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
		go bsAgent.Run(pollingDuration)
		go bsAgent.DownloadAPI(&config.BlockStorageAgentConfig.DownloadAPI)
		go imAgent.Run(pollingDuration)
		go imAgent.DownloadAPI(&config.ImageAgentConfig.DownloadAPI)
		go vmAgent.Run(pollingDuration)
		go vrAgent.Run(pollingDuration)
		go nodeNetAgent.Run(pollingDuration)

	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt)

	<-done
}
