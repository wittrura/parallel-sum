package sum

import (
	"context"
	"slices"
	"sync"
)

func Sum(nums []int) int {
	result := 0

	for _, num := range nums {
		result += num
	}
	return result
}

func ParallelSum(nums []int, workers int) int {
	result, _ := ParallelSumCtx(context.Background(), nums, workers, nil)
	return result
}

func Chunk(nums []int, chunks int) [][]int {
	length := len(nums)
	size := length / chunks
	if length%chunks != 0 {
		size++
	}

	var results [][]int
	for c := range slices.Chunk(nums, size) {
		results = append(results, c)
	}
	return results
}

func ParallelSumCtx(ctx context.Context, nums []int, workers int, beforeProcessing func()) (int, error) {
	if len(nums) == 0 {
		return 0, nil
	}

	if workers < 1 {
		workers = 1
	} else if workers > len(nums) {
		workers = len(nums)
	}

	// Early exit if already canceled.
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	chunks := Chunk(nums, workers)
	numJobs := len(chunks)

	jobs := make(chan []int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go worker(ctx, &wg, jobs, results, beforeProcessing)
	}

	// Close results once all workers are done.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Feed jobs, respecting cancellation.
	go func() {
		defer close(jobs)

		for _, c := range chunks {
			select {
			case <-ctx.Done():
				return
			case jobs <- c:
			}
		}
	}()

	result := 0
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case r, ok := <-results:
			if !ok {
				return result, nil
			}
			result += r
		}
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan []int, results chan<- int, beforeProcessing func()) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	case job, ok := <-jobs:
		if !ok {
			return
		}
		if beforeProcessing != nil {
			beforeProcessing()
		}

		sum := Sum(job)

		select {
		case <-ctx.Done():
			return
		case results <- sum:
		}
	}
}
