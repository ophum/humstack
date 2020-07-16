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
