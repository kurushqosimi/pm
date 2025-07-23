package fsselect

import (
	"github.com/kurushqosimi/pm/pkg/manifest"
	"os"
	"path/filepath"
	"testing"
)

func TestCollect(t *testing.T) {
	tmp := t.TempDir()
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	files := []string{
		"archive_this1/a.txt",
		"archive_this1/b.txt",
		"archive_this1/d.tmp",
		"archive_this2/c.log",
		"archive_this2/b.txt",
		"archive_this2/d.tmp",
	}
	for _, f := range files {
		full := filepath.Join(tmp, f)
		must(os.MkdirAll(filepath.Dir(full), 0o755))
		must(os.WriteFile(full, []byte("x"), 0o644))
	}

	targets := []manifest.Target{
		{Path: "./archive_this1/*.txt"},
		{Path: "./archive_this2/*", Exclude: []string{"*.tmp"}},
	}

	list, err := Collect(tmp, targets)
	if err != nil {
		t.Fatalf("collect err: %v", err)
	}
	exp := []string{
		"archive_this1/a.txt",
		"archive_this1/b.txt",
		"archive_this2/b.txt",
		"archive_this2/c.log",
	}
	if len(list) != len(exp) {
		t.Fatalf("want %v, got %v", exp, list)
	}
	for i := range exp {
		if list[i] != exp[i] {
			t.Fatalf("mismatch at %d: %s vs %s", i, list[i], exp[i])
		}
	}
}
