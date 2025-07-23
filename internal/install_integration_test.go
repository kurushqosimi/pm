package internal

import (
	"github.com/kurushqosimi/pm/internal/mockstorage"
	"github.com/kurushqosimi/pm/pkg/archiver"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadAndExtract(t *testing.T) {
	tmp := t.TempDir()

	repo := mockstorage.FS{Root: tmp}
	name, ver := "packet-1", "1.0.0"

	files := []string{"foo.txt"}
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(tmp, files[0]), []byte("ok"), 0o644))
	tgz, err := archiver.TarGz(files, tmp)
	must(err)

	remotePath := filepath.Join("/repo", name, ver+".tar.gz")
	must(repo.Upload(remotePath, tgz))

	localDir := filepath.Join(tmp, "local")
	if err := DownloadAndExtract(name, ver, "/repo", localDir, repo); err != nil {
		t.Fatalf("install err: %v", err)
	}

	if _, err := os.Stat(filepath.Join(localDir, name, ver, "foo.txt")); err != nil {
		t.Fatalf("file not extracted: %v", err)
	}
}
