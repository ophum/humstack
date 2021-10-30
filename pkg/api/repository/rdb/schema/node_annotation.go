package schema

type NodeAnnotation struct {
	NodeName string
	Key      string
	Value    string
}

func (v NodeAnnotation) TableName() string {
	return "node_annotations"
}

func ToMapNodeAnnotations(v []*NodeAnnotation) map[string]string {
	annotations := map[string]string{}
	for _, vv := range v {
		annotations[vv.Key] = vv.Value
	}
	return annotations
}
