from copy import deepcopy
from random import randrange, random
from typing import List

from demo.chromosome import Chromosome
from demo.genetic_algorithm import GeneticAlgorithm


class SimpleEquation(Chromosome):

    def __init__(self, x: int, y: int) -> None:
        self.x: int = x
        self.y: int = y

    def fitness(self) -> float:
        return 6 * self.x - self.x * self.x + 4 * self.y - self.y * self.y

    # 随机实例
    @classmethod
    def random_instance(cls):
        return SimpleEquation(randrange(100), randrange(100))

    # 交叉
    def crossover(self, other):
        # 深拷贝，确保 child 是原始对象 self 和 other 的独立副本
        # 防止修改 child1 和 child2 后，改变 self 和 other
        child1 = deepcopy(self)
        child2 = deepcopy(other)
        child1.y = other.y
        child2.y = self.y
        return child1, child2

    # 变异
    def mutate(self) -> None:
        if random() > 0.5:
            # mutate x
            if random() > 0.5:
                self.x += 1
            else:
                self.x -= 1
        else:
            # otherwise mutate y
            if random() > 0.5:
                self.y += 1
            else:
                self.y -= 1

    def __str__(self) -> str:
        return f"X: {self.x},Y: {self.y} Fitness: {self.fitness()}"


if __name__ == '__main__':
    initial_population: List[SimpleEquation] = [
        SimpleEquation.random_instance() for _ in range(20)]

    for i in initial_population:
        print(i)
    ga: GeneticAlgorithm[SimpleEquation] = GeneticAlgorithm(
        initial_population=initial_population,
        threshold=13.0,
        max_generations=100,
        mutation_chance=0.1,
        crossover_chance=0.7)
    result = ga.run()
    print(result)
