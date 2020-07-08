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
	meta.Meta `json:"meta" yaml:"meta"`

	Spec NetworkSpec `json:"spec" yaml:"spec"`
}

type ImageTag struct {
	Tag              string `json:""`
	BlockStorageName string `json:"blockStorageName" yaml:"blockStorageName"`
}
type ImageSpec struct {
	Tags []ImageTag `json:"tags" yaml:"tags"`
}

type Image struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec ImageSpec `json:"spec" yaml:"spec"`
}

type BlockStorageFromBaseImage struct {
	ImageName string `json:"imageName" yaml:"imageName"`
	Tag       string `json:"tag" yaml:"tag"`
}

type BlockStorageFromBlockStorage struct {
	Name string `json:"name" yaml:"name"`
}

type BlockStorageFromHTTP struct {
	URL string `json:"url" yaml:"url"`
}

type BlockStorageFromType string

const (
	BlockStorageFromTypeBaseImage    BlockStorageFromType = "BaseImage"
	BlockStorageFromTypeBlockStorage BlockStorageFromType = "BlockStorage"
	BlockStorageFromTypeHTTP         BlockStorageFromType = "HTTP"
	BlockStorageFromTypeEmpty        BlockStorageFromType = "Empty"
)

type BlockStorageFrom struct {
	Type         BlockStorageFromType         `json:"type" yaml:"type"`
	BaseImage    BlockStorageFromBaseImage    `json:"baseImage" yaml:"baseImage"`
	BlockStorage BlockStorageFromBlockStorage `json:"blockStorage" yaml:"blockStorage"`
	HTTP         BlockStorageFromHTTP         `json:"http" yaml:"http"`
}
type BlockStorageSpec struct {
	RequestSize string           `json:"requestSize" yaml:"requestSize"`
	LimitSize   string           `json:"limitSize" yaml:"limitSize"`
	From        BlockStorageFrom `json:"from" yaml:"from"`
}

type BlockStorageState string

const (
	BlockStorageStateActive  BlockStorageState = "Active"
	BlockStorageStateUsed    BlockStorageState = "Used"
	BlockStorageStatePending BlockStorageState = "Pending"
)

type BlockStorageStatus struct {
	State BlockStorageState `json:"state" yaml:"state"`
}
type BlockStorage struct {
	meta.Meta

	Spec   BlockStorageSpec   `json:"spec" yaml:"spec"`
	Status BlockStorageStatus `json:"status" yaml:"status"`
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

type VirtualMachineActionState string

const (
	VirtualMachineActionStateStart  VirtualMachineActionState = "Start"
	VirtualMachineActionStateStop   VirtualMachineActionState = "Stop"
	VirtualMachineActionStateReboot VirtualMachineActionState = "Reboot"
	VirtualMachineActionStateDone   VirtualMachineActionState = "Done"
)

type VirtualMachineSpec struct {
	RequestVcpus string `json:"requestVcpus" yaml:"requestVcpus"`
	LimitVcpus   string `json:"limitVcpus" yaml:"limitVcpus"`

	RequestMemory string `json:"requestMemory" yaml:"requestMemory"`
	LimitMemory   string `json:"limitMemory" yaml:"limitMemory"`

	BlockStorageNames []string `json:"blockStorageNames" yaml:"blockStorageNames"`

	NICs []*VirtualMachineNIC `json:"nics" yaml:"nics"`

	LoginUsers []*VirtualMachineLoginUser `json:"loginUsers"`

	ActionState VirtualMachineActionState `json:"actionState" yaml:"actionState"`
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
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   NodeSpec   `json:"spec" yaml:"spec"`
	Status NodeStatus `json:"status" yaml:"status"`
}
