package cli

import (
	"github.com/kurushqosimi/pm/internal"
	"github.com/kurushqosimi/pm/pkg/sshclient"
	"github.com/spf13/cobra"
)

var (
	local string
)

var updateCmd = &cobra.Command{
	Use:   "update <packages.json>",
	Short: "Download and install packages listed in packages.json",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgPath := args[0]
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

		return internal.Update(pkgPath, repo, local, cli)
	},
}

func init() {
	updateCmd.Flags().StringVar(&local, "local", "./packages", "local install die")
}
