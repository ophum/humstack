package agent

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/client"
)

type NodeAgent struct {
	nodeClient client.INodeClient
}

var processStatusMap = map[string]string{}

func NewNodeAgent(
	nodeClient client.INodeClient,
) *NodeAgent {
	return &NodeAgent{nodeClient}
}

func (a *NodeAgent) Start(ctx context.Context) {
	t := time.NewTicker(time.Second * 2)
	a.process(ctx)

	for {
		select {
		case <-t.C:
			a.process(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (a *NodeAgent) process(ctx context.Context) {
	log.Println("=============START=================")
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		return
	}
	node, err := a.nodeClient.Get(ctx, hostname)
	if err != nil {
		log.Println(err)
		return
	}

	isReady := true
	for _, agent := range node.Agents {
		if _, processed := processStatusMap[agent.Name]; processed {
			log.Printf("%s already started\n", agent.Name)
			continue
		}
		isReady = false
		processStatusMap[agent.Name] = "started"
		go func(ctx context.Context, agent *entity.AgentDaemon) {
			log.Printf("%s starting...\n", agent.Name)
			cmd := exec.CommandContext(ctx, agent.Command, agent.Args...)
			cmd.Stdout = os.Stdout
			cmd.Start()
			cmd.Wait()
			log.Printf("%s exited...\n", agent.Name)
		}(ctx, agent)
	}

	if node.Status == entity.NodeStatusNotReady && isReady {
		a.nodeClient.UpdateStatus(ctx, node.Name, entity.NodeStatusReady)
	} else if node.Status == entity.NodeStatusReady && !isReady {
		a.nodeClient.UpdateStatus(ctx, node.Name, entity.NodeStatusNotReady)
	}
	log.Println("=============FINISH=================")
}
