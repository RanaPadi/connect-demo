## Function `CreateGraphFromMatrix`

### Description
The function `CreateGraphFromMatrix` generates a directed acyclic graph (DAG) from an adjacency matrix. It uses the `graph` package to represent the graph and returns the constructed graph along with an error (if any).

### Parameters
- `adjMatrix` (type: `[][]int`, required): The adjacency matrix representing the connections between vertices. It is assumed that `adjMatrix[i][j]` is 1 if there is an edge from vertex `i+1` to vertex `j+1`, and 0 otherwise.

### Return Type
- A directed acyclic graph (`graph.Graph[string, string]`).
- An error indicating any issues during the graph construction. If the graph is successfully constructed, the error is `nil`.

### Execution Flow
1. Create a new directed acyclic graph (`g`) using the `graph` package.
2. Determine the number of vertices in the graph based on the length of the adjacency matrix.
3. Add vertices to the graph using vertex labels from "1" to the number of vertices.
4. Iterate through the adjacency matrix and add directed edges to the graph for each 1 in the matrix.
5. Return the constructed graph and a `nil` error.

### Snippet
```go

func CreateGraphFromMatrix(adjMatrix [][]int) (graph.Graph[string, string], error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())

	numVertices := len(adjMatrix)
	for i := 1; i <= numVertices; i++ {
		_ = g.AddVertex(strconv.Itoa(i))
	}

	for i := 0; i < numVertices; i++ {
		for j := 0; j < numVertices; j++ {
			if adjMatrix[i][j] == 1 {
				_ = g.AddEdge(strconv.Itoa(i+1), strconv.Itoa(j+1))
			}
		}
	}

	return g, nil
}
