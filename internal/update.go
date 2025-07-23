package internal

import (
	"github.com/kurushqosimi/pm/pkg/manifest"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"log"
	"path"
)

func Update(packagesPath, repoRoot, localDir string, remote sshclient.Storage) error {
	pf, err := manifest.ReadPackages(packagesPath)
	if err != nil {
		return err
	}

	for _, dep := range pf.Packages {
		c, err := manifest.ParseConstraint(dep.Ver)
		if err != nil {
			log.Printf("skip %s: bad constraint %q: %v", dep.Name, dep.Ver, err)
			continue
		}

		remoteDir := path.Join(repoRoot, dep.Name)
		best, ok, err := BestVersion(remoteDir, c, remote)
		if err != nil {
			return err
		}
		if !ok {
			log.Printf("WARN: %s %s – подходящей версии нет", dep.Name, dep.Ver)
			continue
		}

		if ExistsLocally(dep.Name, best, localDir) {
			log.Printf("%s %s уже есть локально", dep.Name, best)
			continue
		}

		if err := DownloadAndExtract(dep.Name, best, repoRoot, localDir, remote); err != nil {
			return err
		}
		log.Printf("%s %s установлена", dep.Name, best)
	}
	return nil
}
