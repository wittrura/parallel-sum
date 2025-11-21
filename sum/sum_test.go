package sum_test

import (
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
