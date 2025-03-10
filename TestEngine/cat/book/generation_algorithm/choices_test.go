package generation_algorithm

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
)

func choices(population []interface{}, weights []float64, k int) []interface{} {
	n := len(population)

	if weights == nil {
		// 如果没有指定权重，则进行均匀随机抽样
		result := make([]interface{}, k)
		for i := 0; i < k; i++ {
			result[i] = population[rand.Intn(n)]
		}
		return result
	}

	if len(weights) != n {
		panic("权重数量与群体大小不匹配")
	}

	// 计算累积权重
	cumWeights := make([]float64, n)
	cumWeights[0] = weights[0]
	for i := 1; i < n; i++ {
		cumWeights[i] = cumWeights[i-1] + weights[i]
	}

	total := cumWeights[n-1]
	if total <= 0.0 {
		panic("权重总和必须大于零")
	}
	if math.IsInf(total, 0) || math.IsNaN(total) {
		panic("权重总和必须是有限值")
	}

	result := make([]interface{}, k)
	for i := 0; i < k; i++ {
		r := rand.Float64() * total
		idx := sort.SearchFloat64s(cumWeights, r)
		result[i] = population[idx]
	}

	return result
}

func TestChoices(t *testing.T) {
	population := []interface{}{"apple", "banana", "orange"}
	weights := []float64{1, 2, 3}

	// 选择 3 个元素，权重分别为 1、2、3
	selected := choices(population, weights, 3)
	fmt.Println(selected) // 可能输出: [orange banana orange]

	// 选择 2 个元素，权重相等
	selected = choices(population, nil, 2)
	fmt.Println(selected) // 可能输出: [apple banana]
}
