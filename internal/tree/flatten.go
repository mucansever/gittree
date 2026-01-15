package tree

import (
	"fmt"

	"github.com/mucansever/gittree/internal/timefmt"
)

type Item struct {
	BranchName string
	Text       string
}

func Flatten(t *Tree) []Item {
	if t == nil || t.Root == nil {
		return nil
	}
	return flattenNode(t.Root, "", true)
}

func flattenNode(node *Node, prefix string, isLast bool) []Item {
	var items []Item

	displayName := node.Name
	if !node.LastCommit.IsZero() {
		displayName = fmt.Sprintf("%s (%s ago)", node.Name, timefmt.RelativeTime(node.LastCommit))
	}

	connector := ""

	lineText := ""
	if prefix != "" {
		connector = "├── "
		if isLast {
			connector = "└── "
		}
		lineText = prefix + connector + displayName
	} else {
		lineText = displayName
	}

	items = append(items, Item{
		BranchName: node.Name,
		Text:       lineText,
	})

	childPrefix := prefix
	if prefix == "" {
		childPrefix = ""
	} else {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range node.Children {
		isChildLast := i == len(node.Children)-1

		if prefix == "" {
			childDisplayName := child.Name
			if !child.LastCommit.IsZero() {
				childDisplayName = fmt.Sprintf("%s (%s ago)", child.Name, timefmt.RelativeTime(child.LastCommit))
			}

			childConnector := "├── "
			if isChildLast {
				childConnector = "└── "
			}

			items = append(items, Item{
				BranchName: child.Name,
				Text:       childConnector + childDisplayName,
			})

			grandchildPrefix := "│   "
			if isChildLast {
				grandchildPrefix = "    "
			}

			for j, grandchild := range child.Children {
				isGrandchildLast := j == len(child.Children)-1
				items = append(items, flattenNode(grandchild, grandchildPrefix, isGrandchildLast)...)
			}

		} else {
			items = append(items, flattenNode(child, childPrefix, isChildLast)...)
		}
	}

	return items
}
