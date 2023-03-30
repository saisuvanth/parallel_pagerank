package main

import (
	"fmt"
	"time"

	"github.com/saisuvanth/parallel_pagerank/pagerank"
)

func main() {
	start := time.Now()
	graph := pagerank.GraphInitFromFile("web-Google.txt")
	fmt.Println("Time taken to read file: ", time.Since(start))

	start = time.Now()
	graph = pagerank.GraphInitFromFileParallel("web-Google.txt", 1000, 64)
	fmt.Println("Time taken to read file parallel: ", time.Since(start))
	fmt.Println(graph.GetNodesLen())
	fmt.Println(graph.GetEdgesLen())
	// pagerank.Rank(&graph, 0.85, 0.0001)
	// graph.PrintRanks()
	// graph.PrintGraph()
}
