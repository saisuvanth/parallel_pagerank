package pagerank

import (
	"sync"
)

type PageRank struct {
	value         float64
	incomingNodes []int32
}

func ParallelPageRank(graph *Graph, numIterations int, numWorkers int, dampingFactor float64, threshold float64) map[int32]float64 {
	numNodes := len(graph.nodes)
	initialPageRank := 1.0 / float64(numNodes)
	pageRanks := make(map[int32]*PageRank)
	for nodeId := range graph.nodes {
		pageRanks[nodeId] = &PageRank{
			value:         initialPageRank,
			incomingNodes: make([]int32, 0),
		}
	}

	for _, edge := range graph.edges {
		dstNode := pageRanks[edge.dst]
		dstNode.incomingNodes = append(dstNode.incomingNodes, edge.src)
	}

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
						incomingPageRank += pageRanks[incomingNode].value / float64(graph.nodes[incomingNode])
					}
					pageRank.value = threshold/float64(numNodes) + dampingFactor*incomingPageRank
				}
			}(parts[j])
		}
		wg.Wait()
	}

	totalPageRank := 0.0
	for _, pageRank := range pageRanks {
		totalPageRank += pageRank.value
	}
	for _, pageRank := range pageRanks {
		pageRank.value /= totalPageRank
	}

	finalPageRanks := make(map[int32]float64)
	for nodeId, pageRank := range pageRanks {
		finalPageRanks[nodeId] = pageRank.value
	}
	return finalPageRanks
}

func GetNodeByIndex(graph *Graph, index int) int32 {
	for node, idx := range graph.nodes {
		if int(idx) == index {
			return node
		}
	}
	return -1
}
