package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mucansever/gittree/internal/list"
)

var rootCmd = &cobra.Command{
	Use:   "gittree",
	Short: "List branches of a git repository in a tree structure",
	Long: `gittree visualizes git branches in a hierarchical tree structure,
showing ancestor-descendant relationships between branches.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(list.NewListCommand())
}
