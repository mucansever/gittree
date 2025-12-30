package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr error
	}{
		{
			name: "valid repository",
			setup: func(t *testing.T) string {
				return createTestRepo(t)
			},
			wantErr: nil,
		},
		{
			name: "non-existent path",
			setup: func(t *testing.T) string {
				return "/non/existent/path"
			},
			wantErr: ErrNotRepository,
		},
		{
			name: "not a git repository",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				return dir
			},
			wantErr: ErrNotRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)

			repo, err := Open(path)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
			}
		})
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *git.Repository)
		want    string
		wantErr error
	}{
		{
			name: "main branch",
			setup: func(t *testing.T, repo *git.Repository) {
				// Default branch is already main/master
			},
			want:    "master",
			wantErr: nil,
		},
		{
			name: "custom branch",
			setup: func(t *testing.T, repo *git.Repository) {
				createBranch(t, repo, "feature")
				checkoutBranch(t, repo, "feature")
			},
			want:    "feature",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestRepo(t)
			gitRepo := openGitRepo(t, path)

			if tt.setup != nil {
				tt.setup(t, gitRepo)
			}

			repo := &Repository{repo: gitRepo}
			branch, err := repo.GetCurrentBranch()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, branch)
			}
		})
	}
}

func TestGetBranches(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, repo *git.Repository)
		wantCount int
		wantNames []string
	}{
		{
			name: "single branch",
			setup: func(t *testing.T, repo *git.Repository) {
				// Just master/main
			},
			wantCount: 1,
			wantNames: []string{"master"},
		},
		{
			name: "multiple branches",
			setup: func(t *testing.T, repo *git.Repository) {
				createBranch(t, repo, "develop")
				createBranch(t, repo, "feature")
			},
			wantCount: 3,
			wantNames: []string{"master", "develop", "feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestRepo(t)
			gitRepo := openGitRepo(t, path)

			if tt.setup != nil {
				tt.setup(t, gitRepo)
			}

			repo := &Repository{repo: gitRepo}
			branches, err := repo.GetBranches()

			require.NoError(t, err)
			assert.Len(t, branches, tt.wantCount)

			names := make([]string, len(branches))
			for i, b := range branches {
				names[i] = b.Name
			}
			assert.ElementsMatch(t, tt.wantNames, names)

			// Verify all branches have valid commits
			for _, b := range branches {
				assert.NotEmpty(t, b.Name)
				assert.NotEqual(t, plumbing.ZeroHash, b.Hash)
				assert.NotNil(t, b.Commit)
			}
		})
	}
}

func TestGetBranchRelationships(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, repo *git.Repository) []Branch
		check func(t *testing.T, relationships map[string]map[string]bool)
	}{
		{
			name: "linear history",
			setup: func(t *testing.T, repo *git.Repository) []Branch {
				// master -> feature (feature is ahead)
				commitFile(t, repo, "file1.txt", "content1")
				createBranch(t, repo, "feature")
				checkoutBranch(t, repo, "feature")
				commitFile(t, repo, "file2.txt", "content2")

				r := &Repository{repo: repo}
				branches, err := r.GetBranches()
				require.NoError(t, err)
				return branches
			},
			check: func(t *testing.T, rel map[string]map[string]bool) {
				// master should be ancestor of feature
				assert.True(t, rel["master"]["feature"])
				// feature should not be ancestor of master
				assert.False(t, rel["feature"]["master"])
			},
		},
		{
			name: "diverged branches",
			setup: func(t *testing.T, repo *git.Repository) []Branch {
				// Create common base
				commitFile(t, repo, "base.txt", "base")

				// Create branch1
				createBranch(t, repo, "branch1")
				checkoutBranch(t, repo, "branch1")
				commitFile(t, repo, "file1.txt", "content1")

				// Create branch2 from master
				checkoutBranch(t, repo, "master")
				createBranch(t, repo, "branch2")
				checkoutBranch(t, repo, "branch2")
				commitFile(t, repo, "file2.txt", "content2")

				r := &Repository{repo: repo}
				branches, err := r.GetBranches()
				require.NoError(t, err)
				return branches
			},
			check: func(t *testing.T, rel map[string]map[string]bool) {
				// master should be ancestor of both
				assert.True(t, rel["master"]["branch1"])
				assert.True(t, rel["master"]["branch2"])
				// branch1 and branch2 should not be ancestors of each other
				assert.False(t, rel["branch1"]["branch2"])
				assert.False(t, rel["branch2"]["branch1"])
			},
		},
		{
			name: "same commit",
			setup: func(t *testing.T, repo *git.Repository) []Branch {
				commitFile(t, repo, "file.txt", "content")
				createBranch(t, repo, "same")

				r := &Repository{repo: repo}
				branches, err := r.GetBranches()
				require.NoError(t, err)
				return branches
			},
			check: func(t *testing.T, rel map[string]map[string]bool) {
				// same commit branches should not be ancestors of each other
				assert.False(t, rel["master"]["same"])
				assert.False(t, rel["same"]["master"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTestRepo(t)
			gitRepo := openGitRepo(t, path)

			branches := tt.setup(t, gitRepo)

			repo := &Repository{repo: gitRepo}
			relationships, err := repo.GetBranchRelationships(branches)

			require.NoError(t, err)
			assert.Len(t, relationships, len(branches))
			tt.check(t, relationships)
		})
	}
}

func TestNormalizeBranchName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "refs/heads/master",
			want:  "master",
		},
		{
			input: "refs/heads/feature/my-feature",
			want:  "feature/my-feature",
		},
		{
			input: "master",
			want:  "master",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeBranchName(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func createTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	w, err := repo.Worktree()
	require.NoError(t, err)

	filename := filepath.Join(dir, "README.md")
	err = os.WriteFile(filename, []byte("# Test"), 0644)
	require.NoError(t, err)

	_, err = w.Add("README.md")
	require.NoError(t, err)

	// FIX: Provide an explicit Author and Committer
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	return dir
}

func openGitRepo(t *testing.T, path string) *git.Repository {
	t.Helper()

	repo, err := git.PlainOpen(path)
	require.NoError(t, err)
	return repo
}

func createBranch(t *testing.T, repo *git.Repository, name string) {
	t.Helper()

	head, err := repo.Head()
	require.NoError(t, err)

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), head.Hash())
	err = repo.Storer.SetReference(ref)
	require.NoError(t, err)
}

func checkoutBranch(t *testing.T, repo *git.Repository, name string) {
	t.Helper()

	w, err := repo.Worktree()
	require.NoError(t, err)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
	})
	require.NoError(t, err)
}

func commitFile(t *testing.T, repo *git.Repository, filename, content string) {
	t.Helper()

	w, err := repo.Worktree()
	require.NoError(t, err)

	path := filepath.Join(w.Filesystem.Root(), filename)
	err = os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	_, err = w.Add(filename)
	require.NoError(t, err)

	_, err = w.Commit("Add "+filename, &git.CommitOptions{})
	require.NoError(t, err)
}
