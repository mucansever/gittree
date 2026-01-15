package tree

import (
	"time"
)

type Node struct {
	Name       string
	Children   []*Node
	LastCommit time.Time
}

type Tree struct {
	Root *Node
}

func NewNode(name string, lastCommit time.Time) *Node {
	return &Node{
		Name:       name,
		Children:   make([]*Node, 0),
		LastCommit: lastCommit,
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
