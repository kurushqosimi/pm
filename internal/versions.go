package internal

import (
	"github.com/Masterminds/semver/v3"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"strings"
)

func ListVersions(dir string, remote sshclient.Storage) ([]*semver.Version, error) {
	entries, err := remote.List(dir)
	if err != nil {
		return nil, err
	}

	var out []*semver.Version
	for _, e := range entries {
		if e.IsDir {
			continue
		}
		if !strings.HasSuffix(e.Name, ".tar.gz") {
			continue
		}
		verStr := strings.TrimSuffix(e.Name, ".tar.gz")
		v, err := semver.NewVersion(verStr)
		if err != nil {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func PickBest(versions []*semver.Version, c *semver.Constraints) (best *semver.Version, ok bool) {
	for _, v := range versions {
		if c.Check(v) {
			if best == nil || v.GreaterThan(best) {
				best = v
				ok = true
			}
		}
	}
	return
}

func BestVersion(dir string, c *semver.Constraints, remote sshclient.Storage) (string, bool, error) {
	vers, err := ListVersions(dir, remote)
	if err != nil {
		return "", false, err
	}

	best, ok := PickBest(vers, c)
	if !ok {
		return "", false, nil
	}
	return best.String(), true, nil
}
