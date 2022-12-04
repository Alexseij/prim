package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Graph struct {
	graph [][]*Edge
}

type Edge struct {
	leftVertex  int
	rightVertex int
	weight      int
}

func NewGraph(amountOfVertices int) *Graph {
	return &Graph{
		graph: make([][]*Edge, amountOfVertices),
	}
}

func (g *Graph) AddEdges(edges ...*Edge) {
	for _, edge := range edges {
		g.graph[edge.leftVertex] = append(g.graph[edge.leftVertex], edge)
	}
}

func main() {
	file, err := os.OpenFile("input.txt", os.O_RDONLY, 0755)
	if err != nil {
		panic(err.Error())
	}

	var amountOfVertices int

	scanner := bufio.NewScanner(file)

	if ok := scanner.Scan(); ok {
		line := scanner.Text()
		_, err = fmt.Sscanf(line, "%d", &amountOfVertices)
		if err != nil {
			panic(err.Error())
		}
	} else {
		panic("scanner not working")
	}

	graph := NewGraph(amountOfVertices)

	for scanner.Scan() {
		var (
			leftVertex  int
			rightVertex int
			weight      int
		)

		line := scanner.Text()
		_, err := fmt.Sscanf(line, "%d %d %d", &leftVertex, &rightVertex, &weight)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err.Error())
		}

		graph.AddEdges(
			&Edge{
				leftVertex:  leftVertex,
				rightVertex: rightVertex,
				weight:      weight,
			},
			&Edge{
				leftVertex:  rightVertex,
				rightVertex: leftVertex,
				weight:      weight,
			},
		)
	}
	resultGraph := prism(graph)
	printGraph(resultGraph)
}

func printGraph(graph *Graph) {
	for i := 0; i < len(graph.graph); i++ {
		if len(graph.graph[i]) == 0 {
			continue
		}
		fmt.Printf("For %d :\n", i)
		for _, edge := range graph.graph[i] {
			fmt.Printf("{%d , %d} -> %d;", edge.leftVertex, edge.rightVertex, edge.weight)
		}
		println()
	}
}

func FindMinEdge(
	edges []*Edge,
	vertexSet map[int]bool,
	ch chan *Edge,
) {
	var currentMinEdge *Edge
	for _, edge := range edges {
		if vertexSet[edge.rightVertex] {
			continue
		}
		if currentMinEdge == nil || edge.weight < currentMinEdge.weight {
			currentMinEdge = edge
		}
	}
	ch <- currentMinEdge
}

func prism(g *Graph) *Graph {
	var edges []*Edge

	vertexSet := map[int]bool{}
	amountOfVertex := len(g.graph)
	ch := make(chan *Edge, amountOfVertex)

	resultGraph := NewGraph(amountOfVertex)

	defer close(ch)

	vertexSet[0] = true
	amountOfVertexInSet := 1 
	for amountOfVertexInSet != amountOfVertex {
		for vertex := range vertexSet {
			go FindMinEdge(
				g.graph[vertex],
				vertexSet,
				ch,
			)
		}

		var currentMinEdge *Edge
		amountOfEdges := 0

		for edge := range ch {
			amountOfEdges++

			if edge != nil && (currentMinEdge == nil || edge.weight < currentMinEdge.weight) {
				currentMinEdge = edge
			}

			if amountOfEdges == amountOfVertexInSet {
				vertexSet[currentMinEdge.rightVertex] = true
				amountOfVertexInSet++
				edges = append(edges, currentMinEdge)
				break
			}
		}
	}

	resultGraph.AddEdges(edges...)
	return resultGraph
}
