package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	host string
	user string
	key  string
	repo string

	rootCmd = &cobra.Command{
		Use:   "pm",
		Short: "Tiny package manager",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost:22", "SSH host:port")
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "SSH username")
	rootCmd.PersistentFlags().StringVar(&key, "key", "", "path to private key")
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "/repo", "remote repository root")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
}
