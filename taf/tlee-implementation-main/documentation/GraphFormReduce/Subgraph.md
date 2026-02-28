## Function `Subgraph`

### Description
The function `Subgraph` creates a subgraph from a given directed acyclic graph (`g`) based on a list of vertices. The vertices are specified as a 2D slice (`vertices`), where each inner slice represents a set of vertices forming a subgraph. The function creates a new directed acyclic subgraph (`subgraph`) and adds the specified vertices along with their associated edges from the original graph. The resulting subgraph is then returned.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed acyclic graph from which the subgraph is created.
- `vertices` (type: `[][]string`, required): A 2D slice where each inner slice represents a set of vertices forming a subgraph.

### Return Type
- A subgraph (`graph.Graph[string, string]`) created from the original graph (`g`) based on the specified vertices.

### Execution Flow
1. Initialize a new directed acyclic subgraph (`subgraph`) with the same characteristics as the original graph (`g`), including string hashing, directed edges, acyclic nature, and weighted edges.
2. Obtain the adjacency map of the original graph using the `AdjacencyMap` function.
3. Create a map (`NodeMap`) to store unique vertices from the specified sets.
4. Iterate through each set of vertices in the 2D slice (`vertices`).
   - For each vertex in a set, add it to the `NodeMap` and the subgraph.
5. Iterate through the unique vertices in the `NodeMap`.
   - For each vertex, add edges from the original graph to the subgraph, preserving the connectivity of the specified vertices.
6. Return the resulting subgraph.

### Snippet
```go
func Subgraph(g graph.Graph[string, string], vertices [][]string) graph.Graph[string, string] {
	subgraph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())
	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var NodeMap map[string]string = make(map[string]string)

	for _, vertice := range vertices {
		for _, v := range vertice {
			NodeMap[v] = v
			subgraph.AddVertex(v)
		}
	}

	for _, v := range NodeMap {
		for _, edge := range adj[v] {
			subgraph.AddEdge(v, edge.Target)
		}

	}

	return subgraph
