package schema

type DiskAnnotation struct {
	DiskName string
	Key      string
	Value    string
}

func (v DiskAnnotation) TableName() string {
	return "disk_annotations"
}

func ToEntityDiskAnnotations(v []*DiskAnnotation) map[string]string {
	annotations := map[string]string{}
	for _, vv := range v {
		annotations[vv.Key] = vv.Value
	}
	return annotations
}
