package main

import (
	"fmt"
	"log"
)

type constraint struct {
	variables []string
	place1    string
	place2    string
}

func (c *constraint) init(place1, place2 string) {
	c.variables = []string{place1, place2}
	c.place1 = place1
	c.place2 = place2
}

func (c *constraint) satisfied(assignment map[string]string) bool {
	_, ok1 := assignment[c.place1]
	_, ok2 := assignment[c.place2]
	if !ok1 || !ok2 {
		return true
	}
	return assignment[c.place1] != assignment[c.place2]
}

type CSP struct {
	variables   []string
	domains     map[string][]string
	constraints map[string][]constraint
}

func (c *CSP) init(variables []string, domains map[string][]string) {
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
			fmt.Println(variable)
			log.Fatal("Variable in constraint not in CSP")
		} else {
			c.constraints[variable] = append(c.constraints[variable], cons)
		}
	}
}

func (c *CSP) consistent(variable string, assignment map[string]string) bool {
	for _, constraint := range c.constraints[variable] {
		if !constraint.satisfied(assignment) {
			return false
		}
	}
	return true
}

func (c *CSP) backtracking_search(assignment map[string]string) map[string]string {
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
		localAssignment := make(map[string]string)
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
	variables := []string{"Western Australia", "Northern Territory", "South Australia", "Queensland", "New South Wales", "Victoria", "Tasmania"}
	domains := make(map[string][]string)
	for _, variable := range variables {
		domains[variable] = []string{"red", "green", "blue"}
	}
	csp := CSP{}
	csp.init(variables, domains)
	constraint := constraint{}
	constraint.init("Western Australia", "Northern Territory")
	csp.addConstraint(constraint)
	constraint.init("Western Australia", "South Australia")
	csp.addConstraint(constraint)
	constraint.init("South Australia", "Northern Territory")
	csp.addConstraint(constraint)
	constraint.init("Queensland", "Northern Territory")
	csp.addConstraint(constraint)
	constraint.init("Queensland", "South Australia")
	csp.addConstraint(constraint)
	constraint.init("Queensland", "New South Wales")
	csp.addConstraint(constraint)
	constraint.init("New South Wales", "South Australia")
	csp.addConstraint(constraint)
	constraint.init("Victoria", "South Australia")
	csp.addConstraint(constraint)
	constraint.init("Victoria", "New South Wales")
	csp.addConstraint(constraint)
	constraint.init("Victoria", "Tasmania")
	csp.addConstraint(constraint)

	fmt.Println(csp.backtracking_search(make(map[string]string)))
}
