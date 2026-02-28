## Function `VertexDegree`

### Description
The function `VertexDegree` calculates the degree of each vertex in a directed graph (`g`). The degree of a vertex is the sum of its in-degree and out-degree. The function concurrently calculates the degree for each vertex in the graph, and the result is a slice containing the degrees of all vertices.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph for which vertex degrees are calculated.

### Return Type
- A slice of integers representing the degrees of each vertex in the graph.

### Execution Flow
1. Initialize an empty slice (`vertexDegree`) to store the degrees of each vertex.
2. Obtain the list of vertices in the graph using the `VertexList` function.
3. Use goroutines to concurrently calculate in-degrees (`vi`) and out-degrees (`vo`) for each vertex.
4. For each vertex, wait for the goroutines to finish and calculate the total degree by summing the in-degrees and out-degrees.
5. Append the total degree to the `vertexDegree` slice for each vertex.
6. Return the final slice containing the degrees of all vertices.

### Snippet
```go
func VertexDegree(g graph.Graph[string, string]) []int {
	var vertexDegree []int

	v := VertexList(g)

	var wg sync.WaitGroup

	var vi []string
	var vo []string

	for _, vertex := range v {
		wg.Add(2)
		go func() {
			defer wg.Done()
			vi = VertexInDegreeList(g, vertex)
		}()

		go func() {
			defer wg.Done()
			vo = VertexOutDegreeList(g, vertex)
		}()
		wg.Wait()

		vertexDegree = append(vertexDegree, len(vi)+len(vo))
	}

	return vertexDegree
}
