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
	"github.com/pkg/errors"
)

func (a *BlockStorageAgent) syncLocalBlockStorage(bs *system.BlockStorage) error {
	dirPath := filepath.Join(a.localBlockStorageDirectory, bs.Group, bs.Namespace)
	path := filepath.Join(dirPath, bs.ID)

	// コピー中・ダウンロード中の場合はskip
	switch bs.Status.State {
	case system.BlockStorageStateCopying, system.BlockStorageStateDownloading:
		log.Println("skip copying or downloading")
		return nil
	}

	// 削除処理
	if bs.DeleteState == meta.DeleteStateDelete {
		if bs.Status.State != "" &&
			bs.Status.State != system.BlockStorageStateError &&
			bs.Status.State != system.BlockStorageStateActive &&
			bs.Status.State != system.BlockStorageStateQueued {
			return nil
		}
		bs.Status.State = system.BlockStorageStateDeleting
		_, err := a.client.SystemV0().BlockStorage().Update(bs)
		if err != nil {
			return err
		}

		if fileIsExists(path) {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		err = a.client.SystemV0().BlockStorage().Delete(bs.Group, bs.Namespace, bs.ID)
		if err != nil {
			return err
		}

		return nil
	}

	if fileIsExists(path) {
		switch bs.Status.State {
		case system.BlockStorageStateError:
			// Stateがエラーなら存在するイメージを消す
			if err := os.Remove(path); err != nil {
				return err
			}
		case "", system.BlockStorageStatePending:
			bs.Status.State = system.BlockStorageStateActive
			return setHash(bs)
		case system.BlockStorageStateActive, system.BlockStorageStateUsed:
			return nil
		}
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
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}
	case system.BlockStorageFromTypeHTTP:
		bs.Status.State = system.BlockStorageStateDownloading
		if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
			return err
		}

		res, err := http.Get(bs.Spec.From.HTTP.URL)
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}
		defer res.Body.Close()

		file, err := os.Create(path)
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}

		//reader := shapeio.NewReader(res.Body)
		//reader.SetRateLimit(1024 * 1024 * 10) // 10KB/sec
		if res.ContentLength >= 0 {
			_, err = io.CopyN(file, res.Body, res.ContentLength)
		} else {
			_, err = io.Copy(file, res.Body)
		}
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}

		err = file.Close()
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}

		command := "qemu-img"
		args := []string{
			"resize",
			path,
			withUnitToWithoutUnit(bs.Spec.LimitSize),
		}
		cmd := exec.Command(command, args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return errors.Wrap(err, string(out))
		}
	case system.BlockStorageFromTypeBaseImage:

		bs.Status.State = system.BlockStorageStateCopying
		if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
			return err
		}

		image, err := a.client.SystemV0().Image().Get(bs.Group, bs.Spec.From.BaseImage.ImageName)
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}

		imageEntity, ok := image.Spec.EntityMap[bs.Spec.From.BaseImage.Tag]
		if !ok {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return fmt.Errorf("Image Entity not found")
		}

		srcDirPath := filepath.Join(a.localImageDirectory, bs.Group)
		if !fileIsExists(srcDirPath) {
			err := os.MkdirAll(srcDirPath, 0755)
			if err != nil {
				bs.Status.State = system.BlockStorageStateError
				if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
					return err
				}
				return err
			}
		}
		srcPath := filepath.Join(srcDirPath, imageEntity)

		// なかったら別のノードから持ってくるようにする
		// 動くはず
		if !fileIsExists(srcPath) {
			err := func() error {
				src, err := os.Create(srcPath)
				if err != nil {
					bs.Status.State = system.BlockStorageStateError
					if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
						return err
					}
					return err
				}
				defer src.Close()

				stream, _, err := a.client.SystemV0().Image().Download(bs.Group, bs.Spec.From.BaseImage.ImageName, bs.Spec.From.BaseImage.Tag)
				if err != nil {
					bs.Status.State = system.BlockStorageStateError
					if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
						return err
					}
					return err
				}
				defer stream.Close()

				if _, err := io.Copy(src, stream); err != nil {
					bs.Status.State = system.BlockStorageStateError
					if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
						return err
					}
					return err
				}
				return nil
			}()
			if err != nil {
				bs.Status.State = system.BlockStorageStateError
				if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
					return err
				}
				return err
			}
		}
		src, err := os.Open(srcPath)
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}
		defer src.Close()

		dest, err := os.Create(path)
		if err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
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
			bs.Status.State = system.BlockStorageStateError
			if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
				return err
			}
			return err
		}
	}

	if bs.Annotations == nil {
		bs.Annotations = map[string]string{}
	}

	bs.Annotations["bs-download-host"] = fmt.Sprintf("%s:%d", a.config.DownloadAPI.AdvertiseAddress, a.config.DownloadAPI.ListenPort)
	// ここに来た時点で処理は終わっているのでActiveにする
	if bs.Status.State == "" ||
		bs.Status.State == system.BlockStorageStatePending ||
		bs.Status.State == system.BlockStorageStateCopying ||
		bs.Status.State == system.BlockStorageStateDownloading {
		bs.Status.State = system.BlockStorageStateActive

		if _, err := a.client.SystemV0().BlockStorage().Update(bs); err != nil {
			return err
		}
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
