package schema

import (
	"encoding/json"
	"log"

	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type NodeAgentDaemon struct {
	NodeName      string
	Name          string
	Command       string
	Args          string
	Envs          string
	RestartPolicy string
}

func (v NodeAgentDaemon) TableName() string {
	return "node_agent_daemons"
}

func ToEntityAgentDaemon(v *NodeAgentDaemon) *entity.AgentDaemon {
	if v == nil {
		return nil
	}
	var args []string
	if err := json.Unmarshal([]byte(v.Args), &args); err != nil {
		log.Println("warning: ", err)
	}
	var envs map[string]string
	if err := json.Unmarshal([]byte(v.Envs), &envs); err != nil {
		log.Println("warning: ", err)
	}
	return &entity.AgentDaemon{
		Name:          v.Name,
		Command:       v.Command,
		Args:          args,
		Envs:          envs,
		RestartPolicy: entity.AgentDaemonRestartPolicy(v.RestartPolicy),
	}
}

func ToEntityAgentDaemons(v []*NodeAgentDaemon) []*entity.AgentDaemon {
	l := []*entity.AgentDaemon{}
	for _, vv := range v {
		if vv == nil {
			continue
		}
		l = append(l, ToEntityAgentDaemon(vv))
	}
	return l
}
