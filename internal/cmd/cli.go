package cmd

import "errors"

var ErrNoEntrypoint error = errors.New("no entrypoint")

type CliNode struct {
	namespace  string
	parent     *CliNode
	children   []CliNode
	entrypoint func([]string) error
}

func NewCliNode(namespace string, entrypoint func([]string) error) *CliNode {
	return &CliNode{
		namespace:  namespace,
		entrypoint: entrypoint,
	}
}

// ["hs", "log", "my command"]
// ["something", "else"]
func (n *CliNode) Search(path []string) (*CliNode, []string) {
	if len(path) < 1 {
		return nil, path
	}
	if n.namespace != path[0] {
		return nil, path
	}

	if len(n.children) == 0 || len(path) == 1 {
		return n, path[1:]
	}

	for _, child := range n.children {
		if found, remainingArgs := child.Search(path[1:]); found != nil {
			return found, remainingArgs
		}
	}
	return nil, path
}

func (n *CliNode) Fqdn() string {
	if n.parent == nil {
		return n.namespace
	}
	return n.parent.Fqdn() + n.namespace
}

func (n *CliNode) Entrypoint(args []string) error {
	if n.entrypoint == nil {
		return ErrNoEntrypoint
	}
	return n.entrypoint(args)
}

func (n *CliNode) WithChildren(children []CliNode) *CliNode {
	n.children = children
	for _, child := range n.children {
		child.parent = n
	}
	return n
}
