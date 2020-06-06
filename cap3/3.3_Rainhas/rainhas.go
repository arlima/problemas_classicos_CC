package main

import (
	"fmt"
	"log"
)

type constraint struct {
	variables []int
	columns   []int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *constraint) init(columns []int) {
	c.variables = columns
	c.columns = columns
}

func (c *constraint) satisfied(assignment map[int]int) bool {
	for q1c, q1r := range assignment {
		for q2c := q1c + 1; q2c < len(c.columns)+1; q2c++ {
			if _, ok := assignment[q2c]; ok {
				q2r := assignment[q2c]
				if q1r == q2r { // Same line ?
					return false
				}
				if abs(q1r-q2r) == abs(q1c-q2c) { // Same diagonal ??
					return false
				}
			}
		}
	}
	return true
}

type CSP struct {
	variables   []int
	domains     map[int][]int
	constraints map[int][]constraint
}

func (c *CSP) init(variables []int, domains map[int][]int) {
	c.variables = variables
	c.domains = domains
	c.constraints = make(map[int][]constraint)

	for _, variable := range c.variables {
		if _, ok := c.domains[variable]; !ok {
			log.Fatal("Every variable should have a domain assigned to it")
		}
	}
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (c *CSP) addConstraint(cons constraint) {
	for _, variable := range cons.variables {
		if !intInSlice(variable, c.variables) {
			log.Fatal("Variable in constraint not in CSP")
		} else {
			c.constraints[variable] = append(c.constraints[variable], cons)
		}
	}
}

func (c *CSP) consistent(variable int, assignment map[int]int) bool {
	for _, constraint := range c.constraints[variable] {
		if !constraint.satisfied(assignment) {
			return false
		}
	}
	return true
}

func (c *CSP) backtracking_search(assignment map[int]int) map[int]int {
	if len(assignment) == len(c.variables) {
		return assignment
	}

	unassigned := []int{}
	for _, variable := range c.variables {
		if _, ok := assignment[variable]; !ok {
			unassigned = append(unassigned, variable)
		}
	}

	first := unassigned[0]
	for _, value := range c.domains[first] {
		localAssignment := make(map[int]int)
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
	columns := []int{1, 2, 3, 4, 5, 6, 7, 8}
	rows := make(map[int][]int)
	for _, column := range columns {
		rows[column] = []int{1, 2, 3, 4, 5, 6, 7, 8}
	}
	csp := CSP{}
	csp.init(columns, rows)
	constraint := constraint{}
	constraint.init(columns)
	csp.addConstraint(constraint)
	fmt.Println(csp.backtracking_search(make(map[int]int)))
}
