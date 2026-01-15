package tree

import (
	"fmt"
	"io"

	"github.com/mucansever/gittree/internal/timefmt"
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
			fmt.Fprintf(p.output, "%s\n", p.formatName(child))
			for j, grandchild := range child.Children {
				p.printNode(grandchild, "", j == len(child.Children)-1)
			}
		}
	} else {
		p.printNode(t.Root, "", true)
	}
}

func (p *Printer) formatName(node *Node) string {
	displayName := node.Name
	if !node.LastCommit.IsZero() {
		displayName = fmt.Sprintf("%s (%s ago)", node.Name, timefmt.RelativeTime(node.LastCommit))
	}
	return displayName
}

func (p *Printer) printNode(node *Node, prefix string, isLast bool) {
	displayName := p.formatName(node)

	if prefix != "" {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Fprintf(p.output, "%s%s%s\n", prefix, connector, displayName)
	} else {
		fmt.Fprintf(p.output, "%s\n", displayName)
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
			fmt.Fprintf(p.output, "%s%s\n", connector, p.formatName(child))
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
