## Function `VertexOutDegreeList`

### Description
The function `VertexOutDegreeList` retrieves a list of vertices in a directed graph (`g`) that are the targets of outgoing edges from a specified vertex. It iterates through the adjacency map of the graph and compiles a list of vertices connected to the specified vertex by outgoing edges.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph in which the out-degree list is calculated.
- `vertex` (type: `string`, required): The source vertex for which outgoing vertices are identified.

### Return Type
- A slice of strings representing the vertices that are the targets of outgoing edges from the specified vertex.

### Execution Flow
1. Obtain the adjacency map of the directed graph (`g`) using the `AdjacencyMap` function.
2. Initialize an empty slice (`vertexList`) to store vertices with outgoing edges.
3. Iterate through the vertices in the adjacency map associated with the specified `vertex` and add them to the `vertexList`.
4. Return the final slice containing vertices that are the targets of outgoing edges from the specified vertex.

### Snippet
```go
func VertexOutDegreeList(g graph.Graph[string, string], vertex string) []string {

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var vertexList []string

	for vertex := range adj[vertex] {
		vertexList = append(vertexList, vertex)
	}

	return vertexList
}
