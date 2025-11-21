package sum

import (
	"slices"
)

func Sum(nums []int) int {
	result := 0

	for _, num := range nums {
		result += num
	}
	return result
}

func ParallelSum(nums []int, workers int) int {
	chunks := Chunk(nums, workers)
	numJobs := len(chunks)

	jobs := make(chan []int, numJobs)
	results := make(chan int, numJobs)

	for range workers {
		go worker(jobs, results)
	}

	for _, chunk := range chunks {
		jobs <- chunk
	}
	close(jobs)

	result := 0
	for i := 1; i <= numJobs; i++ {
		result += <-results
	}
	return result
}

func worker(jobs <-chan []int, results chan<- int) {
	job := <-jobs
	sum := Sum(job)
	results <- sum
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
