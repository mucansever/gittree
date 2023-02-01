package list

import (
	"fmt"
	"strings"
)

type Node struct{
	name string
	children []*Node
}

type Tree struct {
	root *Node
}

func MakeTree(branches map[string]map[string]bool) Tree {
	root := ""
	nodes := make(map[string]*Node)
	for branch, _ := range branches {
		nodes[branch] = &Node{branch, []*Node{}}
	}

	for len(branches) > 0 {
		emptyBranches := make(map[string]bool)
		for branch, children := range branches {
			if len(children) == 0 {
				emptyBranches[branch] = true
				delete(branches, branch)
				root = branch // last to be removed will be master
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

func (tree Tree) Print(head string) {
	printDfs(*tree.root, 0, head)
}

func printDfs(node Node, level int, head string) {
	fmt.Printf(strings.Repeat("  ", level) + node.name[11:])
	if node.name == head { fmt.Println("*") } else { fmt.Println(" ") }
	for _, child := range node.children {
		printDfs(*child, level+1, head)
	}
}