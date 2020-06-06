package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

type chromosome struct {
	x int
	y int
}

func (c chromosome) fitness() float64 {
	//
	return 0
}

func (c chromosome) crossover(other chromosome) []chromosome {
	//
	return nil
}

func (c *chromosome) mutate() {
	//
}

func choicesRoulette(population []chromosome, wheel []float64, qtt int) []chromosome {
	var k int
	var s float64 = 0.0

	ret := []chromosome{}
	if len(population) != len(wheel) {
		log.Fatal("Population size is not equal to wheel size !")
	}

	for i := 0; i < len(population); i++ {
		s += wheel[i]
	}

	for i := 0; i < qtt; i++ {
		r := rand.Float64()
		k = 0
		for {
			r -= (wheel[k] / s)
			if r <= 0 {
				break
			}
			k++
		}
		ret = append(ret, population[k])
	}
	return ret
}

type chromosomeFit struct {
	c chromosome
	f float64
}

type cf []chromosomeFit

func (c cf) Len() int           { return len(c) }
func (c cf) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c cf) Less(i, j int) bool { return c[i].f < c[j].f }

func choicesTournament(population []chromosome, participants int, qtt int) []chromosome {
	part := cf{}
	for i := 0; i < participants; i++ {
		p := population[rand.Intn(len(population))]
		part = append(part, chromosomeFit{p, p.fitness()})
	}
	sort.Reverse(part)
	ret := []chromosome{}

	for i := 0; i < qtt; i++ {
		ret = append(ret, part[i].c)
	}
	return ret
}

type geneticAlgorithm struct {
	population      []chromosome
	threshold       float64
	maxGenerations  int
	mutationChance  float64
	crossoverChance float64
	fitnessKey      func() float64
	selectionType   string
}

func (g geneticAlgorithm) init(initialPopulation []chromosome, threshold float64, maxGenerations int, mutationChance float64, crossoverChance float64) {
	g.population = initialPopulation
	g.threshold = threshold
	g.maxGenerations = maxGenerations
	g.mutationChance = mutationChance
	g.crossoverChance = crossoverChance
	g.fitnessKey = g.population[0].fitness
}

func (g geneticAlgorithm) pickRoulette(wheel []float64) []chromosome {
	return choicesRoulette(g.population, wheel, 2)
}

func (g geneticAlgorithm) pickTournament(participants int) []chromosome {
	return choicesTournament(g.population, participants, 2)
}

func (g *geneticAlgorithm) reproduceAndReplace() {
	newPopulation := []chromosome{}
	w := []float64{}
	parents := []chromosome{}
	for i := 0; i < len(g.population); i++ {
		w = append(w, g.population[i].fitness())
	}
	for len(newPopulation) < len(g.population) {
		if g.selectionType == "ROULETTE" {
			parents = g.pickRoulette(w)
		} else {
			parents = g.pickTournament(len(g.population) / 2)
		}
		if rand.Float64() < g.crossoverChance {
			co := parents[0].crossover(parents[1])
			newPopulation = append(newPopulation, co[0])
			newPopulation = append(newPopulation, co[1])
		} else {
			newPopulation = append(newPopulation, parents[0])
			newPopulation = append(newPopulation, parents[1])
		}
	}
	if len(newPopulation) > len(g.population) {
		newPopulation = newPopulation[:len(newPopulation)-1]
	}
	g.population = newPopulation
}

func (g *geneticAlgorithm) mutate() {
	for _, individual := range g.population {
		if rand.Float64() < g.mutationChance {
			individual.mutate()
		}
	}
}

func (g *geneticAlgorithm) maxFitness() chromosome {
	var max float64
	var ind chromosome
	for _, individual := range g.population {
		f := individual.fitness()
		if f > max {
			max = f
			ind = individual
		}
	}
	return ind
}

func (g geneticAlgorithm) run() chromosome {
	best := g.maxFitness()
	for generation := 0; generation < g.maxGenerations; generation++ {
		if best.fitness() >= g.threshold {
			return best
		}
		fmt.Println("Generation %d, Best %f, Avg %f", generation, best.fitness, 0.0)
		g.reproduceAndReplace()
		g.mutate()
		highest := g.maxFitness()
		if highest.fitness() > best.fitness() {
			best = highest
		}
	}
	return best
}

func main() {
	//
}
