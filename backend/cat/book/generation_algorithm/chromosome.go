package generation_algorithm

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"math/rand"
	"sync"
	"time"
)

var mu sync.Mutex

type Chromosome interface {
	Fitness() float64
	RandomInstance() SimpleEquation
	Crossover(other SimpleEquation) (SimpleEquation, SimpleEquation)
	Mutate()
}

type SimpleEquation struct {
	x          int
	y          int
	Name       string
	FitnessVal float64
}

// Fitness 适应度，越高表示越优解
func (s *SimpleEquation) Fitness() float64 {
	s.FitnessVal = float64(6*s.x - s.x*s.x + 4*s.y - s.y*s.y)
	return s.FitnessVal
}

// RandomInstance 随机实例
func (s *SimpleEquation) RandomInstance() SimpleEquation {
	return SimpleEquation{
		x: GetRandInt(100),
		y: GetRandInt(100),
	}
}

// Crossover 交叉
func (s SimpleEquation) Crossover(other SimpleEquation) (SimpleEquation, SimpleEquation) {
	// 这里需要深拷贝
	child1 := deepcopy.Copy(s).(SimpleEquation)
	child2 := deepcopy.Copy(other).(SimpleEquation)
	child1.x = s.x
	child1.y = other.y
	child2.x = other.x
	child2.y = s.y
	return child1, child2
}

// Mutate 变异
func (s *SimpleEquation) Mutate() SimpleEquation {
	if GetRandFloat64() > 0.5 {
		//  mutate x
		if GetRandFloat64() > 0.5 {
			s.x += 1
		} else {
			s.x -= 1
		}
	} else {
		// otherwise mutate y
		if GetRandFloat64() > 0.5 {
			s.y += 1
		} else {
			s.y -= 1
		}
	}
	//println("s mutate 之后", s.Sprint())
	return *s
}

func (s *SimpleEquation) Sprint() string {
	return fmt.Sprintf("X: %v,Y: %v Fitness: %v", s.x, s.y, s.Fitness())
}

func GetRandInt(val int) int {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu.Lock()
	defer mu.Unlock()
	return r.Intn(val)
}

func GetRandFloat64() float64 {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu.Lock()
	defer mu.Unlock()
	val := r.Float64()
	//fmt.Printf("谁在调用-随机数：%v\n", val)
	return val
}

// 一般的遗传算法里，染色体需要有的方法：
//// Fitness 适应度，越高表示越优解
//func (c Chromosome[T]) Fitness() float64 {
//	panic("implement me")
//}
//
//// RandomInstance 随机实例
//func (c Chromosome[T]) RandomInstance() T {
//	panic("implement me")
//}
//
//// Crossover 交叉
//func (c Chromosome[T]) Crossover(other Chromosome[T]) []Chromosome[T] {
//	panic("implement me")
//}
//
//// Mutate 变异
//func (c Chromosome[T]) Mutate() {
//}
