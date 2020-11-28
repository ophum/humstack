package blockstorage

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

func (a *BlockStorageAgent) syncCephBlockStorage(bs *system.BlockStorage) error {
	// ex. rbd:pool-name/image-name
	path := filepath.Join(fmt.Sprintf("rbd:%s", a.config.CephPoolName), bs.ID)

	// 削除処理
	if bs.DeleteState == meta.DeleteStateDelete {
		if bs.Status.State != "" &&
			bs.Status.State != system.BlockStorageStateError &&
			bs.Status.State != system.BlockStorageStateActive {
			return nil
		}

		bs.Status.State = system.BlockStorageStateDeleting
		_, err := a.client.SystemV0().BlockStorage().Update(bs)
		if err != nil {
			return err
		}

		// TODO: cephのpoolにイメージがあるかチェックする

		err = a.client.SystemV0().BlockStorage().Delete(bs.Group, bs.Namespace, bs.ID)
		if err != nil {
			return err
		}

		return nil
	}

	// コピー中・ダウンロード中の場合はskip
	switch bs.Status.State {
	case system.BlockStorageStateCopying:
	case system.BlockStorageStateDownloading:
		return nil
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
		// TODO: From HTTP
	case system.BlockStorageFromTypeBaseImage:
		// TODO: From BaseImage
	}
	return nil
}
