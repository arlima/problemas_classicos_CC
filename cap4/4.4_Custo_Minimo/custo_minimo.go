package main

import (
	"container/heap"
	"fmt"
)

type vertex struct {
	value string
}

type edge struct {
	u      int
	v      int
	weight float32
}

func (e *edge) init(vu, vv int) {
	e.u = vu
	e.v = vv
}

func (e edge) reversed() edge {
	return edge{e.v, e.u, e.weight}
}

func (e edge) String() string {
	return fmt.Sprintf("%d (%f) -> %d", e.u, e.v, e.weight)
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

func (g *graph) addEdgeByIndices(vu, vv int, w float32) {
	e := edge{vu, vv, w}
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

func (g graph) addEdgeByVertices(first, second vertex, w float32) {
	u := index(g.vertices, first)
	v := index(g.vertices, second)
	g.addEdgeByIndices(u, v, w)
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

type vertexWithWeight struct {
	v vertex
	w float32
}

func (g graph) neighborsForIndexWithWeights(index int) []vertexWithWeight {
	ret := []vertexWithWeight{}
	for _, e := range g.edges[index] {
		ret = append(ret, vertexWithWeight{g.vertexAt(e.v), e.weight})
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
		str += fmt.Sprintf("%v -> %v\n", g.vertexAt(i), g.neighborsForIndexWithWeights(i))
	}
	return str
}

type PriorityQueue []edge

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) empty() bool { return len(pq) == 0 }

func (pq PriorityQueue) Less(i, j int) bool {
	return (pq[i].weight < pq[j].weight)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(edge)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = edge{} // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func totalWeight(path []edge) float32 {
	total := float32(0.0)
	for _, e := range path {
		total += e.weight
	}
	return total
}

func mst(wg graph, start int) []edge {
	if start > (wg.vertexCount()-1) || start < 0 {
		return []edge{}
	}
	result := []edge{}
	pq := PriorityQueue{}
	heap.Init(&pq)
	visited := make([]bool, wg.vertexCount())

	visit := func(index int) {
		visited[index] = true
		for _, e := range wg.edgesForIndex(index) {
			if !visited[e.v] {
				heap.Push(&pq, e)
			}
		}
	}

	visit(start)

	for !pq.empty() {
		e := heap.Pop(&pq).(edge)
		if visited[e.v] {
			continue
		}
		result = append(result, e)
		visit(e.v)
	}
	return result
}

func printWeightedPath(wg graph, wp []edge) {
	for _, e := range wp {
		fmt.Printf("%v %0.2f > %v\n", wg.vertexAt(e.u), e.weight, wg.vertexAt(e.v))
	}
	fmt.Println("Total Weight: ", totalWeight(wp))
}

func main() {
	citiGraph := graph{}
	citiGraph.init([]vertex{{"Seattle"}, {"San Francisco"}, {"Los Angeles"}, {"Riverside"}, {"Phoenix"}, {"Chicago"},
		{"Boston"}, {"New York"}, {"Atlanta"}, {"Miami"}, {"Dallas"}, {"Houston"}, {"Detroit"},
		{"Philadelphia"}, {"Washington"}})
	citiGraph.addEdgeByVertices(vertex{"Seattle"}, vertex{"Chicago"}, 1737)
	citiGraph.addEdgeByVertices(vertex{"Seattle"}, vertex{"San Francisco"}, 678)
	citiGraph.addEdgeByVertices(vertex{"San Francisco"}, vertex{"Riverside"}, 386)
	citiGraph.addEdgeByVertices(vertex{"San Francisco"}, vertex{"Los Angeles"}, 348)
	citiGraph.addEdgeByVertices(vertex{"Los Angeles"}, vertex{"Riverside"}, 50)
	citiGraph.addEdgeByVertices(vertex{"Los Angeles"}, vertex{"Phoenix"}, 357)
	citiGraph.addEdgeByVertices(vertex{"Riverside"}, vertex{"Phoenix"}, 307)
	citiGraph.addEdgeByVertices(vertex{"Riverside"}, vertex{"Chicago"}, 1704)
	citiGraph.addEdgeByVertices(vertex{"Phoenix"}, vertex{"Dallas"}, 887)
	citiGraph.addEdgeByVertices(vertex{"Phoenix"}, vertex{"Houston"}, 1015)
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Chicago"}, 805)
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Atlanta"}, 721)
	citiGraph.addEdgeByVertices(vertex{"Dallas"}, vertex{"Houston"}, 225)
	citiGraph.addEdgeByVertices(vertex{"Houston"}, vertex{"Atlanta"}, 702)
	citiGraph.addEdgeByVertices(vertex{"Houston"}, vertex{"Miami"}, 968)
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Chicago"}, 588)
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Washington"}, 543)
	citiGraph.addEdgeByVertices(vertex{"Atlanta"}, vertex{"Miami"}, 604)
	citiGraph.addEdgeByVertices(vertex{"Miami"}, vertex{"Washington"}, 923)
	citiGraph.addEdgeByVertices(vertex{"Chicago"}, vertex{"Detroit"}, 238)
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"Boston"}, 613)
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"Washington"}, 396)
	citiGraph.addEdgeByVertices(vertex{"Detroit"}, vertex{"New York"}, 482)
	citiGraph.addEdgeByVertices(vertex{"Boston"}, vertex{"New York"}, 190)
	citiGraph.addEdgeByVertices(vertex{"New York"}, vertex{"Philadelphia"}, 81)
	citiGraph.addEdgeByVertices(vertex{"Philadelphia"}, vertex{"Washington"}, 123)

	result := mst(citiGraph, 0)
	printWeightedPath(citiGraph, result)
}
