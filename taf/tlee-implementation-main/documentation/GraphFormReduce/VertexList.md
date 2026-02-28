## Function `VertexList`

### Description
The function `VertexList` retrieves a list of all vertices in a directed graph (`g`). It iterates through the adjacency map of the graph and compiles a list of unique vertices present in the graph.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph for which the vertex list is obtained.

### Return Type
- A slice of strings representing all vertices in the directed graph.

### Execution Flow
1. Obtain the adjacency map of the directed graph (`g`) using the `AdjacencyMap` function.
2. Initialize an empty slice (`vertexList`) to store the vertices.
3. Iterate through each vertex in the adjacency map and add it to the `vertexList`.
4. Return the final slice containing all vertices in the directed graph.

### Snippet
```go
func VertexList(g graph.Graph[string, string]) []string {
	var vertexList []string

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	for vertex := range adj {
		vertexList = append(vertexList, vertex)
	}

	return vertexList
}
