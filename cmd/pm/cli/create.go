package cli

import (
	"github.com/kurushqosimi/pm/internal"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"github.com/spf13/cobra"
	"log"
	"path"
)

var createCmd = &cobra.Command{
	Use:   "create <manifest>",
	Short: "Pack files and upload to remote repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath := args[0]

		cli, err := sshclient.New(sshclient.Config{
			Host:    host,
			User:    user,
			KeyPath: key,
		})
		if err != nil {
			return err
		}
		defer func() {
			_ = cli.Close()
		}()

		outRelPath, err := internal.Create(manifestPath, repo, cli)
		if err != nil {
			return err
		}

		log.Printf("Uploaded to %s\n", path.Join(host, outRelPath))
		return nil
	},
}
