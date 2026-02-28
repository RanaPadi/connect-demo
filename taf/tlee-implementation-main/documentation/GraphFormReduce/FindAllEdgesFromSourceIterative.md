
The `findAllEdgesFromSourceIterative` function iteratively explores a graph to find all edges that originate from a specified `source` node. It returns a list of `dtos.EdgeDTO` objects representing those edges.

## Parameters

- `dspgGraph`: A map representing a graph with nodes and edges, where the keys are node names, and the values are maps of neighboring nodes and their corresponding edges.
- `source`: A string representing the source node from which to find all outgoing edges.

## Flow of Execution

1. The function initializes a set of data structures:
   - `visited`: A map to track visited nodes during the search.
   - `edgesFromSource`: A list to store the `dtos.EdgeDTO` objects found during the search.
   - `visitedEdges`: A map to track visited edges to avoid duplicates.
   - `stack`: A stack used for iterative traversal, initially containing the `source` node.
2. The function enters a loop that continues until the `stack` is empty.
3. Within the loop, it retrieves the last node from the `stack` and marks it as visited by setting `visited[node]` to `true`.
4. It then iterates through the neighbors of the current `node` in the `dspgGraph`.
5. For each neighbor, it creates an `dtos.EdgeDTO` object representing the edge from the current `node` to the neighbor.
6. If the edge has not been visited before (checked using `visitedEdges`), it is added to the `edgesFromSource` list, and the edge is marked as visited in `visitedEdges`.
7. If the neighbor node has not been visited, it is added to the `stack` for further exploration.
8. The loop continues until all nodes and edges originating from the `source` node have been explored.
9. The function returns the list of `edgesFromSource`, which contains all edges originating from the `source` node.

```go
func findAllEdgesFromSourceIterative(dspgGraph map[string]map[string]graph.Edge[string], source string) []dtos.EdgeDTO {
	visited := make(map[string]bool)
	edgesFromSource := []dtos.EdgeDTO{}
	visitedEdges := make(map[dtos.EdgeDTO]bool)
	stack := []string{source}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visited[node] = true

		for neighbor := range dspgGraph[node] {
			edge := dtos.EdgeDTO{FromNode: node, ToNode: neighbor}

			if !visitedEdges[edge] {
				edgesFromSource = append(edgesFromSource, edge)
				visitedEdges[edge] = true
			}

			if !visited[neighbor] {
				stack = append(stack, neighbor)
			}
		}
	}

	return edgesFromSource
}
