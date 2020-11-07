package core

import (
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

const (
	ResourceTypeNamespace meta.ResourceType = "Namespace"
)

type NamespaceSpec struct {
}

type Namespace struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec NamespaceSpec `json:"spec" yaml:"spec"`
}

type GroupSpec struct {
}

type Group struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec GroupSpec `json:"spec" yaml:"spec"`
}

type UserSpec struct {
	Password string `json:"password" yaml:"password"`
}
type User struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec UserSpec `json:"spec" yaml:"spec"`
}

type ExternalIPPoolSpec struct {
	IPv4CIDR string `json:"ipv4CIDR" yaml:"ipv4CIDR"`
	IPv6CIDR string `json:"ipv6CIDR" yaml:"ipv6CIDR"`

	BridgeName     string `json:"bridgeName" yaml:"bridgeName"`
	DefaultGateway string `json:"defaultGateway" yaml:"defaultGateway"`
}

type ExternalIPPoolUsed struct {
	UsedExternalIPID string `json:"usedExternalIPID" yaml:"usedExternalIPID"`
}

type ExternalIPPoolStatus struct {
	UsedIPv4Addresses map[string]ExternalIPPoolUsed `json:"usedIPv4Address" yaml:"usedIPv4Address"`
	UsedIPv6Addresses map[string]ExternalIPPoolUsed `json:"usedIPv6Address" yaml:"usedIPv6Address"`
}

type ExternalIPPool struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec   ExternalIPPoolSpec   `json:"spec" yaml:"spec"`
	Status ExternalIPPoolStatus `json:"status" yaml:"status"`
}

type ExternalIPSpec struct {
	PoolID string `json:"poolID" yaml:"poolID"`

	IPv4Address string `json:"ipv4Address" yaml:"ipv4Address"`
	IPv4Prefix  int32  `json:"ipv4Prefix" yaml:"ipv4Prefix"`
	IPv6Address string `json:"ipv6Address" yaml:"ipv6Address"`
	IPv6Prefix  int32  `json:"ipv6Prefix" yaml:"ipv6Prefix"`
}
type ExternalIP struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec ExternalIPSpec `json:"spec" yaml:"spec"`
}

type NetworkSpec struct {
	Template system.NodeNetwork `json:"template" yaml:"template"`
}

type NetworkState string

const (
	NetworkStateActive   NetworkState = "Active"
	NetworkStateCreating NetworkState = "Creating"
)

type NetworkStatus struct {
	State NetworkState `json:"state" yaml:"state"`
}
type Network struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec NetworkSpec `json:"spec" yaml:"spec"`

	Status NetworkStatus `json:"status" yaml:"status"`
}
