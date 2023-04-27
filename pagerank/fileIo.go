package pagerank

import (
	"bufio"
	"context"
	"os"
)

func ReadFile(ctx context.Context, file *os.File, batchSize int, rowsBatch *[]string) <-chan []string {
	out := make(chan []string)

	scanner := bufio.NewScanner(file)

	go func() {
		defer close(out)

		for {
			scanned := scanner.Scan()

			select {
			case <-ctx.Done():
				return
			default:
				row := scanner.Text()
				// fmt.Print(row)
				if len(*rowsBatch) == batchSize || !scanned {
					out <- *rowsBatch
					*rowsBatch = []string{}
				}
				*rowsBatch = append(*rowsBatch, row)
			}

			if !scanned {
				return
			}
		}
	}()
	return out
}

func SaveMaptoFile(nodes *map[int32]float64, filename string) {
	// file, err := os.Create(filename)
	// if err != nil {
	// panic(err)
	// }
	// defer file.Close()
	// for nodeId, rank := range *nodes {
	// file.WriteString(fmt.Sprintf("%d : %.12f\n", nodeId, rank))
	// }
}
