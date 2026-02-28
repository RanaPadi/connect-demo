## Method `GenerateRandomDSPGGraph`

### Description
This function generates a random DSPG Directed Acyclic Graph (DAG) with a specified number of nodes (`n`) using the gonum/graph library. The generated graph is weighted, acyclic, and may have dynamic edges added or modified based on random processes. The resulting graph is returned.

### Parameters
- `n` (type: `int`, required): The number of nodes in the generated graph.

### Return Type
- `graph.Graph[string, string]`: A directed acyclic graph with string vertices and string edge labels.

### Execution Flow
1. A new directed acyclic graph (`Graph`) is created using the gonum/graph library with string vertices and string edge labels.
2. A map (`mapGraph`) is created to keep track of the edges and their weights in the graph.
3. Source and target vertices (`source` and `target`) are added to the graph.
4. An edge between the source and target vertices is added to the graph with a default weight of 1.0.
5. The mapGraph is updated to reflect the weight of the source-to-target edge.
6. A loop is executed `n` times to add random edges to the graph.
7. Randomly selected edges are obtained from the current graph.
8. A random number (`randomNumber`) is generated to determine which edge to modify.
9. A coin toss (`coin`) is simulated with another random number to decide whether to add a new node or modify an existing edge.
10. If the coin toss results in 0, a new node is added, and edges are modified accordingly.
11. If the coin toss results in 1, the weight of the randomly selected edge is increased by 1.
12. The function ensures acyclicity by removing self-loops from the graph.
13. The final graph is returned.

### Method Snippet
```go
func GenerateRandomDSPGGraph(n int) graph.Graph[string, string] {
	Graph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())
	mapGraph := make(map[string]int)

	source := "s"
	target := "t"

	Graph.AddVertex(source)
	Graph.AddVertex(target)
	Graph.AddEdge(source, target, graph.EdgeWeight(1.0))
	mapGraph[source+target] = 1

	for i := 0; i < n; i++ {
		links, _ := Graph.Edges()
		randomNumber := dspgFunc.GenerateRandomNumber(1, len(links))

		selectedEdge := links[randomNumber-1]
		coin := dspgFunc.GenerateRandomNumber(0, 1)

		if coin == 0 {
			node := strconv.Itoa(i)
			Graph.AddVertex(node)
			Graph.AddEdge(selectedEdge.Source, node)
			Graph.AddEdge(node, selectedEdge.Target)
			mapGraph[selectedEdge.Source+node] = 1
			mapGraph[node+selectedEdge.Target] = 1

			if mapGraph[selectedEdge.Source+selectedEdge.Target] == 1 {
				Graph.RemoveEdge(selectedEdge.Source, selectedEdge.Target)
				delete(mapGraph, selectedEdge.Source+selectedEdge.Target)
			} else {
				mapGraph[selectedEdge.Source+selectedEdge.Target] -= 1
			}

		} else {
			Graph.AddEdge(selectedEdge.Source, selectedEdge.Target, graph.EdgeWeight(1.0))
			mapGraph[selectedEdge.Source+selectedEdge.Target] += 1
		}
	}

	Graph = dspgFunc.RemoveSelfLoop(Graph)

	return Graph
}
