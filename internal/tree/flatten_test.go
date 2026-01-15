package tree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	now := time.Now()

	nodeFeature1 := NewNode("feature1", time.Time{})
	nodeFeature2 := NewNode("feature2", time.Time{})

	nodeMaster := NewNode("master", now.Add(-2*time.Hour))
	nodeMaster.AddChild(nodeFeature1)
	nodeMaster.AddChild(nodeFeature2)

	nodeRoot := NewNode(".", time.Time{})
	nodeRoot.AddChild(nodeMaster)

	tree := &Tree{Root: nodeRoot}

	items := Flatten(tree)

	assert.Len(t, items, 4)

	assert.Equal(t, ".", items[0].BranchName)
	assert.Equal(t, ".", items[0].Text)

	assert.Equal(t, "master", items[1].BranchName)
	assert.Contains(t, items[1].Text, "master")
	assert.Contains(t, items[1].Text, "└── ")

	assert.Equal(t, "feature1", items[2].BranchName)
	assert.Contains(t, items[2].Text, "feature1")
	assert.Contains(t, items[2].Text, "├── ")
	assert.Contains(t, items[2].Text, "    ")

	assert.Equal(t, "feature2", items[3].BranchName)
	assert.Contains(t, items[3].Text, "feature2")
	assert.Contains(t, items[3].Text, "└── ")
	assert.Contains(t, items[3].Text, "    ")
}
