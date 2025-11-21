package sum_test

import (
	"strconv"
	"testing"

	. "example.com/parallel-sum/sum"
)

func TestSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			name: "empty slice",
			nums: []int{},
			want: 0,
		},
		{
			name: "nil slice",
			nums: nil,
			want: 0,
		},
		{
			name: "single element",
			nums: []int{42},
			want: 42,
		},
		{
			name: "multiple positive numbers",
			nums: []int{1, 2, 3, 4, 5},
			want: 15,
		},
		{
			name: "includes negative numbers",
			nums: []int{-5, 10, -3, 4},
			want: 6,
		},
		{
			name: "large slice",
			nums: func() []int {
				nums := make([]int, 1000)
				for i := range nums {
					nums[i] = 1
				}
				return nums
			}(),
			want: 1000,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Sum(tt.nums)
			if got != tt.want {
				t.Fatalf("Sum(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}

func TestParallelSum_MatchesSequentialSum(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name    string
		nums    []int
		workers []int
	}

	tests := []testCase{
		{
			name:    "small_slice_various_workers",
			nums:    []int{1, 2, 3, 4, 5},
			workers: []int{1, 2, 5},
		},
		{
			name:    "single_element_only_one_worker",
			nums:    []int{42},
			workers: []int{1},
		},
		{
			name: "moderate_size_slice",
			nums: func() []int {
				nums := make([]int, 100)
				for i := range nums {
					nums[i] = i + 1 // 1..100
				}
				return nums
			}(),
			workers: []int{1, 2, 4, 8},
		},
		{
			name: "large_slice",
			nums: func() []int {
				nums := make([]int, 10_000)
				for i := range nums {
					nums[i] = 1
				}
				return nums
			}(),
			workers: []int{1, 2, 4, 8},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			want := Sum(tt.nums)

			for _, w := range tt.workers {
				w := w

				// Only test "happy path" where workers > 0 and <= len(nums).
				if w <= 0 || w > len(tt.nums) {
					continue
				}

				t.Run(
					workersSubtestName(w),
					func(t *testing.T) {
						t.Parallel()

						got := ParallelSum(tt.nums, w)
						if got != want {
							t.Fatalf("ParallelSum(nums, workers=%d) = %d, want %d", w, got, want)
						}
					},
				)
			}
		})
	}
}

// workersSubtestName is a tiny helper to keep subtest names readable.
func workersSubtestName(workers int) string {
	return "workers_" + strconv.Itoa(workers)
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		workers int
		want    [][]int
	}{
		{
			name:    "odd length, even workers",
			input:   []int{1, 2, 3, 4, 5, 6, 7},
			workers: 2,
			want:    [][]int{{1, 2, 3, 4}, {5, 6, 7}},
		},
		{
			name:    "odd length, odd workers",
			input:   []int{1, 2, 3, 4, 5, 6, 7},
			workers: 3,
			want:    [][]int{{1, 2, 3}, {4, 5, 6}, {7}},
		},
		{
			name:    "even length",
			input:   []int{1, 2, 3, 4, 5, 6},
			workers: 2,
			want:    [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name:    "even length, odd workers",
			input:   []int{1, 2, 3, 4, 5, 6},
			workers: 3,
			want:    [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Chunk(tt.input, tt.workers)
			if !equal(got, tt.want) {
				t.Fatalf("Chunk got %v, want %d", got, tt.want)
			}
		})
	}
}

func equal(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, temp := range a {
		for j := range temp {
			if len(a[i]) != len(b[i]) {
				return false
			}
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
