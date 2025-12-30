package tree

import (
	"errors"
)

var (
	ErrCyclicDependency = errors.New("cyclic dependency detected")
)

const rootNodeName = "."

type Builder struct {
	relationships map[string]map[string]bool
}

func NewBuilder(relationships map[string]map[string]bool) *Builder {
	return &Builder{
		relationships: relationships,
	}
}

func (b *Builder) Build(currentBranch string) (*Tree, error) {
	working := b.copyRelationships()

	b.pruneRelationships(working)

	if currentBranch != "" {
		markedBranch := currentBranch + "*"
		if children, exists := working[currentBranch]; exists {
			working[markedBranch] = children
			delete(working, currentBranch)

			for branch := range working {
				if working[branch][currentBranch] {
					delete(working[branch], currentBranch)
					working[branch][markedBranch] = true
				}
			}
		}
	}

	topLevelBranches := make(map[string]bool)
	referencedBranches := make(map[string]bool)
	for _, children := range working {
		for child := range children {
			referencedBranches[child] = true
		}
	}
	for branch := range working {
		if !referencedBranches[branch] {
			topLevelBranches[branch] = true
		}
	}

	working[rootNodeName] = topLevelBranches

	nodes := make(map[string]*Node)
	for branch := range working {
		nodes[branch] = NewNode(branch)
	}

	if err := b.buildHierarchy(working, nodes); err != nil {
		return nil, err
	}

	return &Tree{Root: nodes[rootNodeName]}, nil
}

// removes redundant direct links to descendants
func (b *Builder) pruneRelationships(rels map[string]map[string]bool) {
	for _, children := range rels {
		for child1 := range children {
			for child2 := range children {
				if child1 == child2 {
					continue
				}
				// if child2 is reachable from child1, then the link parent -> child2 is redundant
				if b.isDescendant(rels, child1, child2) {
					delete(children, child2)
				}
			}
		}
	}
}

// checks if target is reachable from start using BFS
func (b *Builder) isDescendant(rels map[string]map[string]bool, start, target string) bool {
	visited := make(map[string]bool)
	queue := []string{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for next := range rels[curr] {
			if next == target {
				return true
			}
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

func (b *Builder) buildHierarchy(relationships map[string]map[string]bool, nodes map[string]*Node) error {
	maxIterations := len(relationships) * 2
	iteration := 0

	for len(relationships) > 0 {
		iteration++
		if iteration > maxIterations {
			return ErrCyclicDependency
		}

		emptyBranches := make(map[string]bool)
		for branch, children := range relationships {
			if len(children) == 0 {
				emptyBranches[branch] = true
				delete(relationships, branch)
			}
		}

		if len(emptyBranches) == 0 && len(relationships) > 0 {
			return ErrCyclicDependency
		}

		for branch, children := range relationships {
			for child := range children {
				if emptyBranches[child] {
					delete(children, child)
					nodes[branch].AddChild(nodes[child])
				}
			}
		}
	}

	return nil
}

func (b *Builder) findLeafNodes(relationships map[string]map[string]bool) map[string]bool {
	leaves := make(map[string]bool)
	for branch, descendants := range relationships {
		if len(descendants) == 0 {
			leaves[branch] = true
		}
	}
	return leaves
}

func (b *Builder) copyRelationships() map[string]map[string]bool {
	copy := make(map[string]map[string]bool)
	for branch, descendants := range b.relationships {
		copy[branch] = make(map[string]bool)
		for descendant := range descendants {
			copy[branch][descendant] = true
		}
	}
	return copy
}
