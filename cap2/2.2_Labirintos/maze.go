package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
)

const (
	EMPTY   byte = ' '
	BLOCKED byte = 'X'
	START   byte = 'S'
	GOAL    byte = 'G'
	PATH    byte = '*'
)

type mazeLocation struct {
	column int
	row    int
}

type node struct {
	state     mazeLocation
	parent    *node
	cost      float64
	heuristic float64
}

type PriorityQueue []node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) empty() bool { return len(pq) == 0 }

func (pq PriorityQueue) Less(i, j int) bool {
	return (pq[i].cost + pq[i].heuristic) < (pq[j].cost + pq[j].heuristic)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = node{} // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

type stack struct {
	container []node
}

func (s stack) empty() bool {
	return (len(s.container) == 0)
}

func (s *stack) push(n node) {
	s.container = append(s.container, n)
}

func (s *stack) pop() node {
	p := len(s.container) - 1
	value := s.container[p]
	s.container = s.container[:p]
	return value
}

type queue struct {
	container []node
}

func (q queue) empty() bool {
	return (len(q.container) == 0)
}

func (q *queue) push(n node) {
	q.container = append(q.container, n)
}

func (q *queue) pop() node {
	value := q.container[0]
	q.container = q.container[1:]
	return value
}

type maze struct {
	rows    int
	columns int
	grid    []byte
	start   mazeLocation
	goal    mazeLocation
}

func (m *maze) init(r int, c int, s float64, startPos mazeLocation, goalPos mazeLocation) {
	m.grid = make([]byte, r*c)
	m.rows = r
	m.columns = c
	m.goal = goalPos
	m.start = startPos

	rnd := rand.New(rand.NewSource(36))

	for x := 0; x < c; x++ {
		for y := 0; y < r; y++ {
			if rnd.Float64() < s {
				m.grid[x+y*r] = BLOCKED
			} else {
				m.grid[x+y*r] = EMPTY
			}
		}
	}
	m.grid[startPos.row+startPos.column*r] = START
	m.grid[goalPos.row+goalPos.column*r] = GOAL
}

func (m maze) print() {
	for y := m.rows - 1; y >= 0; y-- {
		for x := 0; x < m.columns; x++ {
			fmt.Printf("%c", m.grid[x+y*m.rows])
		}
		fmt.Printf("\n")
	}
}

func (m maze) goalTest(ml mazeLocation) bool {
	return (ml == m.goal)
}

func (m maze) successors(ml mazeLocation) []mazeLocation {
	locations := []mazeLocation{}
	if ml.row+1 < m.rows && m.grid[ml.column+(ml.row+1)*m.rows] != BLOCKED {
		locations = append(locations, mazeLocation{ml.column, ml.row + 1})
	}
	if ml.row-1 >= 0 && m.grid[ml.column+(ml.row-1)*m.rows] != BLOCKED {
		locations = append(locations, mazeLocation{ml.column, ml.row - 1})
	}
	if ml.column+1 < m.columns && m.grid[(ml.column+1)+ml.row*m.rows] != BLOCKED {
		locations = append(locations, mazeLocation{ml.column + 1, ml.row})
	}
	if ml.column-1 >= 0 && m.grid[(ml.column-1)+ml.row*m.rows] != BLOCKED {
		locations = append(locations, mazeLocation{ml.column - 1, ml.row})
	}
	return locations
}

func (m *maze) mark(path []mazeLocation) {
	for _, val := range path {
		m.grid[val.column+val.row*m.rows] = PATH
	}
	m.grid[m.start.column+m.start.row*m.rows] = START
	m.grid[m.goal.column+m.goal.row*m.rows] = GOAL
}

func (m *maze) clear(path []mazeLocation) {
	for _, val := range path {
		m.grid[val.column+val.row*m.rows] = EMPTY
	}
	m.grid[m.start.column+m.start.row*m.rows] = START
	m.grid[m.goal.column+m.goal.row*m.rows] = GOAL
}

func dfs(m maze) node {
	initial := m.start
	frontier := stack{}
	explored := make(map[mazeLocation]int)
	frontier.push(node{initial, nil, 0, 0})
	explored[initial] = 1

	for !frontier.empty() {
		currentNode := frontier.pop()
		currentState := currentNode.state

		if m.goalTest(currentState) {
			return currentNode
		}

		for _, child := range m.successors(currentState) {
			if _, ok := explored[child]; ok {
				continue
			}
			explored[child] = 1
			frontier.push(node{child, &currentNode, 0, 0})
		}
	}
	return node{}
}

func nodeToPath(n node) []mazeLocation {
	path := []mazeLocation{n.state}
	for n.parent != nil {
		n = *n.parent
		path = append(path, n.state)
	}
	return path
}

func bfs(m maze) node {
	initial := m.start
	frontier := queue{}
	explored := make(map[mazeLocation]int)
	frontier.push(node{initial, nil, 0, 0})
	explored[initial] = 1

	for !frontier.empty() {
		currentNode := frontier.pop()
		currentState := currentNode.state

		if m.goalTest(currentState) {
			return currentNode
		}

		for _, child := range m.successors(currentState) {
			if _, ok := explored[child]; ok {
				continue
			}
			explored[child] = 1
			frontier.push(node{child, &currentNode, 0, 0})
		}
	}
	return node{}
}

func manhattan(init, goal mazeLocation) float64 {
	xdist := math.Abs(float64(init.column) - float64(goal.column))
	ydist := math.Abs(float64(init.row) - float64(goal.row))
	return xdist + ydist
}

type heuristicFn func(pos, goal mazeLocation) float64

func astar(m maze, heuristic heuristicFn) node {
	initial := m.start
	frontier := PriorityQueue{}
	explored := make(map[mazeLocation]float64)
	heap.Init(&frontier)

	heap.Push(&frontier, node{initial, nil, 0, heuristic(initial, m.goal)})

	explored[initial] = 0.0

	for !frontier.empty() {
		currentNode := heap.Pop(&frontier).(node)
		currentState := currentNode.state

		if m.goalTest(currentState) {
			return currentNode
		}

		for _, child := range m.successors(currentState) {
			newCost := currentNode.cost + 1.0
			_, ok := explored[child]
			if !ok || explored[child] > newCost {
				explored[child] = newCost
				heap.Push(&frontier, node{child, &currentNode, newCost, heuristic(child, m.goal)})
			}
		}
	}
	return node{}
}

func main() {
	maze := maze{}
	maze.init(10, 10, 0.2, mazeLocation{0, 0}, mazeLocation{9, 9})
	maze.print()
	solution := dfs(maze)
	path := nodeToPath(solution)
	maze.mark(path)
	fmt.Println("-------------------")
	maze.print()
	maze.clear(path)
	fmt.Println("-------------------")
	solution2 := bfs(maze)
	path2 := nodeToPath(solution2)
	maze.mark(path2)
	maze.print()
	maze.clear(path2)
	fmt.Println("-------------------")
	solution3 := astar(maze, manhattan)
	path3 := nodeToPath(solution3)
	maze.mark(path3)
	maze.print()
}
