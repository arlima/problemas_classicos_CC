package main

import (
	"fmt"
	"log"
)

type constraint struct {
	variables []string
	letters   []string
}

func (c *constraint) init(letters []string) {
	c.variables = letters
	c.letters = letters
}

func (c *constraint) satisfied(assignment map[string]int) bool {
	a := make(map[int]bool)
	for _, value := range assignment {
		a[value] = true
	}
	if len(a) < len(assignment) {
		return false
	}
	if len(assignment) == len(c.letters) {
		s := assignment["S"]
		e := assignment["E"]
		n := assignment["N"]
		d := assignment["D"]
		m := assignment["M"]
		o := assignment["O"]
		r := assignment["R"]
		y := assignment["Y"]
		send := s*1000 + e*100 + n*10 + d
		more := m*1000 + o*100 + r*10 + e
		money := m*10000 + o*1000 + n*100 + e*10 + y
		return send+more == money
	}
	return true
}

type CSP struct {
	variables   []string
	domains     map[string][]int
	constraints map[string][]constraint
}

func (c *CSP) init(variables []string, domains map[string][]int) {
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

func (c *CSP) consistent(variable string, assignment map[string]int) bool {
	for _, constraint := range c.constraints[variable] {
		if !constraint.satisfied(assignment) {
			return false
		}
	}
	return true
}

func (c *CSP) backtracking_search(assignment map[string]int) map[string]int {
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
		localAssignment := make(map[string]int)
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

func main() {
	letters := []string{"S", "E", "N", "D", "M", "O", "R", "Y"}
	possibleDigits := make(map[string][]int)
	for _, letter := range letters {
		possibleDigits[letter] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	}
	possibleDigits["M"] = []int{1}
	csp := CSP{}
	csp.init(letters, possibleDigits)
	constraint := constraint{}
	constraint.init(letters)
	csp.addConstraint(constraint)
	fmt.Println(csp.backtracking_search(make(map[string]int)))
}
