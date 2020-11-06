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

type NetworkState string

const (
	NetworkStatePending   NetworkState = "Pending"
	NetworkStateAvailable NetworkState = "Available"
	NetworkStateDeleting  NetworkState = "Deleting"
)

type NetworkStatusLog struct {
	NodeID   string `json:"nodeID" yaml:"nodeID"`
	Datetime string `json:"datetime" yaml:"datetime"`
	Log      string `json:"log" yaml:"log"`
}
type NetworkStatus struct {
	State              NetworkState                 `json:"state" yaml:"state"`
	AttachedInterfaces map[string]VirtualMachineNIC `json:"attachedInterfaces" yaml:"attachedInterfaces"`
	Logs               []NetworkStatusLog           `json:"logs" yaml:"logs"`
}

type Network struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   NetworkSpec   `json:"spec" yaml:"spec"`
	Status NetworkStatus `json:"status" yaml:"status"`
}

type ImageEntitySource struct {
	Namespace      string `json:"namespace" yaml:"namespace"`
	BlockStorageID string `json:"blockStorageID" yaml:"blockStorageID"`
}
type ImageEntitySpec struct {
	Hash   string            `json:"hash" yaml:"hash"`
	Source ImageEntitySource `json:"source" yaml:"source"`
}

type ImageEntityState string

const (
	ImageEntityStatePending   ImageEntityState = "Pending"
	ImageEntityStateCopying   ImageEntityState = "Copying"
	ImageEntityStateAvailable ImageEntityState = "Available"
	ImageEntityStateDeleting  ImageEntityState = "Deleting"
)

type ImageEntityStatus struct {
	State ImageEntityState `json:"state" yaml:"state"`
}
type ImageEntity struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   ImageEntitySpec   `json:"spec" yaml:"spec"`
	Status ImageEntityStatus `json:"status" yaml:"status"`
}

type ImageSpec struct {
	EntityMap map[string]string `json:"entityMap" yaml:"entityMap"`
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
	BlockStorageStateActive      BlockStorageState = "Active"
	BlockStorageStateUsed        BlockStorageState = "Used"
	BlockStorageStatePending     BlockStorageState = "Pending"
	BlockStorageStateDeleting    BlockStorageState = "Deleting"
	BlockStorageStateCopying     BlockStorageState = "Copying"
	BlockStorageStateDownloading BlockStorageState = "Downloading"
	BlockStorageStateReserved    BlockStorageState = "Reserved"
)

type BlockStorageStatus struct {
	State BlockStorageState `json:"state" yaml:"state"`
}
type BlockStorage struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   BlockStorageSpec   `json:"spec" yaml:"spec"`
	Status BlockStorageStatus `json:"status" yaml:"status"`
}

type VirtualMachineLoginUser struct {
	Username          string   `json:"username" yaml:"username"`
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys" yaml:"sshAuthorizedKeys"`
}

type VirtualMachineNIC struct {
	NetworkID      string   `json:"networkID" yaml:"networkID"`
	MacAddress     string   `json:"macAddress" yaml:"macAddress"`
	IPv4Address    string   `json:"ipv4Address" yaml:"ipv4Address"`
	IPv6Address    string   `json:"ipv6Address" yaml:"ipv6Address"`
	Nameservers    []string `json:"nameservers" yaml:"nameservers"`
	DefaultGateway string   `json:"defaultGateway" yaml:"defaultGateway"`
}

type VirtualMachineActionState string

const (
	VirtualMachineActionStatePowerOn  VirtualMachineActionState = "PowerOn"
	VirtualMachineActionStatePowerOff VirtualMachineActionState = "PowerOff"
)

type VirtualMachineSpec struct {
	UUID string `json:"uuid" yaml:"uuid"`

	RequestVcpus string `json:"requestVcpus" yaml:"requestVcpus"`
	LimitVcpus   string `json:"limitVcpus" yaml:"limitVcpus"`

	RequestMemory string `json:"requestMemory" yaml:"requestMemory"`
	LimitMemory   string `json:"limitMemory" yaml:"limitMemory"`

	BlockStorageIDs []string `json:"blockStorageIDs" yaml:"blockStorageIDs"`

	NICs []*VirtualMachineNIC `json:"nics" yaml:"nics"`

	LoginUsers []*VirtualMachineLoginUser `json:"loginUsers" yaml:"loginUsers"`

	ActionState VirtualMachineActionState `json:"actionState" yaml:"actionState"`
}

type VirtualMachineState string

const (
	VirtualMachineStateRunning  VirtualMachineState = "Running"
	VirtualMachineStatePending  VirtualMachineState = "Pending"
	VirtualMachineStateStopping VirtualMachineState = "Stopping"
	VirtualMachineStateStopped  VirtualMachineState = "Stopped"
)

type VirtualMachineStatus struct {
	State VirtualMachineState `json:"state" yaml:"state"`
}
type VirtualMachine struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   VirtualMachineSpec   `json:"spec" yaml:"spec"`
	Status VirtualMachineStatus `json:"status" yaml:"status"`
}

type NodeSpec struct {
	Address     string `json:"address" yaml:"address"`
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

type NATRuleType string

const (
	NATRuleTypeDNAT NATRuleType = "DNAT"
	NATRuleTypeSNAT NATRuleType = "SNAT"
	NATRuleTypeNAPT NATRuleType = "NAPT"
)

type NATRule struct {
	Type            NATRuleType `json:"type" yaml:"rule"`
	SrcNetwork      string      `json:"srcNetwork" yaml:"srcNetwork"`
	DestNetwork     string      `json:"destNetwork" yaml:"destNetwork"`
	InputInterface  string      `json:"inputInterface" yaml:"inputInterface"`
	OutputInterface string      `json:"outputInterface" yaml:"outputInterface"`
}

type DNATRule struct {
	DestAddress   string `json:"destAddress" yaml:"destAddress"`
	DestPort      int32  `json:"destPort" yaml:"destPort"`
	ToDestAddress string `json:"toDestAddress" yaml:"toDestAddress"`
	ToDestPort    int32  `json:"toDestPort" yaml:"toDestPort"`
}

type VirtualRouterExternalIP struct {
	ExternalIPID            string `json:"externalIPID" yaml:"externalIPID"`
	BindInternalIPv4Address string `json:"bindInternalIPv4Address" yaml:"bindInternalIPv4Address"`
}
type VirtualRouterNIC struct {
	NetworkID   string `json:"networkID" yaml:"networkID"`
	IPv4Address string `json:"ipv4Address" yaml:"ipv4Address"`
}
type VirtualRouterSpec struct {
	ExternalGateway string                    `json:"externalGateway" yaml:"externalGateway"`
	ExternalIPs     []VirtualRouterExternalIP `json:"externalIPs" yaml:"externalIPs"`
	NATGatewayIP    string                    `json:"natGatewayIP" yaml:"natGatewayIP"`
	NICs            []VirtualRouterNIC        `json:"nics" yaml:"nics"`
	NATRules        []NATRule                 `json:"natRules" yaml:"natRules"`
	DNATRules       []DNATRule                `json:"dnatRules" yaml:"dnatRules"`
}

type VirtualRouterState string

const (
	VirtualRouterStatePending VirtualRouterState = "Pending"
	VirtualRouterStateRunning VirtualRouterState = "Running"
)

type VirtualRouterStatus struct {
	State VirtualRouterState `json:"state" yaml:"state"`
}

type VirtualRouter struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   VirtualRouterSpec   `json:"spec" yaml:"spec"`
	Status VirtualRouterStatus `json:"status" yaml:"status"`
}
