package system

import (
	"github.com/ophum/humstack/pkg/api/meta"
)

const (
	ResourceTypeNetwork meta.ResourceType = "Network"
)

type NetworkSpec struct {
	ID       string `json:"id" yaml:"id"`
	IPv4CIDR string `json:"ipv4CIDR" yaml:"ipv4CIDR"`
	IPv6CIDR string `json:"ipv6CIDR" yaml:"ipv6CIDR"`
}

type Network struct {
	meta.Meta

	Spec NetworkSpec `json:"spec" yaml:"spec"`
}

type BlockStorageSpec struct {
	RequestSize string `json:"requestSize" yaml:"requestSize"`
	LimitSize   string `json:"limitSize" yaml:"limitSize"`
}
type BlockStorage struct {
	meta.Meta

	Spec BlockStorageSpec `json:"spec" yaml:"spec"`
}

type VirtualMachineLoginUser struct {
	Username          string   `json:"username" yaml:"username"`
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys" yaml:"sshAuthorizedKeys"`
}

type VirtualMachineNIC struct {
	NetworkName    string `json:"networkName" yaml:"networkName"`
	IPv4Address    string `json:"ipv4Address" yaml:"ipv4Address"`
	IPv6Address    string `json:"ipv6Address" yaml:"ipv6Address"`
	DefaultGateway string `json:"defaultGateway" yaml:"defaultGateway"`
}

type VirtualMachineSpec struct {
	RequestVcpus string `json:"requestVcpus" yaml:"requestVcpus"`
	LimitVcpus   string `json:"limitVcpus" yaml:"limitVcpus"`

	RequestMemory string `json:"requestMemory" yaml:"requestMemory"`
	LimitMemory   string `json:"limitMemory" yaml:"limitMemory"`

	BlockStorageNames []string `json:"blockStorageNames" yaml:"blockStorageNames"`

	NICs []*VirtualMachineNIC `json:"nics" yaml:"nics"`

	LoginUsers []*VirtualMachineLoginUser `json:"loginUsers"`
}

type VirtualMachine struct {
	meta.Meta

	Spec VirtualMachineSpec `json:"spec" yaml:"spec"`
}

type NodeSpec struct {
	LimitVcpus  string `json:"limitVcpus" yaml:"limitVcpus"`
	LimitMemory string `json:"limitMemory" yaml:"limitMemory"`
}

type NodeState string

const (
	NodeStateNotReady NodeState = "NotReady"
	NodeStateReady    NodeState = "Ready"
)

type NodeStatus struct {
	State           NodeState `json:"state" yaml:"state"`
	RequestedVcpus  string    `json:"requestedVcpus" yaml:"requestedVcpus"`
	RequestedMemory string    `json:"requestedMemory" yaml:"requestedMemory"`
}

type Node struct {
	meta.Meta

	Spec   NodeSpec   `json:"spec" yaml:"spec"`
	Status NodeStatus `json:"status" yaml:"status"`
}
