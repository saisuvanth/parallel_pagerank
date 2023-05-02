package pagerank

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
)

type workerStruct struct {
	edges []Edge
	nodes []int
}

func NormalizeRanks(graph *Graph) {
	sum := float64(0)
	for _, rank := range graph.nodes {
		sum += rank
	}
	for nodeId, rank := range graph.nodes {
		graph.nodes[nodeId] = rank / sum
	}
}

func Rank(graph *Graph, dampFactor float64, iterations int) {
	n := len(graph.nodes)
	for i := 0; i < iterations; i++ {
		newRanks := make(map[int32]float64)
		for _, edge := range graph.edges {
			newRanks[edge.dst] += graph.nodes[edge.src] / float64(n)
		}
		maxDiff := float64(0)
		for nodeId, rank := range graph.nodes {
			newRank := (1-dampFactor)/float64(n) + dampFactor*newRanks[nodeId]
			diff := newRank - rank
			if diff > maxDiff {
				maxDiff = diff
			}
			graph.nodes[nodeId] = newRank
		}

	}
	NormalizeRanks(graph)
}

func GraphInitFromFile(filename string) Graph {
	graph := Graph{}
	graph = graph.Init()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		nodes := strings.Fields(row)
		src, err := strconv.Atoi(nodes[0])
		if err != nil {
			panic(err)
		}
		dst, err := strconv.Atoi(nodes[1])
		if err != nil {
			panic(err)
		}
		graph.AddEdge(int32(src), int32(dst))
		graph.AddNodes(src, dst)
	}
	return graph
}

func GraphInitFromFileParallel(filename string, batchSize int, workerThreads int) Graph {
	graph := Graph{}
	graph = graph.Init()
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	worker := func(ctx context.Context, rowBatch <-chan []string) <-chan workerStruct {
		out := make(chan workerStruct)
		go func() {
			defer close(out)

			p := workerStruct{}
			for rowBatch := range rowBatch {
				for _, row := range rowBatch {
					nodes := strings.Fields(row)
					src, err := strconv.Atoi(nodes[0])
					if err != nil {
						panic(err)
					}
					dst, err := strconv.Atoi(nodes[1])
					if err != nil {
						panic(err)
					}
					p.edges = append(p.edges, Edge{int32(src), int32(dst)})
					p.nodes = append(p.nodes, src, dst)
				}
			}
			out <- p
		}()
		return out
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rowsBatch := []string{}
	rowBatch := ReadFile(ctx, file, batchSize, &rowsBatch)

	workerChannels := make([]<-chan workerStruct, workerThreads)

	for i := 0; i < workerThreads; i++ {
		workerChannels[i] = worker(ctx, rowBatch)
	}

	for res := range Combiner(ctx, workerChannels...) {
		graph.AddNodes(res.nodes...)
		graph.AddEdges(res.edges...)
	}

	return graph
}
