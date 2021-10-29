package qemu_img

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type QemuImgType string

const (
	QemuImgTypeQcow2 QemuImgType = "qcow2"
)

type Unit struct {
	bytes int64
}

func ParseUnit(v string) (*Unit, error) {
	if strings.HasSuffix(v, "G") {
		n, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		return &Unit{
			bytes: n * 1024 * 1024 * 1024,
		}, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Unit{
		bytes: n,
	}, nil
}

func (u *Unit) Int64() int64 {
	return u.bytes
}

func (u *Unit) String() string {
	if u.bytes%(1024*1024*1024) == 0 {
		return fmt.Sprintf("%dG", u.bytes/(1024*1024*1024))
	}
	return fmt.Sprint(u.bytes)
}

type QemuImg struct {
	Path     string
	Type     QemuImgType
	Size     Unit
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
	size, err := ParseUnit(fmt.Sprint(info.VirtualSize))
	if err != nil {
		return &QemuImg{Path: path, IsExists: false}
	}
	return &QemuImg{
		Size:     *size,
		Type:     QemuImgType(info.Format),
		IsExists: true,
	}
}

func (q *QemuImg) Create(ctx context.Context) error {
	command := "qemu-img"
	args := []string{
		"create", "-f", string(q.Type), q.Path, q.Size.String(),
	}
	cmd := exec.CommandContext(ctx, command, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}
