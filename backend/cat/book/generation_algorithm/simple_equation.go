package generation_algorithm

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"math/rand"
	"time"
)

var r2 = rand.New(rand.NewSource(time.Now().UnixNano()))

type SimpleEquationV2 struct {
	x int
	y int
}

// Fitness 适应度，越高表示越优解
func (s *SimpleEquationV2) Fitness() float64 {
	return float64(6*s.x - s.x*s.x + 4*s.y - s.y*s.y)
}

// RandomInstance 随机实例
func (s *SimpleEquationV2) RandomInstance() SimpleEquationV2 {
	return SimpleEquationV2{
		x: r2.Intn(100),
		y: r2.Intn(100),
	}
}

// Crossover 交叉
func (s *SimpleEquationV2) Crossover(other SimpleEquationV2) (SimpleEquationV2, SimpleEquationV2) {
	// 这里需要深拷贝
	child1 := deepcopy.Copy(s).(SimpleEquationV2)
	child2 := deepcopy.Copy(other).(SimpleEquationV2)
	child1.y = other.y
	child2.y = s.y
	return child1, child2
}

// Mutate 变异
func (s *SimpleEquationV2) Mutate() {
	if r2.Float64() > 0.5 {
		//  mutate x
		if r2.Float64() > 0.5 {
			s.x += 1
		} else {
			s.x -= 1
		}
	} else {
		// otherwise mutate y
		if r2.Float64() > 0.5 {
			s.y += 1
		} else {
			s.y -= 1
		}
	}
}

func (s *SimpleEquationV2) Sprint() string {
	return fmt.Sprintf("X: %v,Y: %v Fitness: %v", s.x, s.y, s.Fitness())
}
