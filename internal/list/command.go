package list

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/mucansever/gittree/internal/git"
	"github.com/mucansever/gittree/internal/tree"
)

const (
	defaultPath = "."
)

type Options struct {
	Path   string
	Output io.Writer
}

func NewListCommand() *cobra.Command {
	opts := &Options{
		Output: os.Stdout,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List branches in a tree structure",
		Long: `List all branches of a git repository in a hierarchical tree structure.
The tree shows ancestor-descendant relationships between branches.
The current HEAD branch is marked with an asterisk (*).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Path, "path", "p", defaultPath,
		"Path to the git repository")

	return cmd
}

func runList(opts *Options) error {
	repo, err := git.Open(opts.Path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	branches, err := repo.GetBranches()
	if err != nil {
		return fmt.Errorf("failed to get branches: %w", err)
	}

	if len(branches) == 0 {
		fmt.Fprintln(opts.Output, "No branches found")
		return nil
	}

	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	relationships, err := repo.GetBranchRelationships(branches)
	if err != nil {
		return fmt.Errorf("failed to analyze branch relationships: %w", err)
	}

	builder := tree.NewBuilder(relationships)
	t, err := builder.Build(currentBranch)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	printer := tree.NewPrinter(opts.Output)
	printer.Print(t)

	return nil
}
