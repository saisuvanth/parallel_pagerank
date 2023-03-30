package pagerank

import "fmt"

type Graph struct {
	nodes map[int32]float32
	edges []Edge
}

func (g Graph) Init() Graph {
	new_graph := Graph{nodes: make(map[int32]float32), edges: make([]Edge, 0)}
	return new_graph
}

func (g *Graph) CheckNodeExist(id int32) bool {
	_, ok := g.nodes[id]
	return ok
}

func (g *Graph) AddNodes(ids ...int) {
	for _, id := range ids {
		g.nodes[int32(id)] = 1.0
	}
}

func (g *Graph) AddEdge(src int32, dst int32) {
	new_edge := Edge{src: src, dst: dst}
	g.edges = append(g.edges, new_edge)
}

func (g *Graph) AddEdges(edges ...Edge) {
	g.edges = append(g.edges, edges...)
}

func (g *Graph) InitRanks() {
	n := len(g.nodes)
	for nodeId := range g.nodes {
		g.nodes[nodeId] = 1.0 / float32(n)
	}
}

func (g *Graph) GetNodesLen() int {
	return len(g.nodes)
}

func (g *Graph) GetEdgesLen() int {
	return len(g.edges)
}

func (g *Graph) PrintRanks() {
	for nodeId, rank := range g.nodes {
		fmt.Println(nodeId, rank)
	}
}
