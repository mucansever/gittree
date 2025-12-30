package git

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	refPrefix = "refs/heads/"
)

var (
	ErrNotRepository = errors.New("not a git repository")
	ErrDetachedHead  = errors.New("HEAD is detached")
)

type Repository struct {
	repo *git.Repository
}

type Branch struct {
	Name   string
	Hash   plumbing.Hash
	Commit *object.Commit
}

func Open(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, ErrNotRepository
		}
		return nil, fmt.Errorf("opening repository: %w", err)
	}

	return &Repository{repo: repo}, nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	headRef, err := r.repo.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return "", ErrDetachedHead
		}
		return "", fmt.Errorf("getting HEAD: %w", err)
	}

	return normalizeBranchName(headRef.Name().String()), nil
}

func (r *Repository) GetBranches() ([]Branch, error) {
	branchIter, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("getting branches: %w", err)
	}
	defer branchIter.Close()

	var branches []Branch

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		commit, err := r.repo.CommitObject(ref.Hash())
		if err != nil {
			return fmt.Errorf("getting commit for %s: %w", ref.Name(), err)
		}

		branches = append(branches, Branch{
			Name:   normalizeBranchName(ref.Name().String()),
			Hash:   ref.Hash(),
			Commit: commit,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("iterating branches: %w", err)
	}

	return branches, nil
}

// Returns a map where each branch name maps to a set of its descendants
func (r *Repository) GetBranchRelationships(branches []Branch) (map[string]map[string]bool, error) {
	relationships := make(map[string]map[string]bool)

	for i := range branches {
		relationships[branches[i].Name] = make(map[string]bool)
	}

	for i := 0; i < len(branches); i++ {
		for j := i + 1; j < len(branches); j++ {
			// Skip if same commit
			if branches[i].Hash == branches[j].Hash {
				continue
			}

			// Check if i is ancestor of j
			isAncestor, err := branches[i].Commit.IsAncestor(branches[j].Commit)
			if err != nil {
				return nil, fmt.Errorf("checking ancestry %s->%s: %w",
					branches[i].Name, branches[j].Name, err)
			}
			if isAncestor {
				relationships[branches[i].Name][branches[j].Name] = true
			}

			// Check if j is ancestor of i
			isAncestor, err = branches[j].Commit.IsAncestor(branches[i].Commit)
			if err != nil {
				return nil, fmt.Errorf("checking ancestry %s->%s: %w",
					branches[j].Name, branches[i].Name, err)
			}
			if isAncestor {
				relationships[branches[j].Name][branches[i].Name] = true
			}
		}
	}

	return relationships, nil
}

func normalizeBranchName(refName string) string {
	return strings.TrimPrefix(refName, refPrefix)
}
