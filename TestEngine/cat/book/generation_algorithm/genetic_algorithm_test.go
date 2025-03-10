package generation_algorithm

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"math/rand"
	"testing"
	"time"
)

var (
	population = []SimpleEquation{
		{1, 2, "A", 10},
		{1, 2, "B", 20},
		{1, 2, "C", 30},
		{1, 2, "D", 15},
		{1, 2, "E", 1},
		{1, 2, "F", 2},
		{1, 2, "G", 3},
		{1, 2, "H", 4},
	}
	weights = []float64{10, 20, 30, 15, 1, 2, 3, 4}
)

var pop1 = []SimpleEquation{

	{88, 2, "1", -7212},
	{18, 62, "2", -3812},
	{81, 29, "3", -6800},
	{85, 49, "4", -8920},
	{9, 34, "5", -1047},
	{33, 45, "6", -2736},
	{10, 59, "7", -3285},
	{29, 7, "8", -688},
	{87, 71, "9", -11804},
	{63, 7, "10", -3612},
	{55, 16, "11", -2887},
	{94, 3, "12", -8269},
	{48, 6, "13", -2028},
	{26, 32, "14", -1416},
	{73, 69, "15", -9376},
	{76, 12, "16", -5416},
	{8, 40, "17", -1456},
	{95, 19, "18", -8740},
	{14, 62, "19", -3708},
	{26, 54, "20", -3220},
}

func TestPickRoulette(t *testing.T) {
	g := NewGeneticAlgorithm(pop1, 10)
	for i := 0; i < 10; i++ {
		println("====")
		a := g.PickRoulette()
		println(a[0].Name)
		println(a[1].Name)
	}

}

func TestPickTournament(t *testing.T) {
	g := NewGeneticAlgorithm(population, 10)
	num := len(population) / 2
	fmt.Printf("num: %v\n", num)
	a := g.PickTournament(num)
	println(a[0].Name)
	println(a[1].Name)
}

func TestRandIntn(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		//a := rand.Intn(9)
		//b := rand.Intn(9)
		a := r.Float64()
		b := r.Float64()
		fmt.Printf("a: %v\n", a)
		fmt.Printf("b: %v\n", b)
	}
}

func TestDeepCopy(t *testing.T) {
	type Person struct {
		Name    string
		Age     int
		Friends []string
	}
	p1 := Person{
		Name:    "Alice",
		Age:     30,
		Friends: []string{"Bob", "Charlie"},
	}

	p2 := deepcopy.Copy(p1).(Person)
	p2.Friends[0] = "David"                      // 修改副本
	fmt.Println("Original Friends:", p1.Friends) // ["Bob", "Charlie"]
	fmt.Println("Copied Friends:", p2.Friends)   // ["David", "Charlie"]
}

func TestCrossover(t *testing.T) {
	s1 := SimpleEquation{
		x: 1,
		y: 5,
	}

	s2 := SimpleEquation{
		x: 2,
		y: 3,
	}
	s1.Crossover(s2)
	a, b := s1.Crossover(s2)
	println(a.Sprint())
	println(b.Sprint())
}

func TestMutate(t *testing.T) {
	initialPopulation := make([]SimpleEquation, 0)
	var s SimpleEquation
	for i := 0; i < 10; i++ {
		initialPopulation = append(initialPopulation, s.RandomInstance())
	}

	for _, i := range initialPopulation {
		fmt.Printf("x: %v,y: %v fitness: %v\n", i.x, i.y, i.Fitness())
	}
	ga := NewGeneticAlgorithm(initialPopulation, 13)
	ga.Mutate()
	println("==========")
	println(ga.Sprint())
}
