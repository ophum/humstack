package meta

type APIType string

const (
	APITypeNodeV0           APIType = "systemv0/node"
	APITypeNetworkV0        APIType = "systemv0/network"
	APITypeBlockStorageV0   APIType = "systemv0/blockstorage"
	APITypeVirtualMachineV0 APIType = "systemv0/virtualmachine"
	APITypeNamespaceV0      APIType = "corev0/namespace"
)

type ResourceType string

type Meta struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Namespace    string            `json:"namespace" yaml:"namespace"`
	Annotations  map[string]string `json:"annotations" yaml:"annotations"`
	Labels       map[string]string `json:"labels" yaml:"labels"`
	ResourceHash string            `json:"resourceHash" yaml:"resourceHash"`
	APIType      APIType           `json:"apiType" yaml:"apiType"`
}
