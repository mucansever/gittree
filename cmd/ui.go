package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/mucansever/gittree/internal/git"
	"github.com/mucansever/gittree/internal/tree"
	"github.com/mucansever/gittree/internal/tui"
)

var uiPath string

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Interactive branch tree UI",
	Long:  `Visualizes git branches in an interactive tree structure, allowing you to navigate and checkout branches.`,
	RunE:  runUI,
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().StringVarP(&uiPath, "path", "p", ".", "Path to the git repository")
}

func runUI(cmd *cobra.Command, args []string) error {
	repo, err := git.Open(uiPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	branches, err := repo.GetBranches()
	if err != nil {
		return fmt.Errorf("failed to get branches: %w", err)
	}

	if len(branches) == 0 {
		fmt.Println("No branches found")
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

	meta := make(map[string]time.Time)
	for _, b := range branches {
		meta[b.Name] = b.Commit.Committer.When
	}

	builder := tree.NewBuilder(relationships, meta)
	t, err := builder.Build(currentBranch)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	items := tree.Flatten(t)

	p := tea.NewProgram(tui.NewModel(items, repo))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
