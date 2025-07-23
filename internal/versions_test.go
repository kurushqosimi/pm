package internal

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/kurushqosimi/pm/internal/mockstorage"
	"os"
	"path/filepath"
	"testing"
)

func TestPickBest(t *testing.T) {
	tmp := t.TempDir()
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	write := func(name string) {
		full := filepath.Join(tmp, name)
		must(os.MkdirAll(filepath.Dir(full), 0o755))
		must(os.WriteFile(full, []byte("x"), 0o644))
	}
	write("1.0.tar.gz")
	write("1.10.tar.gz")
	write("1.12.tar.gz")
	write("nonsense.txt")

	remote := mockstorage.FS{Root: tmp}
	c, _ := semver.NewConstraint(">=1.10")

	best, ok, err := BestVersion("", c, remote)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("parsed version:", ok)

	if !ok || best != "1.12.0" {
		t.Fatalf("want 1.12.0, got %q ok=%v", best, ok)
	}
}
