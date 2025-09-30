package cli

import (
	"context"
	"errors"
)

var ErrNoEntrypoint error = errors.New("no entrypoint")

type Entrypoint func(context.Context, []string) error

type Node struct {
	namespace  string
	parent     *Node
	Children   []Node
	entrypoint Entrypoint
}

type Option func(*Node)

func WithChildren(children []Node) Option {
	return func(cn *Node) {
		cn.Children = children
		for _, child := range cn.Children {
			child.parent = cn
		}
	}
}

func New(namespace string, entrypoint Entrypoint, opts ...Option) *Node {
	n := &Node{
		namespace:  namespace,
		entrypoint: entrypoint,
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

// ["hs", "log", "my command"]
// ["something", "else"]
func (n *Node) Search(path []string) (*Node, []string) {
	if len(path) < 1 {
		return nil, path
	}
	if n.namespace != path[0] {
		return nil, path
	}

	if len(n.Children) == 0 || len(path) == 1 {
		return n, path[1:]
	}

	for _, child := range n.Children {
		if found, remainingArgs := child.Search(path[1:]); found != nil {
			return found, remainingArgs
		}
	}
	return nil, path
}

func (n *Node) Fqdn() string {
	if n.parent == nil {
		return n.namespace
	}
	return n.parent.Fqdn() + n.namespace
}

func (n *Node) Entrypoint(ctx context.Context, args []string) error {
	if n.entrypoint == nil {
		return ErrNoEntrypoint
	}
	return n.entrypoint(ctx, args)
}
