package day01

import "testing"

type TreeNode struct {
	val   int
	left  *TreeNode
	right *TreeNode
}

// 初始化二叉树
func InitBinaryTree() *TreeNode {
	root := &TreeNode{val: 1}
	root.left = &TreeNode{val: 2}
	root.right = &TreeNode{val: 3}

	root.left.left = &TreeNode{val: 4}
	root.left.right = &TreeNode{val: 5}

	root.right.left = &TreeNode{val: 6}
	root.right.right = &TreeNode{val: 7}
	return root
}

// 广度优先遍历二叉树
func BFS(root *TreeNode) []int {
	queue := []*TreeNode{root}
	var res []int

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		res = append(res, node.val)

		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}
	return res
}

func TestBFS(t *testing.T) {
	root := InitBinaryTree()
	res := BFS(root)
	for _, i := range res {
		t.Log(i)
	}

}
