package main

import (
	"fmt"
	"log"
	"math/rand"
)

type grid [][]string
type gridLocation struct {
	col int
	row int
}

func generateGrid(cols, rows int) grid {
	g := grid{}
	for x := 0; x < cols; x++ {
		g = append(g, []string{})
		for y := 0; y < rows; y++ {
			c := byte(rand.Intn(24) + 'a')
			g[x] = append(g[x], string(c))
		}
	}
	return g
}

func printGrid(g grid, cols, rows int) {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			fmt.Printf(g[x][y])
		}
		fmt.Printf("\n")
	}
}

func generateDomain(word string, g grid) [][]gridLocation {
	domain := [][]gridLocation{}
	width := len(g)
	height := len(g[0])
	length := len(word)

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if col+length <= width {
				d := []gridLocation{}
				for c := col; c < col+length; c++ {
					d = append(d, gridLocation{c, row})
				}
				domain = append(domain, d)

				if row+length <= height {
					d := []gridLocation{}
					for r := row; r < row+length; r++ {
						d = append(d, gridLocation{col + (r - row), r})
					}
					domain = append(domain, d)
				}
			}

			if row+length <= height {
				d := []gridLocation{}
				for r := row; r < row+length; r++ {
					d = append(d, gridLocation{col, r})
				}
				domain = append(domain, d)

				if col-length >= 0 {
					d := []gridLocation{}
					for r := row; r < row+length; r++ {
						d = append(d, gridLocation{col - (r - row), r})
					}
					domain = append(domain, d)
				}
			}
		}
	}
	return domain
}

type constraint struct {
	variables []string
	words     []string
}

func (c *constraint) init(words []string) {
	c.variables = words
	c.words = words
}

func unique(array []gridLocation) []gridLocation {
	keys := make(map[gridLocation]bool)
	list := []gridLocation{}
	for _, entry := range array {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (c *constraint) satisfied(assignment map[string][]gridLocation) bool {
	allLocations := []gridLocation{}
	for _, values := range assignment {
		for _, locs := range values {
			allLocations = append(allLocations, locs)
		}
	}
	return len(unique(allLocations)) == len(allLocations)
}

type CSP struct {
	variables   []string
	domains     map[string][][]gridLocation
	constraints map[string][]constraint
}

func (c *CSP) init(variables []string, domains map[string][][]gridLocation) {
	c.variables = variables
	c.domains = domains
	c.constraints = make(map[string][]constraint)

	for _, variable := range c.variables {
		if _, ok := c.domains[variable]; !ok {
			log.Fatal("Every variable should have a domain assigned to it")
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (c *CSP) addConstraint(cons constraint) {
	for _, variable := range cons.variables {
		if !stringInSlice(variable, c.variables) {
			log.Fatal("Variable in constraint not in CSP")
		} else {
			c.constraints[variable] = append(c.constraints[variable], cons)
		}
	}
}

func (c *CSP) consistent(variable string, assignment map[string][]gridLocation) bool {
	for _, constraint := range c.constraints[variable] {
		if !constraint.satisfied(assignment) {
			return false
		}
	}
	return true
}

func (c *CSP) backtracking_search(assignment map[string][]gridLocation) map[string][]gridLocation {
	if len(assignment) == len(c.variables) {
		return assignment
	}

	unassigned := []string{}
	for _, variable := range c.variables {
		if _, ok := assignment[variable]; !ok {
			unassigned = append(unassigned, variable)
		}
	}

	first := unassigned[0]
	for _, value := range c.domains[first] {
		localAssignment := make(map[string][]gridLocation)
		for k, v := range assignment {
			localAssignment[k] = v
		}
		localAssignment[first] = value
		if c.consistent(first, localAssignment) {
			result := c.backtracking_search(localAssignment)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

func reverse(g []gridLocation) []gridLocation {
	for i, j := 0, len(g)-1; i < j; i, j = i+1, j-1 {
		g[i], g[j] = g[j], g[i]
	}
	return g
}

func main() {
	rand.Seed(86)
	g := generateGrid(9, 9)
	printGrid(g, 9, 9)

	words := []string{"MATTHEW", "JOE", "MARY", "SARAH", "SALLY", "ADRIANO", "THATIANA", "GABRIEL", "LEONARDO"}
	locations := make(map[string][][]gridLocation)

	for _, word := range words {
		locations[word] = generateDomain(word, g)
	}
	csp := CSP{}
	csp.init(words, locations)
	constraint := constraint{}
	constraint.init(words)
	csp.addConstraint(constraint)
	solution := csp.backtracking_search(make(map[string][]gridLocation))
	if solution != nil {
		for word, gridLocations := range solution {
			if rand.Float32() < 0.5 {
				gridLocations = reverse(gridLocations)
			}
			for index, letter := range word {
				row, col := gridLocations[index].row, gridLocations[index].col
				g[col][row] = string(letter)
			}
		}
		fmt.Print("\n\n")
		printGrid(g, 9, 9)
	} else {
		fmt.Println("No solution found!")
	}
}
