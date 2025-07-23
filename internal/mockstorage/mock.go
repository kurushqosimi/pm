package mockstorage

import (
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"os"
	"path/filepath"
)

type FS struct{ Root string }

func (f FS) Upload(dst string, data []byte) error {
	full := filepath.Join(f.Root, dst)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, data, 0o644)
}
func (f FS) Download(src string) ([]byte, error) {
	return os.ReadFile(filepath.Join(f.Root, src))
}

func (f FS) List(dirPath string) ([]sshclient.DirEntry, error) {
	return nil, nil
}
