package graph

import "github.com/seiyab/gost/utils"

type Directed[T comparable] struct {
	Nodes utils.Set[T]
	Edges []DEdge[T]
}

type DEdge[T comparable] struct {
	From T
	To   T
}

func NewDirected[T comparable]() Directed[T] {
	return Directed[T]{Nodes: utils.NewSet[T]()}
}

func (g *Directed[T]) AddNode(node T) {
	g.Nodes.Add(node)
}

func (g *Directed[T]) AddEdge(from, to T) {
	g.Edges = append(g.Edges, DEdge[T]{From: from, To: to})
}

func (g *Directed[T]) LookupBackward(nodes ...T) utils.Set[T] {
	lookup := map[T]utils.Set[T]{}
	for _, edge := range g.Edges {
		if _, ok := lookup[edge.To]; !ok {
			lookup[edge.To] = utils.NewSet[T]()
		}
		lookup[edge.To].Add(edge.From)
	}

	visited := utils.NewSet[T]()
	result := utils.NewSet[T]()
	var rec func(T)
	rec = func(node T) {
		if visited.Has(node) {
			return
		}
		visited.Add(node)
		result.Add(node)
		for parent := range lookup[node] {
			rec(parent)
		}
	}

	for _, n := range nodes {
		rec(n)
	}

	return result
}
