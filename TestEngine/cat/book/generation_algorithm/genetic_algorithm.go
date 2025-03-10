package generation_algorithm

import (
	"fmt"
	"math"
	"sort"
)

type SelectionType uint8

const (
	// SelectionTypeUnknown 为了避免零值的问题
	SelectionTypeUnknown SelectionType = iota
	ROULETTE
	TOURNAMENT
)
const (
	DefaultMaxGenerations  int     = 100
	DefaultMutationChance  float64 = 0.1
	DefaultCrossoverChance float64 = 0.7
)

type GeneticAlgorithm struct {
	Population      []SimpleEquation
	Threshold       float64
	MaxGenerations  int
	MutationChance  float64
	CrossoverChance float64
	SelectionType   SelectionType
}

func NewGeneticAlgorithm(initialPopulation []SimpleEquation, threshold float64) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		Population:      initialPopulation,
		Threshold:       threshold,
		MaxGenerations:  DefaultMaxGenerations,
		MutationChance:  DefaultMutationChance,
		CrossoverChance: DefaultCrossoverChance,
		SelectionType:   ROULETTE,
	}
}

// PickRoulette 轮盘式
func (g *GeneticAlgorithm) PickRoulette() []SimpleEquation {
	weights := make([]float64, 0)
	totalFitness := 0.0
	for _, v := range g.Population {
		totalFitness += v.Fitness()
		weights = append(weights, v.Fitness())
	}
	if totalFitness <= 0 {
		weights = make([]float64, 0)
		for _, v := range g.Population {
			weights = append(weights, v.Fitness()/totalFitness)
		}
	}
	//return g.choicesRoulette(weights)
	return g.choices(g.Population, weights, 2)
}

// PickTournament 锦标赛式
func (g *GeneticAlgorithm) PickTournament(numParticipants int) []SimpleEquation {
	return g.choicesTournament(numParticipants)
}

// ReproduceAndReplace 替换群，使用新一代
func (g *GeneticAlgorithm) ReproduceAndReplace() {
	newPopulation := make([]SimpleEquation, 0)
	for len(newPopulation) < len(g.Population) {
		// 选择 2 parents
		//var parent1, parent2 SimpleEquation
		parents := make([]SimpleEquation, 0)
		if g.SelectionType == ROULETTE {
			// 引入 best，保证每一代都有 best
			best := g.findBest()
			parents = append(parents, best)
			parents = append(parents, g.PickRoulette()...)
			//fmt.Printf("选择的父母：\nparent1: \nX: %v Y: %v fitness:%v\n parent2:\nX: %v Y: %v fitness:%v\n",
			//	parent1.x,
			//	parent1.y,
			//	parent1.Fitness(),
			//	parent2.x,
			//	parent2.y,
			//	parent2.Fitness())

		} else if g.SelectionType == TOURNAMENT {
			// 引入 best，保证每一代都有 best
			best := g.findBest()
			parents = append(parents, best)
			parents = append(parents, g.PickTournament(len(g.Population)/2)...)
		}

		// 交叉 2 parents
		val := GetRandFloat64()
		//fmt.Printf("交叉-随机数：%v\n", val)
		if val < g.CrossoverChance {
			parents[0], parents[1] = parents[0].Crossover(parents[1])
			//fmt.Printf("交叉后：\nparent1: \nX: %v Y: %v fitness:%v\n parent2:\nX: %v Y: %v fitness:%v\n",
			//	parent1.x,
			//	parent1.y,
			//	parent1.Fitness(),
			//	parent2.x,
			//	parent2.y,
			//	parent2.Fitness())
			newPopulation = append(newPopulation, parents[0], parents[1])
		} else {
			newPopulation = append(newPopulation, parents[0], parents[1])
		}
	}
	if len(newPopulation) > len(g.Population) {
		newPopulation = newPopulation[:len(newPopulation)-1]
	}
	// 替换为新的群
	for i := range g.Population {
		g.Population[i] = newPopulation[i]
	}
}

// Mutate 变异
func (g *GeneticAlgorithm) Mutate() {
	for i := range g.Population {
		if GetRandFloat64() < g.MutationChance {
			g.Population[i] = g.Population[i].Mutate()
		}
	}
}

func (g *GeneticAlgorithm) Run() SimpleEquation {
	best := g.findBest()

	for i := 0; i < g.MaxGenerations; i++ {
		if best.Fitness() >= g.Threshold {
			return best
		}
		fmt.Printf("Generation %v Best %v Avg %v \n", i, best.Fitness(), g.average())
		g.ReproduceAndReplace()
		g.Mutate()
		//fmt.Printf("run 变异后：%v\n", g.Sprint())
		highest := g.findBest()
		if highest.Fitness() > best.Fitness() {
			best = highest
		}
	}
	return best
}

func (g *GeneticAlgorithm) choicesRoulette(weights []float64) (SimpleEquation, SimpleEquation) {
	// 默认选择2个元素
	// 累计权重
	selected := make([]SimpleEquation, 2)
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}
	if totalWeight <= 0.0 {
		newtTotalWeight := 0.0
		newWeights := make([]float64, len(weights))
		newCumulativeWeight := make([]float64, len(weights))
		for i, w := range weights {
			newWeights[i] = w / totalWeight
			newtTotalWeight += newWeights[i]
			//fmt.Printf("newtTotalWeight: %v\n", newtTotalWeight)
			newCumulativeWeight[i] = newtTotalWeight
		}
		//for _, v := range newCumulativeWeight {
		//	fmt.Printf("newCumulativeWeight: %v\n", v)
		//}

		for i := 0; i < 2; i++ {
			val := GetRandFloat64()
			randomVal := val * newtTotalWeight
			index := binarySearch(newCumulativeWeight, randomVal)

			//fmt.Printf("newCumulativeWeight index: %v\n", index)
			selected[i] = g.Population[index]
			//fmt.Printf("total<0: %v\n", selected[i])
		}
	} else {
		cumulativeWeight := make([]float64, len(weights))
		for i, w := range weights {
			totalWeight += w
			cumulativeWeight[i] = totalWeight
		}
		//for _, v := range cumulativeWeight {
		//	fmt.Printf("cumulativeWeight: %v\n", v)
		//}

		for i := 0; i < 2; i++ {
			val := GetRandFloat64()
			randomVal := val * totalWeight
			index := binarySearch(cumulativeWeight, randomVal)

			//fmt.Printf("cumulativeWeight index: %v\n", index)
			selected[i] = g.Population[index]
			//fmt.Printf("total>0: %v\n", selected[i])
		}
	}

	return selected[0], selected[1]
}

func (g *GeneticAlgorithm) choicesTournament(numParticipants int) []SimpleEquation {
	// 随机选择参赛者
	participants := make([]SimpleEquation, numParticipants)
	for i := 0; i < numParticipants; i++ {
		randVal := GetRandInt(len(g.Population))
		//fmt.Printf("随机数: %v\n", randVal)
		participants[i] = g.Population[randVal]
	}

	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Fitness() > participants[j].Fitness()
	})
	//for _, v := range participants {
	//	fmt.Printf("participants: %v\n", v)
	//}

	//fmt.Printf("结果：participants[0]: %v, participants[1]: %v \n", participants[0], participants[1])
	return participants
}

func (g *GeneticAlgorithm) choices(population []SimpleEquation, weights []float64, k int) []SimpleEquation {
	n := len(population)

	if weights == nil {
		// 如果没有指定权重，则进行均匀随机抽样
		result := make([]SimpleEquation, k)
		for i := 0; i < k; i++ {
			result[i] = population[GetRandInt(n)]
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

	result := make([]SimpleEquation, k)
	for i := 0; i < k; i++ {
		r := GetRandFloat64() * total
		idx := sort.SearchFloat64s(cumWeights, r)
		result[i] = population[idx]
	}

	return result
}

func binarySearch(cumulativeWeight []float64, val float64) int {
	low, high := 0, len(cumulativeWeight)-1
	for low < high {
		mid := (low + high) / 2
		if cumulativeWeight[mid] < val {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}

func (g *GeneticAlgorithm) findBest() SimpleEquation {
	mu.Lock()
	defer mu.Unlock()
	var best = g.Population[0]
	for _, individual := range g.Population {
		//fmt.Printf("X: %v Y: %v individual.Fitness: %v\n", individual.x, individual.y, individual.Fitness())
		if individual.Fitness() > best.Fitness() {
			best = individual
		}
	}
	return best

}

func (g *GeneticAlgorithm) average() float64 {
	if len(g.Population) == 0 {
		return 0
	}
	var totalFitness float64 = 0

	for _, individual := range g.Population {
		totalFitness += individual.Fitness()
	}
	return totalFitness / float64(len(g.Population))
}

func (g *GeneticAlgorithm) Sprint() string {
	var s string
	for _, i := range g.Population {
		s += fmt.Sprintf("X: %v Y: %v Fitness: %v\n", i.x, i.y, i.Fitness())
	}
	return s
}
