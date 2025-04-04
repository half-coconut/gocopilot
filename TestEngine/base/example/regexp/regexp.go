package regexp

import (
	"fmt"
	"github.com/dlclark/regexp2"
)

func forRange() {
	arr := [3]int{1, 2, 3}

	for _, v := range arr {
		println(v)
	}
	for idx := range arr {
		println(idx)
	}

	m := map[string]string{
		"name": "egg yolk",
		"age":  "5 month",
	}
	for _, value := range m {
		println(value)
	}
	for key := range m {
		println(key)
	}
}

func demo() {
	arr := []int{2, 4, 6, 8, 10}
	// 左闭右开
	fmt.Printf("%+v\n", arr[1:3])
	fmt.Printf("%+v\n", arr[1:len(arr)-1])
}

// const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// 只允许邮箱域名输入小写
const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
const passwordRegex = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*\W).{8,}$`

func RegexpEmail(str string) {
	//match, _ := regexp.MatchString(emailRegex, str)
	match, _ := regexp2.MustCompile(emailRegex, regexp2.None).MatchString(str)
	fmt.Printf("%v\n", match)
}

func RegexpPassword(str string) {
	//match, _ := regexp.MatchString(passwordRegex, str)
	match, _ := regexp2.MustCompile(passwordRegex, regexp2.None).MatchString(str)
	fmt.Printf("%v\n", match)
}
