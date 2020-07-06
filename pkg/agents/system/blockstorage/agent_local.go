package blockstorage

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/ophum/humstack/pkg/api/system"
)

func (a *BlockStorageAgent) syncLocalBlockStorage(bs *system.BlockStorage) error {
	dirPath := filepath.Join(a.localBlockStorageDirectory, bs.Namespace)
	path := filepath.Join(dirPath, bs.ID)
	log.Printf("BLOCKSTORAGE: %s\n", bs.Name)
	log.Printf("==> %s", path)

	if fileIsExists(path) {
		return nil
	}

	if !fileIsExists(dirPath) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	switch bs.Spec.From.Type {
	case system.BlockStorageFromTypeEmpty:
		command := "qemu-img"
		args := []string{
			"create",
			"-f",
			"qcow2",
			path,
			withUnitToWithoutUnit(bs.Spec.LimitSize),
		}
		cmd := exec.Command(command, args...)
		if _, err := cmd.CombinedOutput(); err != nil {
			log.Println(err.Error())
			return err
		}

	}

	if bs.Status.State == "" || bs.Status.State == system.BlockStorageStatePending {
		bs.Status.State = system.BlockStorageStateActive
	}
	return setHash(bs)
}

func fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

const (
	UnitGigabyte = 'G'
	UnitMegabyte = 'M'
	UnitKilobyte = 'K'
)

func withUnitToWithoutUnit(numberWithUnit string) string {
	length := len(numberWithUnit)
	if numberWithUnit[length-1] >= '0' && numberWithUnit[length-1] <= '9' {
		return numberWithUnit
	}

	number, err := strconv.ParseInt(numberWithUnit[:length-1], 10, 64)
	if err != nil {
		return "0"
	}
	switch numberWithUnit[length-1] {
	case UnitGigabyte:
		return fmt.Sprintf("%d", number*1024*1024*1024)
	case UnitMegabyte:
		return fmt.Sprintf("%d", number*1024*1024)
	case UnitKilobyte:
		return fmt.Sprintf("%d", number*1024)
	}
	return "0"
}
