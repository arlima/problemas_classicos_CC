package main

import "fmt"

type stack struct {
	container []int
}

func (s *stack) push(item int) {
	s.container = append(s.container, item)
}

func (s *stack) pop() int {
	p := len(s.container) - 1
	value := s.container[p]
	s.container = s.container[:p]
	return value
}

func hanoi(begin *stack, end *stack, temp *stack, n int) {
	if n == 1 {
		end.push(begin.pop())
	} else {
		hanoi(begin, temp, end, n-1)
		hanoi(begin, end, temp, 1)
		hanoi(temp, end, begin, n-1)
	}
}

func main() {
	towerA := stack{}
	towerB := stack{}
	towerC := stack{}
	for i := 1; i <= 10; i++ {
		towerA.push(i)
	}
	fmt.Println("------ BEFORE -----------")
	fmt.Println("Tower A: ", towerA.container)
	fmt.Println("Tower B: ", towerB.container)
	fmt.Println("Tower C: ", towerC.container)
	hanoi(&towerA, &towerC, &towerB, 10)
	fmt.Println("------ AFTER -----------")
	fmt.Println("Tower A: ", towerA.container)
	fmt.Println("Tower B: ", towerB.container)
	fmt.Println("Tower C: ", towerC.container)
}
