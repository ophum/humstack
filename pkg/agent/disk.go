package agent

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/ophum/humstack/v1/pkg/agent/driver/qemu_img"
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/client"
)

type DiskAgentConfig struct {
	HostBasePath string `yaml:"hostBasePath"`
}

type DiskAgent struct {
	diskClient client.IDiskClient
	config     *DiskAgentConfig
}

func NewDiskAgent(
	diskClient client.IDiskClient,
	config *DiskAgentConfig,
) *DiskAgent {
	return &DiskAgent{diskClient, config}
}

func (a *DiskAgent) Start(ctx context.Context) {
	t := time.NewTicker(time.Second * 2)
	a.process(ctx)

	for {
		select {
		case <-t.C:
			a.process(ctx)
		case <-ctx.Done():
			return
		}

	}
}

func (a *DiskAgent) getPath(disk *entity.Disk) string {
	return filepath.Join(a.config.HostBasePath, disk.Name)
}

func (a *DiskAgent) process(ctx context.Context) {
	log.Println("=============START=================")
	disks, err := a.diskClient.List(ctx)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, disk := range disks {
		log.Println(disk.Name)
		path := a.getPath(disk)
		img := qemu_img.Open(ctx, path)
		if img.IsExists {
			log.Println("already created, ", img)
			if disk.Status == entity.DiskStatusPending || disk.Status == entity.DiskStatusProcessing {
				a.diskClient.UpdateStatus(ctx, disk.Name, entity.DiskStatusActive)
			}
			continue
		}

		if err := a.create(ctx, disk); err != nil {
			log.Println("failed to create disk: ", err)
			continue
		}
		a.diskClient.UpdateStatus(ctx, disk.Name, entity.DiskStatusActive)
	}
	log.Println("=============FINISH=================")
}

func (a *DiskAgent) create(ctx context.Context, disk *entity.Disk) error {
	size, err := qemu_img.ParseUnit(fmt.Sprint(disk.LimitBytes))
	if err != nil {
		return err
	}
	path := a.getPath(disk)
	img := &qemu_img.QemuImg{
		Path: path,
		Type: qemu_img.QemuImgTypeQcow2,
		Size: *size,
	}
	if err := img.Create(ctx); err != nil {
		return err
	}
	return nil
}
