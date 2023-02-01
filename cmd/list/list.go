/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mucansever/gittree/cmd/tree"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Set of commands for listing git branches",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		list()	
	},
}

func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func list() {
	repo, err := git.PlainOpen(".")
	CheckIfError(err)
	
	branchChildren := make(map[string]map[string]bool)
	var branchNames []string
	branches, err := repo.Branches()
	CheckIfError(err)

	var commits []*object.Commit
	for ref, err := branches.Next(); err == nil; ref, err = branches.Next() {
		commit, err := repo.CommitObject(ref.Hash())
		CheckIfError(err)

		commits = append(commits, commit)
		branchName := ref.Name().String()
		branchNames = append(branchNames, branchName)
		branchChildren[branchName] = make(map[string]bool)
	}

	for i := 0; i<len(commits); i++ {
		for j := i+1; j<len(commits); j++ {
			if commits[i].Hash == commits[j].Hash {
				continue
			}
			isAncestor, err := commits[i].IsAncestor(commits[j])
			CheckIfError(err)
			if isAncestor { branchChildren[branchNames[i]][branchNames[j]] = true }

			isAncestor, err = commits[j].IsAncestor(commits[i])
			CheckIfError(err)
			if isAncestor { branchChildren[branchNames[j]][branchNames[i]] = true }
		}
	}

	tree := tree.MakeTree(branchChildren)
	tree.Print()
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
