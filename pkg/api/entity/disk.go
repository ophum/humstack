package entity

import "time"

type DiskType string

const (
	DiskTypeHostQcow2 DiskType = "HostQcow2"
)

type DiskStatus string

const (
	DiskStatusPending    DiskStatus = "Pending"
	DiskStatusProcessing DiskStatus = "Processing"
	DiskStatusActive     DiskStatus = "Active"
	DiskStatusUsed       DiskStatus = "Used"
	DiskStatusDelete     DiskStatus = "Delete"
	DiskStatusDeleting   DiskStatus = "Deleting"
)

type Disk struct {
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`

	Type DiskType `json:"type"`

	RequestSize ByteUnit `json:"request_size"`
	LimitSize   ByteUnit `json:"limit_size"`

	Status    DiskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
