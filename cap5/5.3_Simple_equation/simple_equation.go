package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

type chromosome struct {
	x int
	y int
}

func (c *chromosome) init(xi int, yi int) {
	c.x = xi
	c.y = yi
}

func (c chromosome) randomInstance() chromosome {
	return chromosome{rand.Intn(100), rand.Intn(100)}
}

func (c chromosome) fitness() float64 {
	return 6.0*float64(c.x) - float64(c.x)*float64(c.x) + 4.0*float64(c.y) - float64(c.y)*float64(c.y)
}

func (c chromosome) crossover(other chromosome) []chromosome {
	child1 := c
	child2 := other
	child1.y = other.y
	child2.y = c.y
	return []chromosome{child1, child2}
}

func (c *chromosome) mutate() {
	if rand.Float64() > 0.5 {
		if rand.Float64() > 0.5 {
			c.x++
		} else {
			c.x--
		}
	} else {
		if rand.Float64() > 0.5 {
			c.y++
		} else {
			c.y--
		}
	}
}

func (c chromosome) String() string {
	return fmt.Sprintf("X: %d, Y: %d, Fitness: %f\n", c.x, c.y, c.fitness())
}

func (c chromosome) copyChromossome(dst *chromosome) {
	dst.x = c.x
	dst.y = c.y
}

type chromosomeFit struct {
	c chromosome
	f float64
}

type cf []chromosomeFit

func (c cf) Len() int           { return len(c) }
func (c cf) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c cf) Less(i, j int) bool { return c[i].f > c[j].f }

func choicesTournament(population []chromosome, participants int, qtt int) []chromosome {
	part := cf{}
	for i := 0; i < participants; i++ {
		p := population[rand.Intn(len(population))]
		part = append(part, chromosomeFit{p, p.fitness()})
	}
	sort.Sort(part)
	ret := []chromosome{}

	for i := 0; i < qtt; i++ {
		r := chromosome{}.randomInstance()
		part[i].c.copyChromossome(&r)
		ret = append(ret, r)
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

func (g *geneticAlgorithm) init(initialPopulation []chromosome, threshold float64, maxGenerations int, mutationChance float64, crossoverChance float64, selectionType string) {
	g.population = initialPopulation
	g.threshold = threshold
	g.maxGenerations = maxGenerations
	g.mutationChance = mutationChance
	g.crossoverChance = crossoverChance
	g.fitnessKey = g.population[0].fitness
	g.selectionType = selectionType
}

// RouletteDrawN draws n numbers randomly from a probability mass function (PMF) defined by weights in p.
// RouletteDrawN implements the Roulette Wheel Draw a.k.a. Fitness Proportionate Selection:
// - https://en.wikipedia.org/wiki/Fitness_proportionate_selection
// - http://www.keithschwarz.com/darts-dice-coins/
// It returns a slice of n indices into the vector p.
// It fails with error if p is empty or nil.
func RouletteDrawN(p []float64, n int) ([]int, error) {
	if p == nil || len(p) == 0 {
		return nil, fmt.Errorf("Invalid probability weights: %v", p)
	}
	// Initialization: create the discrete CDF
	// We know that cdf is sorted in ascending order
	cdf := make([]float64, len(p))
	floats.CumSum(cdf, p)
	// Generation:
	// 1. Generate a uniformly-random value x in the range [0,1)
	// 2. Using a binary search, find the index of the smallest element in cdf larger than x
	var val float64
	indices := make([]int, n)
	for i := range indices {
		// multiply the sample with the largest CDF value; easier than normalizing to [0,1)
		val = distuv.UnitUniform.Rand() * cdf[len(cdf)-1]
		// Search returns the smallest index i such that cdf[i] > val
		indices[i] = sort.Search(len(cdf), func(i int) bool { return cdf[i] > val })
	}

	return indices, nil
}

func (g geneticAlgorithm) pickRoulette(wheel []float64) []chromosome {
	nw := []float64{}
	min := math.MaxFloat64
	for _, v := range wheel {
		if v < min {
			min = v
		}
	}

	for k := 0; k < len(wheel); k++ {
		nw = append(nw, wheel[k]-min+1.0)
	}

	indices, _ := RouletteDrawN(nw, 2)

	ret := []chromosome{}
	for i := 0; i < 2; i++ {
		r := chromosome{}.randomInstance()
		g.population[indices[i]].copyChromossome(&r)
		ret = append(ret, r)
	}
	return ret
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
	if len(newPopulation) != len(g.population) {
		log.Fatal("Problem. New population has different size.")
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
	max := g.population[0].fitness()
	ind := g.population[0]
	for _, individual := range g.population {
		f := individual.fitness()
		if f > max {
			max = f
			ind = individual
		}
	}
	return ind
}

func (g *geneticAlgorithm) avgFitness() float64 {
	var sum float64
	for _, individual := range g.population {
		sum += individual.fitness()
	}
	return sum / float64(len(g.population))
}

func (g *geneticAlgorithm) run() chromosome {
	best := chromosome{}.randomInstance()
	highest := g.maxFitness()
	highest.copyChromossome(&best)
	for generation := 0; generation < g.maxGenerations; generation++ {
		if best.fitness() >= g.threshold {
			return best
		}
		fmt.Printf("Generation %d, Best %f, Avg %f\n", generation, best.fitness(), g.avgFitness())
		g.reproduceAndReplace()
		g.mutate()
		highest = g.maxFitness()
		if highest.fitness() > best.fitness() {
			highest.copyChromossome(&best)
		}
	}
	return best
}

func main() {
	rand.Seed(100)
	initialPopulation := []chromosome{}
	for i := 0; i < 20; i++ {
		initialPopulation = append(initialPopulation, chromosome{}.randomInstance())
	}

	ga := geneticAlgorithm{}
	ga.init(initialPopulation, 13.0, 1000, 0.1, 0.7, "TOURNAMENT")
	result := ga.run()
	fmt.Println(result)
}
