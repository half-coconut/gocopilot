package book

func Contains[T comparable](nums []T, target T) bool {
	for _, v := range nums {
		if target == v {
			return true
		}
	}
	return false
}
