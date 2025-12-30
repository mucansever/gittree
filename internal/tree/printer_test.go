package tree

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinter_Print(t *testing.T) {
	tests := []struct {
		name string
		tree *Tree
		want string
	}{
		{
			name: "simple tree",
			tree: &Tree{
				Root: &Node{
					Name: "master",
					Children: []*Node{
						{Name: "feature", Children: []*Node{}},
					},
				},
			},
			want: "master\n└── feature\n",
		},
		{
			name: "branching tree",
			tree: &Tree{
				Root: &Node{
					Name: "master",
					Children: []*Node{
						{Name: "feature1", Children: []*Node{}},
						{Name: "feature2", Children: []*Node{}},
					},
				},
			},
			want: "master\n├── feature1\n└── feature2\n",
		},
		{
			name: "deep tree",
			tree: &Tree{
				Root: &Node{
					Name: "master",
					Children: []*Node{
						{
							Name: "develop",
							Children: []*Node{
								{Name: "feature", Children: []*Node{}},
							},
						},
					},
				},
			},
			want: "master\n└── develop\n    └── feature\n",
		},
		{
			name: "complex tree",
			tree: &Tree{
				Root: &Node{
					Name: "main",
					Children: []*Node{
						{Name: "fix/important-bug", Children: []*Node{}},
						{
							Name: "feat/feature-1",
							Children: []*Node{
								{Name: "chore/document-change*", Children: []*Node{}},
							},
						},
						{Name: "chore/no-commits-yet", Children: []*Node{}},
					},
				},
			},
			want: "main\n├── fix/important-bug\n├── feat/feature-1\n│   └── chore/document-change*\n└── chore/no-commits-yet\n",
		},
		{
			name: "nil tree",
			tree: nil,
			want: "",
		},
		{
			name: "empty root",
			tree: &Tree{Root: nil},
			want: "",
		},
		{
			name: "virtual root with multiple roots",
			tree: &Tree{
				Root: &Node{
					Name: "",
					Children: []*Node{
						{Name: "main", Children: []*Node{}},
						{Name: "develop", Children: []*Node{}},
					},
				},
			},
			want: "main\ndevelop\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printer := NewPrinter(&buf)

			printer.Print(tt.tree)

			got := buf.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrinter_PrintNode(t *testing.T) {
	node := &Node{
		Name: "parent",
		Children: []*Node{
			{Name: "child1", Children: []*Node{}},
			{Name: "child2", Children: []*Node{}},
		},
	}

	var buf bytes.Buffer
	printer := NewPrinter(&buf)
	printer.printNode(node, "", true)

	output := buf.String()
	assert.Contains(t, output, "parent")
	assert.Contains(t, output, "├── child1")
	assert.Contains(t, output, "└── child2")
}
