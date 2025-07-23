package manifest

import (
	"os"
	"testing"
)

func TestReadManifest_JSON(t *testing.T) {
	m, err := Read("testdata/packet.json")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "packet-1" || m.Ver != "1.10" {
		t.Fatalf("unexpected content %#v", m)
	}
}

func TestReadManifest_YAML(t *testing.T) {
	m, err := Read("testdata/packet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Targets) == 0 {
		t.Fatalf("targets not parsed")
	}
}

func TestReadManifest_FileNotFound(t *testing.T) {
	_, err := Read("testdata/not_exist.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadManifest_UnsupportedExtension(t *testing.T) {
	f := "testdata/packet.txt"
	err := os.WriteFile(f, []byte(`dummy`), 0644)
	defer os.Remove(f)

	_, err = Read(f)
	if err == nil || err.Error() != `unsupported manifest extension ".txt"` {
		t.Fatalf("expected unsupported extension error, got: %v", err)
	}
}

func TestReadManifest_InvalidJSON(t *testing.T) {
	f := "testdata/invalid.json"
	err := os.WriteFile(f, []byte(`{invalid json}`), 0644)
	defer os.Remove(f)

	_, err = Read(f)
	if err == nil || err.Error()[:5] != "json:" {
		t.Fatalf("expected json error, got: %v", err)
	}
}

func TestReadManifest_EmptyFile(t *testing.T) {
	f := "testdata/empty.json"
	err := os.WriteFile(f, []byte(``), 0644)
	defer os.Remove(f)

	_, err = Read(f)
	if err == nil {
		t.Fatal("expected error for empty file, got nil")
	}
}

func TestReadManifest_JSON_Fields(t *testing.T) {
	m, err := Read("testdata/packet.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Targets) != 2 || len(m.Targets[1].Exclude) != 1 {
		t.Fatal("exclude field not parsed correctly")
	}
	if len(m.Packets) != 1 || m.Packets[0].Name != "packet-3" {
		t.Fatal("packets field not parsed correctly")
	}
}
