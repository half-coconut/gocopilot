package grpc

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"math/rand"
	"sync"
	"testing"
)

type Node struct {
	name          string
	weight        int
	currentWeight int
}

func (n *Node) Invoke() {}

func TestSmoothWRR(t *testing.T) {
	nodes := []*Node{
		{
			name:          "A",
			weight:        10,
			currentWeight: 10,
		},
		{
			name:          "B",
			weight:        20,
			currentWeight: 20,
		},
		{
			name:          "C",
			weight:        30,
			currentWeight: 30,
		},
	}

	b := &Balance{
		nodes: nodes,
		t:     t,
	}
	for i := 1; i < 6; i++ {
		t.Log(fmt.Sprintf("第 %d 个请求挑选前，nodes: %v", i, slice.Map(nodes, func(idx int, src *Node) Node {
			return *src
		})))
		target := b.wrr()
		target.Invoke()
		t.Log(fmt.Sprintf("第 %d 个请求挑选前，nodes: %v", i, slice.Map(nodes, func(idx int, src *Node) Node {
			return *src
		})))
	}

}

type Balance struct {
	nodes []*Node
	lock  sync.Mutex
	t     *testing.T
}

func (b *Balance) wrr() *Node {
	b.lock.Lock()
	defer b.lock.Unlock()

	total := 0
	// 总权重
	for _, v := range b.nodes {
		total += v.weight
	}
	// 更新了当前权重
	for _, v := range b.nodes {
		v.currentWeight = v.currentWeight + v.weight
	}

	var target *Node
	for _, v := range b.nodes {
		if target == nil {
			target = v
		} else {
			if target.currentWeight < v.currentWeight {
				target = v
			}
		}
	}

	b.t.Log("更新了当前权重后", slice.Map(b.nodes, func(idx int, src *Node) Node {
		return *src
	}))
	b.t.Log("选中了", target)
	target.currentWeight = target.currentWeight - total
	b.t.Log("选中的节点的当前权重，减去总权重后", target)
	return target
}

func (b *Balance) randomPick() *Node {
	var total int32 = 60
	r := rand.Int31n(total)
	for _, n := range b.nodes {
		r -= int32(n.weight)
		if r < 0 {
			return n
		}
	}
	panic("abc")
}
