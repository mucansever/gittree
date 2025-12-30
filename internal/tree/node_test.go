package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := NewNode("test")

	assert.Equal(t, "test", node.Name)
	assert.NotNil(t, node.Children)
	assert.Empty(t, node.Children)
}

func TestNode_AddChild(t *testing.T) {
	parent := NewNode("parent")
	child1 := NewNode("child1")
	child2 := NewNode("child2")

	parent.AddChild(child1)
	parent.AddChild(child2)

	assert.Len(t, parent.Children, 2)
	assert.Equal(t, child1, parent.Children[0])
	assert.Equal(t, child2, parent.Children[1])
}
