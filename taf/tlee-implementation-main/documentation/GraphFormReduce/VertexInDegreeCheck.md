## Function `VertexInDegreeList`

### Description
The function `VertexInDegreeList` retrieves a list of vertices in a directed graph (`g`) that have an incoming edge from a specified vertex. In other words, it returns the vertices that are the source of incoming edges to the given vertex. The function iterates through the adjacency map of the graph and identifies vertices connected to the specified vertex by incoming edges.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph in which the in-degree list is calculated.
- `vertex` (type: `string`, required): The target vertex for which incoming vertices are identified.

### Return Type
- A slice of strings representing the vertices that have incoming edges to the specified vertex.

### Execution Flow
1. Obtain the adjacency map of the directed graph (`g`) using the `AdjacencyMap` function.
2. Initialize an empty slice (`vertexList`) to store vertices with incoming edges.
3. Iterate through each node in the adjacency map.
   - For each node, iterate through its outgoing edges.
   - If the target of an edge matches the specified `vertex`, add the source node to the `vertexList`.
4. Return the final slice containing vertices with incoming edges to the specified vertex.

### Snippet
```go
func VertexInDegreeList(g graph.Graph[string, string], vertex string) []string {

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var vertexList []string

	for node, edges := range adj {
		for edge := range edges {
			if edge == vertex {
				vertexList = append(vertexList, node)
			}
		}
	}

	return vertexList
}
