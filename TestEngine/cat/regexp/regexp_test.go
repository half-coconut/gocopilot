package regexp

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	heights := []int{4, 5, 3, 2, 1, 7, 8, 9}
	heights = append([]int{0}, heights...)
	heights = append(heights, 0)
	stake := []int{0}
	fmt.Printf("%v\n", heights)
	fmt.Printf("%v\n", stake)
}

func TestRegexp(t *testing.T) {
	email := "123@Qq.com"
	password := "Cc12345!"
	RegexpEmail(email)
	RegexpPassword(password)
}
