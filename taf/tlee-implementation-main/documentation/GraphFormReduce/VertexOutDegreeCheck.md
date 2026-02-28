## Function `VertexOutDegreeCheck`

### Description
The function `VertexOutDegreeCheck` checks whether a specified vertex in a directed graph (`g`) has outgoing edges. It does so by examining the adjacency map of the graph and determining if the given vertex has any outgoing edges.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph in which the out-degree of the vertex is checked.
- `vertex` (type: `string`, required): The target vertex for which the presence of outgoing edges is checked.

### Return Type
- A boolean value indicating whether the specified vertex has outgoing edges (`true`) or not (`false`).

### Execution Flow
1. Obtain the adjacency map of the directed graph (`g`) using the `AdjacencyMap` function.
2. Check if the length of the outgoing edges associated with the specified `vertex` in the adjacency map is greater than 0.
3. Return `true` if the vertex has outgoing edges, and `false` otherwise.

### Example Usage
```go
func VertexOutDegreeCheck(g graph.Graph[string, string], vertex string) bool {
	adj, err := g.AdjacencyMap()
	utils.Must(err)

	return len(adj[vertex]) > 0
}
