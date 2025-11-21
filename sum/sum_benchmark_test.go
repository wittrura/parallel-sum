package sum_test

import (
	"strconv"
	"testing"

	. "example.com/parallel-sum/sum"
)

// helper to build a deterministic slice of ints
func makeInts(size int) []int {
	nums := make([]int, size)
	for i := range nums {
		nums[i] = i + 1 // 1..size
	}
	return nums
}

func benchmarkSum(b *testing.B, size int) {
	nums := makeInts(size)

	for b.Loop() {
		_ = Sum(nums)
	}
}

func benchmarkParallelSum(b *testing.B, size, workers int) {
	nums := makeInts(size)

	for b.Loop() {
		_ = ParallelSum(nums, workers)
	}
}

func BenchmarkSum_Sizes(b *testing.B) {
	sizes := []int{
		100,
		10_000,
		1_000_000,
	}

	for _, size := range sizes {
		b.Run("size_"+strconv.Itoa(size), func(b *testing.B) {
			benchmarkSum(b, size)
		})
	}
}

func BenchmarkParallelSum_SizesAndWorkers(b *testing.B) {
	type cfg struct {
		size    int
		workers int
	}

	configs := []cfg{
		{size: 100, workers: 1},
		{size: 100, workers: 2},
		{size: 100, workers: 4},

		{size: 10_000, workers: 1},
		{size: 10_000, workers: 2},
		{size: 10_000, workers: 4},
		{size: 10_000, workers: 8},

		{size: 1_000_000, workers: 1},
		{size: 1_000_000, workers: 2},
		{size: 1_000_000, workers: 4},
		{size: 1_000_000, workers: 8},
	}

	for _, c := range configs {
		name := "size_" + strconv.Itoa(c.size) + "_workers_" + strconv.Itoa(c.workers)

		b.Run(name, func(b *testing.B) {
			benchmarkParallelSum(b, c.size, c.workers)
		})
	}
}
