package list

import (
	"fmt"
	"strings"
)

type Node struct {
	name     string
	children []*Node
}

type Tree struct {
	root *Node
}

// CheckIfError panics if the error is not nil.
func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// MakeTree creates a tree from a map of branch names to their descendants.
func MakeTree(branches map[string]map[string]bool) Tree {
	root := "."
	allBranches := make(map[string]bool)
	for branch := range branches {
		allBranches[branch] = true
	}
	branches[root] = allBranches

	nodes := make(map[string]*Node)
	for branch := range branches {
		nodes[branch] = &Node{branch, []*Node{}}
	}

	for len(branches) > 0 {
		emptyBranches := make(map[string]bool)
		for branch, children := range branches {
			if len(children) == 0 {
				emptyBranches[branch] = true
				delete(branches, branch)
			}
		}

		for branch, children := range branches {
			var currRemoved []string
			for child := range children {
				if _, ok := emptyBranches[child]; ok {
					currRemoved = append(currRemoved, child)
					delete(children, child)
				}
			}
			if len(children) == 0 {
				for _, child := range currRemoved {
					nodes[branch].children = append(nodes[branch].children, nodes[child])
				}
			}
		}
	}
	return Tree{nodes[root]}
}

// Print prints the tree in a DFS manner.
func (tree *Tree) Print() {
	printDfs(*tree.root, -1)
}

func printDfs(node Node, level int) {
	if level >= 0 {
		fmt.Println(strings.Repeat("  ", level) + node.name)
	}
	for _, child := range node.children {
		printDfs(*child, level+1)
	}
}
