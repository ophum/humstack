package entity

import "time"

type NodeStatus string

const (
	NodeStatusNotReady NodeStatus = "NotReady"
	NodeStatusReady    NodeStatus = "Ready"
)

type Node struct {
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`

	Hostname string `json:"hostname"`

	Agents []*AgentDaemon `json:"agents"`
	Status NodeStatus     `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgentDaemonRestartPolicy string

const (
	AgentDaemonRestartPolicyAlways    AgentDaemonRestartPolicy = "Always"
	AgentDaemonRestartPolicyOnFailure AgentDaemonRestartPolicy = "OnFailure"
	AgentDaemonRestartPolicyNever     AgentDaemonRestartPolicy = "Never"
)

type AgentDaemon struct {
	Name          string                   `json:"name"`
	Command       string                   `json:"command"`
	Args          []string                 `json:"args"`
	Envs          map[string]string        `json:"envs"`
	RestartPolicy AgentDaemonRestartPolicy `json:"restart_policy"`
}

type AgentDaemonStatus string

const (
	AgentDaemonStatusPending AgentDaemonStatus = "Pending"
	AgentDaemonStatusRunning AgentDaemonStatus = "Running"
	AgentDaemonStatusError   AgentDaemonStatus = "Error"
)
