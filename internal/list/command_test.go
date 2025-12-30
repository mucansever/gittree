package list

import (
	"bytes"
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

func TestRunList(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		wantErr  bool
		checkOut func(t *testing.T, output string)
	}{
		{
			name: "valid repository",
			setup: func(t *testing.T) string {
				return createTestRepo(t)
			},
			wantErr: false,
			checkOut: func(t *testing.T, output string) {
				assert.Contains(t, output, "master*")
			},
		},
		{
			name: "multiple branches",
			setup: func(t *testing.T) string {
				path := createTestRepo(t)
				repo := openRepo(t, path)
				createBranch(t, repo, "feature")
				return path
			},
			wantErr: false,
			checkOut: func(t *testing.T, output string) {
				assert.Contains(t, output, "master*")
				assert.Contains(t, output, "feature")
			},
		},
		{
			name: "invalid repository",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
			checkOut: func(t *testing.T, output string) {
				// Error should be returned, not printed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			var buf bytes.Buffer

			opts := &Options{
				Path:   path,
				Output: &buf,
			}

			err := runList(opts)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.checkOut(t, buf.String())
			}
		})
	}
}

func TestNewListCommand(t *testing.T) {
	cmd := NewListCommand()

	assert.Equal(t, "list", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)

	flag := cmd.Flags().Lookup("path")
	require.NotNil(t, flag)
	assert.Equal(t, "p", flag.Shorthand)
	assert.Equal(t, defaultPath, flag.DefValue)
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

func openRepo(t *testing.T, path string) *git.Repository {
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
