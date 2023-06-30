package _type

func InIntArray(k int, arr []int) bool {
	for _, i := range arr {
		if i == k {
			return true
		}
	}
	return false
}
