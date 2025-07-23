package internal

import (
	"github.com/kurushqosimi/pm/internal/mockstorage"
	"path/filepath"
	"testing"
)

func TestCreateEnd2End(t *testing.T) {
	tmp := t.TempDir()

	stub := mockstorage.FS{Root: tmp}
	manifest := filepath.Join("..", "pkg", "manifest", "testdata", "packet.json")

	remotePath, err := Create(manifest, "/repo", stub)
	if err != nil {
		t.Fatalf("create err: %v", err)
	}

	if _, err := stub.Download(remotePath); err != nil {
		t.Fatalf("archive not uploaded: %v", err)
	}
}
