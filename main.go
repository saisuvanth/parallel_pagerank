package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/saisuvanth/parallel_pagerank/pagerank"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <num_of_threads>")
		return
	}
	num_threads, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Usage: go run main.go <num_of_threads>")
		return
	}

	start := time.Now()
	graph := pagerank.GraphInitFromFileParallel("web-Google.txt", 1000, 64)
	fmt.Println("Time taken to read file parallel: ", time.Since(start))
	fmt.Println(graph.GetNodesLen())
	fmt.Println(graph.GetEdgesLen())
	start = time.Now()
	pagerank.ParallelPageRank(&graph, 25, num_threads, 0.85, 0.15)
	fmt.Println(time.Since(start))
	start = time.Now()
	pagerank.Rank(&graph, 0.85, 25)
	fmt.Println(time.Since(start))
	graph.SaveToFile("ranks.txt")
	// pagerank.SaveMaptoFile(&ranks, "ranks.txt")
}
