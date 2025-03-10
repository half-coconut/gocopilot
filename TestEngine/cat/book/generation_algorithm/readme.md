踩的坑：

1. 在 Go 语言中，使用 for...range 循环遍历切片、数组或映射时，循环变量只是元素的一个副本，而不是元素本身的引用。如果需要修改元素的值，需要使用索引访问元素。

```go
func (g GeneticAlgorithm) Mutate() {
    for i := range g.Population {
        if GetRandFloat64() < g.MutationChance {
            g.Population[i] = g.Population[i].Mutate()
        }
    }
}

```

2. Crossover 时，child1 和 child2 的 x 值总是 0, =_=

```go
// Crossover 交叉
func (s SimpleEquation) Crossover(other SimpleEquation) (SimpleEquation, SimpleEquation) {
    // 这里需要深拷贝
    child1 := deepcopy.Copy(s).(SimpleEquation)
    child2 := deepcopy.Copy(other).(SimpleEquation)
    child1.x = s.x // 明确再赋值一遍 =_=
    child1.y = other.y
    child2.x = other.x // 明确再赋值一遍 =_=
    child2.y = s.y
    return child1, child2
}

```

3. FindBest 加锁，重要！

```go
    mu.Lock()
    defer mu.Unlock()
```

4. 使用指针 重要 g.Population 通过指针修改原值，在不同的方法间调用，尤其是 findBest
5. 在 ReproduceAndReplace 中，加入 best 优先，可以减少执行代际的次数