package core

import (
	"github.com/ophum/humstack/pkg/api/meta"
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
	ExternalIPPoolID string `json:"externalIPPoolID" yaml:"externalIPPoolID"`

	ExternalIPv4Address string `json:"externalIPv4Address" yaml:"externalIPv4Address"`
	ExternalIPv4Prefix  int32  `json:"externalIPv4Prefix" yaml:"externalIPv4Prefix"`
	ExternalIPv6Address string `json:"externalIPv6Address" yaml:"externalIPv6Address"`
	ExternalIPv6Prefix  int32  `json:"externalIPv6Prefix" yaml:"externalIPv6Prefix"`
}
type ExternalIP struct {
	meta.Meta `json:"meta" yaml:"meta"`

	Spec ExternalIPSpec `json:"spec" yaml:"spec"`
}
