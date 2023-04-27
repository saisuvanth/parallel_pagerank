package main

import (
	"fmt"
	"time"

	"github.com/saisuvanth/parallel_pagerank/pagerank"
)

func main() {

	start := time.Now()
	graph := pagerank.GraphInitFromFileParallel("web-Google.txt", 1000, 64)
	fmt.Println("Time taken to read file parallel: ", time.Since(start))
	fmt.Println(graph.GetNodesLen())
	fmt.Println(graph.GetEdgesLen())
	start = time.Now()
	ranks := pagerank.ParallelPageRank(&graph, 25, 64, 0.85, 0.15)
	fmt.Println("Time taken to Parallel Implementation: ", time.Since(start))
	start = time.Now()
	pagerank.Rank(&graph, 0.85, 25)
	fmt.Println("Time taken to Serial Implementation: ", time.Since(start))
	graph.SaveToFile("ranks.txt")
	pagerank.SaveMaptoFile(&ranks, "ranks.txt")
	// graph.PrintRanks()
	// graph.PrintGraph()
}
