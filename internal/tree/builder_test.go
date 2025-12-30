package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name          string
		relationships map[string]map[string]bool
		currentBranch string
		check         func(t *testing.T, tree *Tree)
		wantErr       error
	}{
		{
			name: "simple linear tree",
			relationships: map[string]map[string]bool{
				"master":  {"feature": true},
				"feature": {},
			},
			currentBranch: "master",
			check: func(t *testing.T, tree *Tree) {
				assert.NotNil(t, tree.Root)
				assert.Equal(t, ".", tree.Root.Name)
				assert.Len(t, tree.Root.Children, 1)

				master := tree.Root.Children[0]
				assert.Equal(t, "master*", master.Name)
				assert.Len(t, master.Children, 1)
				assert.Equal(t, "feature", master.Children[0].Name)
			},
			wantErr: nil,
		},
		{
			name: "branching tree",
			relationships: map[string]map[string]bool{
				"master":   {"feature1": true, "feature2": true},
				"feature1": {},
				"feature2": {},
			},
			currentBranch: "feature1",
			check: func(t *testing.T, tree *Tree) {
				master := tree.Root.Children[0]
				assert.Equal(t, "master", master.Name)
				assert.Len(t, master.Children, 2)

				// Check both features exist
				childNames := []string{master.Children[0].Name, master.Children[1].Name}
				assert.Contains(t, childNames, "feature1*")
				assert.Contains(t, childNames, "feature2")
			},
			wantErr: nil,
		},
		{
			name: "complex tree",
			relationships: map[string]map[string]bool{
				"master":   {"develop": true},
				"develop":  {"feature1": true, "feature2": true},
				"feature1": {"bugfix": true},
				"feature2": {},
				"bugfix":   {},
			},
			currentBranch: "develop",
			check: func(t *testing.T, tree *Tree) {
				master := tree.Root.Children[0]
				assert.Equal(t, "master", master.Name)
				assert.Len(t, master.Children, 1)

				develop := master.Children[0]
				assert.Equal(t, "develop*", develop.Name)
				assert.Len(t, develop.Children, 2)
			},
			wantErr: nil,
		},
		{
			name: "no current branch",
			relationships: map[string]map[string]bool{
				"master":  {"feature": true},
				"feature": {},
			},
			currentBranch: "",
			check: func(t *testing.T, tree *Tree) {
				master := tree.Root.Children[0]
				assert.Equal(t, "master", master.Name)
			},
			wantErr: nil,
		},
		{
			name: "single branch",
			relationships: map[string]map[string]bool{
				"master": {},
			},
			currentBranch: "master",
			check: func(t *testing.T, tree *Tree) {
				assert.Len(t, tree.Root.Children, 1)
				assert.Equal(t, "master*", tree.Root.Children[0].Name)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.relationships)
			tree, err := builder.Build(tt.currentBranch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tree)
				tt.check(t, tree)
			}
		})
	}
}

func TestBuilder_Build_Transitive(t *testing.T) {
	tests := []struct {
		name          string
		relationships map[string]map[string]bool
		currentBranch string
		check         func(t *testing.T, tree *Tree)
	}{
		{
			name: "nested transitive branches",
			relationships: map[string]map[string]bool{
				"main":                  {"feat/feature-1": true, "chore/document-change": true},
				"feat/feature-1":        {"chore/document-change": true},
				"chore/document-change": {},
			},
			currentBranch: "main",
			check: func(t *testing.T, tree *Tree) {
				// main should only have one direct child: feat/feature-1
				main := tree.Root.Children[0]
				assert.Equal(t, "main*", main.Name)
				require.Len(t, main.Children, 1)

				feat := main.Children[0]
				assert.Equal(t, "feat/feature-1", feat.Name)
				require.Len(t, feat.Children, 1)

				chore := feat.Children[0]
				assert.Equal(t, "chore/document-change", chore.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.relationships)
			tree, err := builder.Build(tt.currentBranch)
			require.NoError(t, err)
			tt.check(t, tree)
		})
	}
}

func TestBuilder_CyclicDependency(t *testing.T) {
	relationships := map[string]map[string]bool{
		"branch1": {"branch2": true},
		"branch2": {"branch1": true}, // cycle
	}

	builder := NewBuilder(relationships)
	_, err := builder.Build("")

	assert.ErrorIs(t, err, ErrCyclicDependency)
}

func TestBuilder_copyRelationships(t *testing.T) {
	original := map[string]map[string]bool{
		"master":  {"feature": true},
		"feature": {},
	}

	builder := NewBuilder(original)
	copy := builder.copyRelationships()

	copy["master"]["new"] = true
	assert.False(t, original["master"]["new"])

	delete(copy, "feature")
	assert.Contains(t, original, "feature")
}

func TestBuilder_findLeafNodes(t *testing.T) {
	relationships := map[string]map[string]bool{
		"master":  {"feature": true},
		"feature": {},
		"develop": {},
	}

	builder := NewBuilder(nil)
	leaves := builder.findLeafNodes(relationships)

	assert.Len(t, leaves, 2)
	assert.True(t, leaves["feature"])
	assert.True(t, leaves["develop"])
	assert.False(t, leaves["master"])
}
