package internal

import (
	"bytes"
	"fmt"
	"github.com/kurushqosimi/pm/pkg/archiver"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"os"
	"path"
)

func DownloadAndExtract(name, ver, repoRoot, localDir string, remote sshclient.Storage) error {
	remotePath := path.Join(repoRoot, name, ver+".tar.gz")

	data, err := remote.Download(remotePath)
	if err != nil {
		return fmt.Errorf("download %s: %w", remotePath, err)
	}

	dst := path.Join(localDir, name, ver)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	r := bytes.NewReader(data)
	if err := archiver.ExtractTarGz(r, dst); err != nil {
		return fmt.Errorf("untar: %w", err)
	}
	return nil
}

func ExistsLocally(name, ver, localDir string) bool {
	_, err := os.Stat(path.Join(localDir, name, ver))
	return err == nil
}
