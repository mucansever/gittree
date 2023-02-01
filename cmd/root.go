package cmd

import (
	"os"

	"github.com/mucansever/gittree/cmd/list"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
	Use:   "gittree",
	Short: "List branches of a git repository in a tree structure",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) { 
		list.ListCmd.Run(cmd, args)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addCommands() {
	rootCmd.AddCommand(list.ListCmd)
}

func init() {
	addCommands()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


