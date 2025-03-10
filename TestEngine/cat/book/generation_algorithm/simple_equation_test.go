package generation_algorithm

import (
	"fmt"
	"testing"
)

func Test_ReproduceAndReplace(t *testing.T) {
	initialPopulation := make([]SimpleEquation, 0)
	var s SimpleEquation
	for i := 0; i < 20; i++ {
		initialPopulation = append(initialPopulation, s.RandomInstance())
	}

	for _, i := range initialPopulation {
		fmt.Printf("x: %v, y:%v, finess: %v \n", i.x, i.y, i.Fitness())
	}
	ga := NewGeneticAlgorithm(initialPopulation, 13)
	ga.ReproduceAndReplace()
	println("====")
	println(ga.Sprint())
}

func Test_Simple_equation(t *testing.T) {
	initialPopulation := make([]SimpleEquation, 0)
	var s SimpleEquation
	for i := 0; i < 20; i++ {
		initialPopulation = append(initialPopulation, s.RandomInstance())
	}

	for _, i := range initialPopulation {
		fmt.Printf("x: %v, y:%v, finess: %v \n", i.x, i.y, i.Fitness())
	}

	ga := NewGeneticAlgorithm(initialPopulation, 13)
	run := ga.Run()
	println(run.Sprint())
}
