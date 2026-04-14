package twosum

func TwoSum(nums []int, target int) []int {
	indexByValue := make(map[int]int, len(nums))

	for i, v := range nums {
		if j, ok := indexByValue[target-v]; ok {
			return []int{j, i}
		}
		indexByValue[v] = i
	}

	return nil
}
