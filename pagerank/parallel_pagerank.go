package pagerank

import (
	"fmt"
	"math"
	"sync"
)

type PageRank struct {
	value         float64
	incomingNodes []int32
}

func ParallelPageRank(graph *Graph, numIterations int, numWorkers int) {
	// Initialize the PageRank values for each node
	numNodes := len(graph.nodes)
	initialPageRank := 1.0 / float64(numNodes)
	pageRanks := make(map[int32]*PageRank)
	for nodeId := range graph.nodes {
		pageRanks[nodeId] = &PageRank{
			value:         initialPageRank,
			incomingNodes: make([]int32, 0),
		}
	}

	// Compute the incoming nodes for each node
	for _, edge := range graph.edges {
		dstNode := pageRanks[edge.dst]
		dstNode.incomingNodes = append(dstNode.incomingNodes, edge.src)
	}

	// Divide the nodes into parts for parallel processing
	nodeIds := make([]int32, 0, numNodes)
	for nodeId := range graph.nodes {
		nodeIds = append(nodeIds, nodeId)
	}
	parts := make([][]int32, numWorkers)
	for i := range parts {
		start := i * numNodes / numWorkers
		end := (i + 1) * numNodes / numWorkers
		parts[i] = nodeIds[start:end]
	}

	// Perform the PageRank computation in parallel
	for i := 0; i < numIterations; i++ {
		var wg sync.WaitGroup
		wg.Add(numWorkers)
		for j := 0; j < numWorkers; j++ {
			go func(part []int32) {
				defer wg.Done()
				for _, nodeId := range part {
					pageRank := pageRanks[nodeId]
					incomingNodes := pageRank.incomingNodes
					incomingPageRank := 0.0
					for _, incomingNode := range incomingNodes {
						// if index < 12 {
						// 	fmt.Println(incomingNode)
						// }
						incomingPageRank += pageRanks[incomingNode].value / float64(graph.nodes[incomingNode])
					}
					pageRank.value = 0.15/float64(numNodes) + 0.85*incomingPageRank
				}
			}(parts[j])
		}
		wg.Wait()
	}

	// Normalize the PageRank values
	// totalPageRank := 0.0
	// for _, pageRank := range pageRanks {
	// 	totalPageRank += pageRank.value
	// }
	// for _, pageRank := range pageRanks {
	// 	pageRank.value /= totalPageRank
	// }

	// Convert the map of PageRank values to a map of float64 values
	finalPageRanks := make(map[int32]float64)
	for nodeId, pageRank := range pageRanks {
		finalPageRanks[nodeId] = pageRank.value
	}
	SaveMaptoFile(&finalPageRanks)
}

func getNodeByIndex(graph Graph, index int) int32 {
	for node, idx := range graph.nodes {
		if int(idx) == index {
			return node
		}
	}
	return -1
}

// getOutgoingEdges returns the edges that originate from the given node.
func getOutgoingEdges(graph Graph, node int32) []Edge {
	edges := make([]Edge, 0)
	for _, edge := range graph.edges {
		if edge.src == node {
			edges = append(edges, edge)
		}
	}
	return edges
}

func Parallel_pagerank(graph Graph, dampFactor float64, threshold float64) {
	n := len(graph.nodes)
	pagerank := make([]float64, n)
	for i := range pagerank {
		pagerank[i] = 1 / float64(n)
	}

	iteration := 0
	diff := float64(1)
	for diff > threshold {
		iteration++

		// Compute the contribution of each neighbor to the Pagerank score in parallel
		contributions := make([]float64, n)
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				defer wg.Done()
				node := getNodeByIndex(graph, i)
				outgoingEdges := getOutgoingEdges(graph, node)
				for _, e := range outgoingEdges {
					contributions[i] += pagerank[e.src] / float64(len(outgoingEdges))
				}
			}(i)
		}
		wg.Wait()

		// Compute the new Pagerank score for each node in parallel
		newPagerank := make([]float64, n)
		var wg2 sync.WaitGroup
		wg2.Add(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				defer wg2.Done()
				// node := getNodeByIndex(graph, i)
				newPagerank[i] = (1-dampFactor)/float64(n) + dampFactor*contributions[i]
			}(i)
		}
		wg2.Wait()

		// Compute the difference between the old and new Pagerank scores
		diff = float64(0)
		for i := range pagerank {
			diff += math.Abs(newPagerank[i] - pagerank[i])
		}

		// Update the Pagerank scores
		pagerank = newPagerank
	}
	fmt.Println(pagerank)
}
