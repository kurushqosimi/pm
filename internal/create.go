package internal

import (
	"github.com/kurushqosimi/pm/pkg/archiver"
	"github.com/kurushqosimi/pm/pkg/fsselect"
	"github.com/kurushqosimi/pm/pkg/manifest"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"path"
)

func Create(manifestPath, repoRoot string, ssh sshclient.Storage) (remotePath string, _ error) {
	m, err := manifest.ReadManifest(manifestPath)
	if err != nil {
		return "", err
	}

	baseDir := path.Dir(manifestPath)
	files, err := fsselect.Collect(baseDir, m.Targets)
	if err != nil {
		return "", err
	}

	data, err := archiver.TarGz(files, baseDir)
	if err != nil {
		return "", err
	}

	remotePath = path.Join(repoRoot, m.Name, m.Ver+".tar.gz")
	if err := ssh.Upload(remotePath, data); err != nil {
		return "", err
	}
	return remotePath, nil
}
