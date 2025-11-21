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

func TestParallelSum_EdgeCases(t *testing.T) {
	t.Parallel()

	type args struct {
		nums    []int
		workers int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty_slice_workers_1",
			args: args{
				nums:    []int{},
				workers: 1,
			},
			want: 0,
		},
		{
			name: "empty_slice_workers_greater_than_len",
			args: args{
				nums:    []int{},
				workers: 8,
			},
			want: 0,
		},
		{
			name: "nil_slice_workers_positive",
			args: args{
				nums:    nil,
				workers: 4,
			},
			want: 0,
		},
		{
			name: "nil_slice_workers_zero",
			args: args{
				nums:    nil,
				workers: 0,
			},
			want: 0,
		},
		{
			name: "non_empty_workers_zero_treated_as_one",
			args: args{
				nums:    []int{1, 2, 3, 4, 5},
				workers: 0,
			},
			want: Sum([]int{1, 2, 3, 4, 5}),
		},
		{
			name: "non_empty_workers_negative_treated_as_one",
			args: args{
				nums:    []int{10, 20, 30},
				workers: -3,
			},
			want: Sum([]int{10, 20, 30}),
		},
		{
			name: "workers_greater_than_length_still_correct",
			args: args{
				nums:    []int{1, 2, 3},
				workers: 10,
			},
			want: Sum([]int{1, 2, 3}),
		},
		{
			name: "workers_equal_length",
			args: args{
				nums:    []int{5, 5, 5, 5},
				workers: 4,
			},
			want: Sum([]int{5, 5, 5, 5}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ParallelSum(tt.args.nums, tt.args.workers)
			if got != tt.want {
				t.Fatalf("ParallelSum(nums=%v, workers=%d) = %d, want %d",
					tt.args.nums, tt.args.workers, got, tt.want)
			}
		})
	}
}

func FuzzParallelSum_MatchesSequential(f *testing.F) {
	// Seed fuzz cases using []byte → converted to []int inside the fuzzer.
	f.Add([]byte{}, 0)
	f.Add([]byte{1}, 1)
	f.Add([]byte{1, 2, 3, 4, 5}, 2)
	f.Add([]byte{10, 200, 30, 4}, 4)

	f.Fuzz(func(t *testing.T, bnums []byte, workers int) {
		// Convert []byte → []int to approximate arbitrary slices.
		nums := make([]int, len(bnums))
		for i, b := range bnums {
			// Keep values small and stable to avoid overflow.
			nums[i] = int(b)
		}

		// Clamp to prevent pathological runs.
		if len(nums) > 10_000 {
			nums = nums[:10_000]
		}

		want := Sum(nums)
		got := ParallelSum(nums, workers)

		if got != want {
			t.Fatalf("ParallelSum(nums=%v, workers=%d) = %d, want %d",
				nums, workers, got, want)
		}
	})
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
