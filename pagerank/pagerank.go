package pagerank

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type workerStruct struct {
	edges []Edge
	nodes []int
}

func NormalizeRanks(graph *Graph) {
	sum := float32(0)
	for _, rank := range graph.nodes {
		sum += rank
	}
	for nodeId, rank := range graph.nodes {
		graph.nodes[nodeId] = rank / sum
	}
}

func Rank(graph *Graph, dampFactor float32, threshold float32) {
	n := len(graph.nodes)
	graph.InitRanks()
	fmt.Println(n)
	for {
		newRanks := make(map[int32]float32)
		for _, edge := range graph.edges {
			newRanks[edge.dst] += graph.nodes[edge.src] / float32(n)
		}
		fmt.Println("Completed")
		maxDiff := float32(0)
		for nodeId, rank := range graph.nodes {
			newRank := (1-dampFactor)/float32(n) + dampFactor*newRanks[nodeId]
			diff := newRank - rank
			if diff > maxDiff {
				maxDiff = diff
			}
			graph.nodes[nodeId] = newRank
		}

		if maxDiff < threshold {
			break
		}
	}
	fmt.Println("Completed 1")

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

	workerChan := make([]<-chan workerStruct, workerThreads)

	for i := 0; i < workerThreads; i++ {
		workerChan[i] = worker(ctx, rowBatch)
	}

	for res := range Combiner(ctx, workerChan...) {
		graph.AddNodes(res.nodes...)
		graph.AddEdges(res.edges...)
	}

	return graph
}
