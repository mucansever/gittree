package tree

import (
	"fmt"
	"io"
)

type Printer struct {
	output io.Writer
}

func NewPrinter(output io.Writer) *Printer {
	return &Printer{output: output}
}

func (p *Printer) Print(t *Tree) {
	if t == nil || t.Root == nil {
		return
	}
	if t.Root.Name == "" {
		for _, child := range t.Root.Children {
			fmt.Fprintf(p.output, "%s\n", child.Name)
			for j, grandchild := range child.Children {
				p.printNode(grandchild, "", j == len(child.Children)-1)
			}
		}
	} else {
		p.printNode(t.Root, "", true)
	}
}

func (p *Printer) printNode(node *Node, prefix string, isLast bool) {
	if prefix != "" {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Fprintf(p.output, "%s%s%s\n", prefix, connector, node.Name)
	} else {
		fmt.Fprintf(p.output, "%s\n", node.Name)
	}

	if len(node.Children) == 0 {
		return
	}

	for i, child := range node.Children {
		childIsLast := i == len(node.Children)-1

		if prefix == "" {
			connector := "├── "
			if childIsLast {
				connector = "└── "
			}
			fmt.Fprintf(p.output, "%s%s\n", connector, child.Name)
			for j, grandchild := range child.Children {
				grandchildIsLast := j == len(child.Children)-1
				grandchildPrefix := "    "
				if !childIsLast {
					grandchildPrefix = "│   "
				}
				p.printNode(grandchild, grandchildPrefix, grandchildIsLast)
			}
		} else {
			childPrefix := prefix
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
			p.printNode(child, childPrefix, childIsLast)
		}
	}
}
