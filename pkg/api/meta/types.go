package meta

type APIType string

const (
	APITypeNodeV0           APIType = "systemv0/node"
	APITypeNetworkV0        APIType = "systemv0/network"
	APITypeBlockStorageV0   APIType = "systemv0/blockstorage"
	APITypeVirtualMachineV0 APIType = "systemv0/virtualmachine"
	APITypeNamespaceV0      APIType = "corev0/namespace"
	APITypeGroupV0          APIType = "corev0/group"
)

type ResourceType string

type DeleteState string

const (
	DeleteStateNone   DeleteState = "None"
	DeleteStateDelete DeleteState = "Delete"
	DeleteStateDone   DeleteState = "Done"
)

type Meta struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Namespace    string            `json:"namespace" yaml:"namespace"`
	Group        string            `json:"group" yaml:"group"`
	Annotations  map[string]string `json:"annotations" yaml:"annotations"`
	Labels       map[string]string `json:"labels" yaml:"labels"`
	ResourceHash string            `json:"resourceHash" yaml:"resourceHash"`
	DeleteState  DeleteState       `json:"deleteState" yaml:"deleteState"`
	APIType      APIType           `json:"apiType" yaml:"apiType"`
}
