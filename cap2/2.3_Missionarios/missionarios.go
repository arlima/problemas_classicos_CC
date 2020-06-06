package main

import "fmt"

const (
	MAX_NUM int = 3
)

type MCState struct {
	wm   int  // Missionaries in the west margin
	wc   int  // Cannibals in the west margin
	em   int  // Missionaries in the east margin
	ec   int  // Cannibals in the east margin
	boat bool // Is the boat in the west margin ?
}

func (m *MCState) init(missionaries, cannibals int, boat bool) {
	m.wm = missionaries
	m.wc = cannibals
	m.boat = boat
	m.em = MAX_NUM - missionaries
	m.ec = MAX_NUM - cannibals
}

func (m MCState) String() string {
	str := fmt.Sprintf("On the west bank there are %d missionaries and %d cannibals.\n", m.wm, m.wc)
	str += fmt.Sprintf("On the east bank there are %d missionaries and %d cannibals.\n", m.em, m.ec)
	str += fmt.Sprintf("The boat is on the ")
	if m.boat {
		str += fmt.Sprintf("west bank.\n")
	} else {
		str += fmt.Sprintf("east bank.\n")
	}
	return str
}

func (m MCState) isLegal() bool {
	if m.wm < m.wc && m.wm > 0 {
		return false
	}
	if m.em < m.ec && m.em > 0 {
		return false
	}
	return true
}

func (m MCState) goalTest() bool {
	return m.isLegal() && m.em == MAX_NUM && m.ec == MAX_NUM
}

func (m MCState) successors() []MCState {
	sucs := []MCState{}
	ret := []MCState{}

	if m.boat {
		if m.wm > 1 {
			sucs = append(sucs, MCState{m.wm - 2, m.wc, MAX_NUM - (m.wm - 2), MAX_NUM - m.wc, !m.boat})
		}
		if m.wm > 0 {
			sucs = append(sucs, MCState{m.wm - 1, m.wc, MAX_NUM - (m.wm - 1), MAX_NUM - m.wc, !m.boat})
		}
		if m.wc > 1 {
			sucs = append(sucs, MCState{m.wm, m.wc - 2, MAX_NUM - m.wm, MAX_NUM - (m.wc - 2), !m.boat})
		}
		if m.wc > 0 {
			sucs = append(sucs, MCState{m.wm, m.wc - 1, MAX_NUM - m.wm, MAX_NUM - (m.wc - 1), !m.boat})
		}
		if m.wc > 0 && m.wm > 0 {
			sucs = append(sucs, MCState{m.wm - 1, m.wc - 1, MAX_NUM - (m.wm - 1), MAX_NUM - (m.wc - 1), !m.boat})
		}
	} else {
		if m.em > 1 {
			sucs = append(sucs, MCState{m.wm + 2, m.wc, MAX_NUM - (m.wm + 2), MAX_NUM - m.wc, !m.boat})
		}
		if m.em > 0 {
			sucs = append(sucs, MCState{m.wm + 1, m.wc, MAX_NUM - (m.wm + 1), MAX_NUM - m.wc, !m.boat})
		}
		if m.ec > 1 {
			sucs = append(sucs, MCState{m.wm, m.wc + 2, MAX_NUM - m.wm, MAX_NUM - (m.wc + 2), !m.boat})
		}
		if m.ec > 0 {
			sucs = append(sucs, MCState{m.wm, m.wc + 1, MAX_NUM - m.wm, MAX_NUM - (m.wc + 1), !m.boat})
		}
		if m.ec > 0 && m.em > 0 {
			sucs = append(sucs, MCState{m.wm + 1, m.wc + 1, MAX_NUM - (m.wm + 1), MAX_NUM - (m.wc + 1), !m.boat})
		}
	}

	for _, val := range sucs {
		if val.isLegal() {
			ret = append(ret, val)
		}
	}
	return ret
}

type node struct {
	state  MCState
	parent *node
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

func bfs(start MCState) node {
	frontier := queue{}
	explored := make(map[MCState]int)
	frontier.push(node{start, nil})
	explored[start] = 1

	for !frontier.empty() {
		currentNode := frontier.pop()
		currentState := currentNode.state

		if currentState.goalTest() {
			return currentNode
		}

		for _, child := range currentState.successors() {
			if _, ok := explored[child]; ok {
				continue
			}
			explored[child] = 1
			frontier.push(node{child, &currentNode})
		}
	}
	return node{}
}

func nodeToPath(n node) []MCState {
	path := []MCState{n.state}
	for n.parent != nil {
		n = *n.parent
		path = append(path, n.state)
	}
	return path
}

func reverse(nodes []MCState) []MCState {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}

func displaySolution(path []MCState) {
	if len(path) == 0 {
		return
	}
	oldState := path[0]
	fmt.Println(oldState)

	for _, currentState := range path {
		if currentState.boat {
			fmt.Printf("%d missionaries and %d cannibals moved from the east bank to the west bank.\n", oldState.em-currentState.em, oldState.ec-currentState.ec)
		} else {
			fmt.Printf("%d missionaries and %d cannibals moved from the west bank to the east bank.\n", oldState.wm-currentState.wm, oldState.wc-currentState.wc)
		}
		fmt.Println(currentState)
		oldState = currentState
	}
}

func main() {
	MCState := MCState{}
	MCState.init(MAX_NUM, MAX_NUM, true)
	fmt.Println(MCState)
	result := bfs(MCState)
	displaySolution(reverse(nodeToPath(result)))
}
