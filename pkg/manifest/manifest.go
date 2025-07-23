package manifest

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type Target struct {
	Path    string   `json:"path" yaml:"path"`
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

type Dep struct {
	Name string `json:"name" yaml:"name"`
	Ver  string `json:"ver,omitempty" yaml:"ver,omitempty"`
}

type Manifest struct {
	Name    string   `json:"name" yaml:"name"`
	Ver     string   `json:"ver" yaml:"ver"`
	Targets []Target `json:"targets" yaml:"targets"`

	Packets []Dep `json:"packets,omitempty" yaml:"packets,omitempty"`
}

func Read(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var m Manifest
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", "yml":
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("yaml: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("json: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported manifest extension %q", ext)
	}

	return &m, nil
}
