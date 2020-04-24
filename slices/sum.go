package slices

func Sum(nums []int) (sum int) {
	for _, num := range nums {
		sum += num
	}
	return
}

func SumAll(numbersToSum ...[]int) []int {
	sums := make([]int, len(numbersToSum))

	for i, nums := range numbersToSum {
		sums[i] = Sum(nums)
	}

	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	count := len(numbersToSum)
	if count == 0 {
		return nil
	}
	sums := make([]int, count)
	for i, nums := range numbersToSum {
		if len(nums) == 0 {
			sums[i] = 0
		} else {
			sums[i] = Sum(nums[1:])
		}
	}

	return sums
}
