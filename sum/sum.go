package sum

func Sum(nums []int) int {
	result := 0

	for _, num := range nums {
		result += num
	}
	return result
}
