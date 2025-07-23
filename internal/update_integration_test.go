package internal

import (
	"github.com/kurushqosimi/pm/internal/mockstorage"
	"github.com/kurushqosimi/pm/pkg/archiver"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateEnd2End(t *testing.T) {
	tmp := t.TempDir()
	repoRoot := "/repo"
	localDir := filepath.Join(tmp, "local")

	ms := mockstorage.FS{Root: tmp}
	makePkg := func(name, ver string) {
		data, _ := archiver.TarGz([]string{}, tmp)
		_ = ms.Upload(filepath.Join(repoRoot, name, ver+".tar.gz"), data)
	}
	makePkg("packet-1", "1.11.0")
	makePkg("packet-2", "0.5.0")

	pkgs := `{"packages":[{"name":"packet-1","ver":">=1.10"},{"name":"packet-2"}]}`
	pkgFile := filepath.Join(tmp, "packages.json")
	_ = os.WriteFile(pkgFile, []byte(pkgs), 0o644)

	if err := Update(pkgFile, repoRoot, localDir, ms); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if !ExistsLocally("packet-1", "1.11.0", localDir) ||
		!ExistsLocally("packet-2", "0.5.0", localDir) {
		t.Fatalf("packages not installed")
	}
}
