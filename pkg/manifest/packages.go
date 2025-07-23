package manifest

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type PackagesFile struct {
	Packages []Dep `json:"packages" yaml:"packages"`
}

func ReadPackages(path string) (*PackagesFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var pf PackagesFile
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(raw, &pf); err != nil {
			return nil, fmt.Errorf("yaml: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(raw, &pf); err != nil {
			return nil, fmt.Errorf("json: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported extension %q", ext)
	}

	return &pf, nil
}
