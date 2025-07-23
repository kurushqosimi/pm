package fsselect

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/kurushqosimi/pm/pkg/manifest"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Collect(baseDir string, targets []manifest.Target) ([]string, error) {
	selected := make(map[string]struct{})

	fsRoot := os.DirFS(baseDir)

	for _, target := range targets {
		pattern := strings.TrimPrefix(target.Path, "./")
		includes, err := doublestar.Glob(fsRoot, pattern)
		if err != nil {
			return nil, err
		}

		var exPatterns []string
		//for _, ex := range target.Exclude {
		//	exPatterns = append(exPatterns, ex)
		//}

		exPatterns = append(exPatterns, target.Exclude...)

		for _, relPath := range includes {
			skip := false
			for _, ex := range exPatterns {
				filename := filepath.Base(relPath)
				match, _ := doublestar.PathMatch(ex, filename)
				if match {
					skip = true
					break
				}
			}
			if !skip {
				selected[relPath] = struct{}{}
			}
		}
	}

	out := make([]string, 0, len(selected))
	for p := range selected {
		if fi, err := fs.Stat(fsRoot, p); err == nil && !fi.IsDir() {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out, nil
}
