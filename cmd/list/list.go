package list

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

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

	tree := MakeTree(branchChildren)
	tree.Print()
}

func init() {
}
