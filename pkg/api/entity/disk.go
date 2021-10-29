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

	RequestBytes int `json:"request_bytes"`
	LimitBytes   int `json:"limit_bytes"`

	Status    DiskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
