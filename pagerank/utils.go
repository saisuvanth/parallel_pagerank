package pagerank

import (
	"context"
	"sync"
)

func Combiner(ctx context.Context, inputs ...<-chan workerStruct) <-chan workerStruct {
	out := make(chan workerStruct)

	var wg sync.WaitGroup
	multiplexer := func(p <-chan workerStruct) {
		defer wg.Done()

		for in := range p {
			select {
			case <-ctx.Done():
			case out <- in:
			}
		}
	}

	// add length of input channels to be consumed by mutiplexer
	wg.Add(len(inputs))
	for _, in := range inputs {
		go multiplexer(in)
	}

	// close channel after all inputs channels are closed
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
