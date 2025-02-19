package book

import (
	"testing"
)

// 算法图解
// 递归
// 快速排序
// BFS

func degui(nums []int) (res int) {
	if len(nums) == 0 {
		return 0
	} else {
		stake := make([]int, 0)
		for _, num := range nums {
			stake = append(stake, num)
		}
		if len(stake) != 0 {
			top := stake[len(stake)-1]
			pop := stake[:len(stake)-1]
			res = top + degui(pop)
		}
		return res
	}
}

func deguiV2(nums []int) (res int) {
	if len(nums) == 0 {
		return 0
	} else {
		s := NewStack[int]()
		s.Make(nums)

		if !s.IsEmpty() {
			top := s.Pop()
			res = top + degui(s.items)
		}
		return res
	}
}

func TestDegui(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	a := degui(nums)
	println(a)
}

func TestDeguiV2(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	a := deguiV2(nums)
	println(a)
}

func quickSort(nums []int) []int {
	if len(nums) < 2 || len(nums) == 0 {
		return nums
	} else {
		// 1. 分而治之，选取某个值作为比较值，进行分组，分为大数组和小数组，
		// 2. 递归，返回排序结果
		match := nums[0]
		greater := make([]int, 0)
		less := make([]int, 0)
		for _, v := range nums[1:] {
			if v > match {
				greater = append(greater, v)
			} else {
				less = append(less, v)
			}
		}
		return append(append(quickSort(less), match), quickSort(greater)...)
	}
}

func TestQuickSort(t *testing.T) {
	nums := []int{3, 2, 1, 5, 4}
	a := quickSort(nums)
	for _, v := range a {
		println(v)
	}
}

func Graph() map[string][]string {
	// 图数据
	graph := make(map[string][]string, 0)
	graph["you"] = []string{"alice", "bob", "claire"}
	graph["bob"] = []string{"anuj", "peggy"}
	graph["alice"] = []string{"peggy"}
	graph["claire"] = []string{"thom", "jonny"}
	graph["anuj"] = []string{}
	graph["peggy"] = []string{}
	graph["thom"] = []string{}
	graph["jonny"] = []string{}
	return graph
}

// BFS Breadth-First Search, BFS
func BFS(name string) bool {
	m := Graph()
	q := NewQueue[string]()
	q.Make(m[name])
	searched := make([]string, 0)

	for !q.IsEmpty() {
		person := q.Get()
		// 只找没有找过的
		if !Contains(searched, person) {
			if person_is_seller(person) {
				println("找到了！" + person + " is a mongo seller!")
				return true
			} else {
				println("没找到，" + person + " is not a mongo seller")
				// 继续搜索，包含第二层的关系
				q.Make(m[person])
				searched = append(searched, person)
			}
		}
	}
	return false
}

func person_is_seller(name string) bool {
	c := []rune(name)
	return string(c[len(c)-1]) == "m"
}

func TestBreadthFirstSearch(t *testing.T) {
	BFS("you")
}

func Bgrid() [][]byte {
	sg := [][]string{
		{"1", "1", "0", "0", "0"},
		{"1", "1", "0", "0", "0"},
		{"0", "0", "1", "0", "0"},
		{"0", "0", "0", "1", "1"},
	}

	bg := make([][]byte, len(sg))
	for i := range sg {
		bg[i] = make([]byte, len(sg[i]))
		for j := range sg[i] {
			bg[i][j] = sg[i][j][0]
		}
	}
	return bg
}

func TestNumsLands(t *testing.T) {
	g := Bgrid()
	a := numIslands(g)
	println(a)
}

// numIslands 需要判断上下左右的元素都是0 时，判断为一个岛屿
func numIslands(grid [][]byte) int {
	nr := len(grid)
	nc := len(grid[0])
	if nr == 0 {
		return 0
	}
	numbers_islands := 0

	for r := 0; r < nr; r++ {
		for c := 0; c < nc; c++ {
			val := string(grid[r][c])
			if val == "1" {
				numbers_islands += 1
				grid[r][c] = 0
				neighbors := NewQueue[[]int]()
				neighbors.Put([]int{r, c})
				for !neighbors.IsEmpty() {
					v := neighbors.Get()
					row, col := v[0], v[1]
					directions := [][]int{
						{row - 1, col},
						{row + 1, col},
						{row, col - 1},
						{row, col + 1}}
					for _, dir := range directions {
						x, y := dir[0], dir[1]
						if isValid(x, y, grid) {
							neighbors.Put([]int{x, y})
							grid[x][y] = 0
						}
					}
				}
			}
		}
	}
	return numbers_islands
}

func isValid(row, col int, grid [][]byte) bool {
	return row >= 0 && row < len(grid) && col >= 0 && col < len(grid[0]) && string(grid[row][col]) == "1"
}
