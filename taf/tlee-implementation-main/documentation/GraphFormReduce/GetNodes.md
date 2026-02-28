The `GetNodes` function extracts and returns the unique source nodes from a list of `dtos.EdgeDTO` objects. It collects and returns a list of these source nodes.

## Parameters

- `processingEdges`: A slice of `dtos.EdgeDTO` objects from which source nodes are to be extracted.

## Flow of Execution

1. The function calculates the length of the `processingEdges` slice and assigns it to `n`.
2. It initializes an empty string slice called `nodes` to store the unique source nodes found.
3. The function enters a loop that iterates through the `processingEdges`, excluding the last element (i.e., `n-1` iterations).
4. For each iteration, it appends the `FromNode` property of the `dtos.EdgeDTO` object to the `nodes` slice.
5. After processing all elements except the last one, the function returns the `nodes` slice, which contains the unique source nodes.

This function is used to extract and collect the source nodes from a list of `dtos.EdgeDTO` objects, which can be useful for various graph-related operations.

```go
func GetNodes(processingEdges []dtos.EdgeDTO) []string {
	n := len(processingEdges)
	nodes := make([]string, 0, n-1)

	for i := 0; i < n-1; i++ {
		nodes = append(nodes, processingEdges[i].FromNode)
	}

	return nodes
}