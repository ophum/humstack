package qemu_img

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type QemuImgType string

const (
	QemuImgTypeQcow2 QemuImgType = "qcow2"
)

type QemuImg struct {
	Path     string
	Type     QemuImgType
	Size     int64
	IsExists bool
}

type QemuImgInfo struct {
	VirtualSize int64  `json:"virtual-size"`
	Format      string `json:"format"`
}

func Open(ctx context.Context, path string) *QemuImg {
	command := "qemu-img"
	args := []string{
		"info", path, "--output", "json",
	}
	cmd := exec.CommandContext(ctx, command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &QemuImg{Path: path, IsExists: false}
	}
	var info QemuImgInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return &QemuImg{Path: path, IsExists: false}
	}
	return &QemuImg{
		Size:     info.VirtualSize,
		Type:     QemuImgType(info.Format),
		IsExists: true,
	}
}

func (q *QemuImg) Create(ctx context.Context) error {
	command := "qemu-img"
	args := []string{
		"create", "-f", string(q.Type), q.Path, fmt.Sprint(q.Size),
	}
	cmd := exec.CommandContext(ctx, command, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}
