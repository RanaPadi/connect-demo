## Function `DiffVertexDegree`

### Description
The function `DiffVertexDegree` calculates the difference between the in-degree and out-degree of vertices in a subgraph (`subgraph`) compared to their corresponding degrees in a larger graph (`g`). It selectively considers vertices in the subgraph based on certain degree conditions and computes the degree differences.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The larger graph containing the overall vertex connections.
- `subgraph` (type: `graph.Graph[string, string]`, required): The subgraph for which the vertex degree differences are calculated.

### Return Type
- An integer representing the sum of the differences between in-degree and out-degree for selected vertices in the subgraph.

### Execution Flow
1. Identify vertices in the subgraph (`subgraph`) that satisfy both in-degree and out-degree conditions using helper functions.
2. For each selected vertex, concurrently retrieve its in-degree and out-degree lists from the larger graph (`g`) using goroutines.
3. Calculate the difference between the lengths of the in-degree and out-degree lists for each vertex and store the differences in a slice (`diff`).
4. Sum the differences in the `diff` slice to obtain the final result.
5. Return the result, which represents the sum of the differences between in-degree and out-degree for selected vertices in the subgraph.


### Snippet
```go
func DiffVertexDegree(g graph.Graph[string, string], subgraph graph.Graph[string, string]) int {
	var v []string

	for _, vertex := range VertexList(subgraph) {
		if VertexInDegreeCheck(subgraph, vertex) && VertexOutDegreeCheck(subgraph, vertex) {
			v = append(v, vertex)
		}
	}

	var diff []int

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

		diff = append(diff, len(vi)-len(vo))
	}

	var result int
	for _, value := range diff {
		result += value
	}

	return result
}
