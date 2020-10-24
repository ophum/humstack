package blockstorage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

func (a *BlockStorageAgent) syncLocalBlockStorage(bs *system.BlockStorage) error {
	dirPath := filepath.Join(a.localBlockStorageDirectory, bs.Group, bs.Namespace)
	path := filepath.Join(dirPath, bs.ID)
	log.Printf("[BS] %s\n", bs.Name)
	log.Printf("[BS] ==> %s", path)

	if fileIsExists(path) {
		// 削除処理
		if bs.DeleteState == meta.DeleteStateDelete {
			log.Println("[BS] ==> DELETING")
			bs.Status.State = system.BlockStorageStateDeleting
			_, err := a.client.SystemV0().BlockStorage().Update(bs)
			if err != nil {
				return err
			}

			err = os.Remove(path)
			if err != nil {
				return err
			}

			err = a.client.SystemV0().BlockStorage().Delete(bs.Group, bs.Namespace, bs.ID)
			if err != nil {
				return err
			}

			log.Println("[BS] ====> DELETED")
			return nil
		}
		if bs.Status.State == "" || bs.Status.State == system.BlockStorageStatePending {
			bs.Status.State = system.BlockStorageStateActive
		}
		log.Println("[BS] ====> ALREADY ACTIVE")

		return setHash(bs)
	}

	if !fileIsExists(dirPath) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	switch bs.Spec.From.Type {
	case system.BlockStorageFromTypeEmpty:
		log.Println("[BS] ====> CREATE EMPTY IMAGE")
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
	case system.BlockStorageFromTypeHTTP:
		log.Printf("[BS] ====> DOWNLOAD: %s\n", bs.Spec.From.HTTP.URL)
		res, err := http.Get(bs.Spec.From.HTTP.URL)
		if err != nil {
			log.Println(err)
			return err
		}
		defer res.Body.Close()

		file, err := os.Create(path)
		if err != nil {
			log.Println(err)
			return err
		}

		_, err = io.Copy(file, res.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		err = file.Close()
		if err != nil {
			log.Println(err)
			return err
		}

		command := "qemu-img"
		args := []string{
			"resize",
			path,
			withUnitToWithoutUnit(bs.Spec.LimitSize),
		}
		cmd := exec.Command(command, args...)
		if _, err := cmd.CombinedOutput(); err != nil {
			log.Println(err.Error())
			return err
		}
	case system.BlockStorageFromTypeBaseImage:
		log.Printf("[BS] Copy from base image.")
		image, err := a.client.SystemV0().Image().Get(bs.Group, bs.Spec.From.BaseImage.ImageName)
		if err != nil {
			return err
		}

		log.Println(image)

		imageEntity, ok := image.Spec.EntityMap[bs.Spec.From.BaseImage.Tag]
		if !ok {
			return fmt.Errorf("Image Entity not found")
		}

		srcPath := filepath.Join(a.localImageDirectory, bs.Group, imageEntity)

		// TODO: なかったら別のノードから持ってくるようにする
		src, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.Create(path)
		if err != nil {
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			return err
		}

		command := "qemu-img"
		args := []string{
			"resize",
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
