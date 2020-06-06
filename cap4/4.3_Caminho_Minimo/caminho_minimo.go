package main

import "fmt"

type vertex struct {
	value string
}

type edge struct {
	u int
	v int
}

func (e *edge) init(vu, vv int) {
	e.u = vu
	e.v = vv
}

func (e edge) reversed() edge {
	return edge{e.v, e.u}
}

func (e edge) String() string {
	return fmt.Sprintf("{%v} -> {%v}", e.u, e.v)
}

type graph struct {
	vertices []vertex
	edges    [][]edge
}

func (g *graph) init(v []vertex) {
	g.vertices = v
	for range v {
		g.edges = append(g.edges, []edge{})
	}
}

func (g *graph) edgeCount() int {
	s := 0
	for _, edges := range g.edges {
		s += len(edges)
	}
	return s
}

func (g *graph) vertexCount() int {
	return len(g.vertices)
}

func (g *graph) addVertex(v vertex) int {
	g.vertices = append(g.vertices, v)
	g.edges = append(g.edges, []edge{})
	return len(g.vertices) - 1
}

func (g *graph) addEdge(e edge) {
	g.edges[e.u] = append(g.edges[e.u], e)
	g.edges[e.v] = append(g.edges[e.v], e.reversed())
}

func (g *graph) addEdgeByIndices(vu, vv int) {
	e := edge{vu, vv}
	g.addEdge(e)
}

func index(list []vertex, value vertex) int {
	for k, v := range list {
		if v == value {
			return k
		}
	}
	return 0
}

func (g graph) addEdgeByVertices(first, second vertex) {
	u := index(g.vertices, first)
	v := index(g.vertices, second)
	g.addEdgeByIndices(u, v)
}

func (g graph) vertexAt(index int) vertex {
	return g.vertices[index]
}

func (g graph) indexOf(v vertex) int {
	return index(g.vertices, v)
}

func (g graph) neighborsForIndex(index int) []vertex {
	ret := []vertex{}
	for _, e := range g.edges[index] {
		ret = append(ret, g.vertexAt(e.v))
	}
	return ret
}

func (g graph) neighborsForVertex(v vertex) []vertex {
	return g.neighborsForIndex(g.indexOf(v))
}

func (g graph) edgesForIndex(index int) []edge {
	return g.edges[index]
}

func (g graph) edgesForVertex(v vertex) []edge {
	return g.edgesForIndex(g.indexOf(v))
}

func (g graph) String() string {
	str := ""
	for i := 0; i < g.vertexCount(); i++ {
		str += fmt.Sprintf("%v -> %v\n", g.vertexAt(i), g.neighborsForIndex(i))
	}
	return str
}

type node struct {
	state  vertex
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

func bfs(initial vertex, final vertex, g graph) node {
	frontier := queue{}
	explored := make(map[vertex]int)
	frontier.push(node{initial, nil})
	explored[initial] = 1

	for !frontier.empty() {
		currentNode := frontier.pop()
		currentState := currentNode.state

		if currentState == final {
			return currentNode
		}

		for _, child := range g.neighborsForVertex(currentState) {
			if _, ok := explored[child]; ok {
				continue
			}
			explored[child] = 1
			frontier.push(node{child, &currentNode})
		}
	}
	return node{}
}

func reverse(g []vertex) []vertex {
	for i, j := 0, len(g)-1; i < j; i, j = i+1, j-1 {
		g[i], g[j] = g[j], g[i]
	}
	return g
}

func nodeToPath(n node) []vertex {
	path := []vertex{n.state}
	for n.parent != nil {
		n = *n.parent
		path = append(path, n.state)
	}
	return reverse(path)
}

func main() {
	citiGraph := graph{}
	citiGraph.init([]vertex{{"Seattle"}, {"San Francisco"}, {"Los Angeles"}, {"Riverside"}, {"Phoenix"}, {"Chicago"},
		{"Boston"}, {"New York"}, {"Atlanta"}, {"Miami"}, {"Dallas"}, {"Houston"}, {"Detroit"},
		{"Philadelphia"}, {"Washington"}})
	citiGraph.addEdgeByVertices(vertex{"Seattle"}, vertex{"Chicago"})
	citiGraph.addEdgeByVertices(vertex{"Seattle"}, vertex{"San Francisco"})
	citiGraph.addEdgeByVertices(vertex{"San Francisco"}, vertex{"Riverside"})
	citiGraph.addEdgeByVertices(vertex{"San Francisco"}, vertex{"Los Angeles"})
	citiGraph.addEdgeByVertices(vertex{"Los Angeles"}, vertex{"Riverside"})
	citiGraph.addEdgeByVertices(vertex{"Los Angeles"}, vertex{"Phoenix"})
	citiGraph.addEdgeByVertices(vertex{"Riverside"}, vertex{"Phoenix"})
	citiGraph.addEdgeByVertices(vertex{"Riverside"}, vertex{"Chicago"})
	citiGraph.addEdgeByVertices(vertex{"Phoenix"}, vertex{"Dallas"})
	citiGraph.addEdgeByVertices(vertex{"Phoenix"}, vertex{"Houston"})
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Chicago"})
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Atlanta"})
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Houston"})
	citiGraph.addEdgeByVertices(vertex{"Houston"}, vertex{"Atlanta"})
	citiGraph.addEdgeByVertices(vertex{"Houston"}, vertex{"Miami"})
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Chicago"})
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Washington"})
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Miami"})
	citiGraph.addEdgeByVertices(vertex{"Miami"}, vertex{"Washington"})
	citiGraph.addEdgeByVertices(vertex{"Chicago"}, vertex{"Detroit"})
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"Boston"})
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"Washington"})
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"New York"})
	citiGraph.addEdgeByVertices(vertex{"Boston"}, vertex{"New York"})
	citiGraph.addEdgeByVertices(vertex{"New York"}, vertex{"Philadelphia"})
	citiGraph.addEdgeByVertices(vertex{"Philadelphia"}, vertex{"Washington"})

	bfsResult := bfs(vertex{"Boston"}, vertex{"Miami"}, citiGraph)

	fmt.Println(nodeToPath(bfsResult))
}
