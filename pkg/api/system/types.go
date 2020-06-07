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
