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
	meta.Meta

	Spec NamespaceSpec `json:"spec" yaml:"spec"`
}
