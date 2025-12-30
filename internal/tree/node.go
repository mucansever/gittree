package tree

type Node struct {
	Name     string
	Children []*Node
}

type Tree struct {
	Root *Node
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make([]*Node, 0),
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}
