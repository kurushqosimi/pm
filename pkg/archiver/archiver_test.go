package archiver

import (
	"bytes"
	"github.com/bmatcuk/doublestar/v4"
	"os"
	"path/filepath"
	"testing"
)

func TestTarGz(t *testing.T) {
	tmp := t.TempDir()

	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	write := func(rel, data string) {
		full := filepath.Join(tmp, rel)
		must(os.MkdirAll(filepath.Dir(full), 0755))
		must(os.WriteFile(full, []byte(data), 0644))
	}

	write("a.txt", "hello")
	write("dir/b.log", "bye")

	files, _ := doublestar.Glob(os.DirFS(tmp), "**/*")
	var relFiles []string
	for _, p := range files {
		if fi, _ := os.Stat(filepath.Join(tmp, p)); fi.Mode().IsRegular() {
			relFiles = append(relFiles, p)
		}
	}

	data, err := TarGz(relFiles, tmp)
	if err != nil {
		t.Fatalf("tar err: %v", err)
	}

	out := filepath.Join(tmp, "out")
	must(os.MkdirAll(out, 0755))
	if err := ExtractTarGz(bytes.NewBuffer(data), out); err != nil {
		t.Fatalf("untar err: %v", err)
	}

	got, _ := os.ReadFile(filepath.Join(out, "a.txt"))
	if string(got) != "hello" {
		t.Fatalf("unexpected content: %s", got)
	}
}
