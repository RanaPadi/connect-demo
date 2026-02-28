
The `findAllEdgesToTargetIterative` function iteratively explores a graph to find all edges that lead to a specified `target` node. It returns a list of `dtos.EdgeDTO` objects representing those edges.

## Parameters

- `dspgGraph`: A map representing a graph with nodes and edges, where the keys are node names, and the values are maps of neighboring nodes and their corresponding edges.
- `target`: A string representing the target node to which incoming edges are to be found.

## Flow of Execution

1. The function initializes several data structures:
   - `edgesToTarget`: A list to store the `dtos.EdgeDTO` objects found during the search.
   - `stack`: A stack used for iterative traversal, initially containing the `target` node.
   - `visited`: A map to track visited nodes during the search.
   - `visitedEdges`: A map to track visited edges to avoid duplicates.
2. The function enters a loop that continues until the `stack` is empty.
3. Within the loop, it retrieves the last node from the `stack` and marks it as visited by setting `visited[node]` to `true`.
4. If the node has already been visited, the function skips it and continues to the next node.
5. It then iterates through the nodes and edges in the `dspgGraph` to check for edges leading to the current `node`.
6. For each node that has an edge leading to the current `node`, it creates an `dtos.EdgeDTO` object representing the edge from that node to the current `node`.
7. If the edge has not been visited before (checked using `visitedEdges`), it is added to the `edgesToTarget` list, and the edge is marked as visited in `visitedEdges`.
8. The parent node, from which the edge originates, is pushed onto the `stack` for further exploration.
9. The loop continues until all nodes and edges leading to the `target` node have been explored.
10. The function returns the list of `edgesToTarget`, which contains all edges leading to the `target` node.

```go
func findAllEdgesToTargetIterative(dspgGraph map[string]map[string]graph.Edge[string], target string) []dtos.EdgeDTO {
	edgesToTarget := []dtos.EdgeDTO{}
	stack := []string{target}
	visited := make(map[string]bool)
	visitedEdges := make(map[dtos.EdgeDTO]bool)

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[node] {
			continue
		}

		visited[node] = true

		for parentNode, edges := range dspgGraph {
			for neighbor := range edges {
				if neighbor == node {
					edge := dtos.EdgeDTO{FromNode: parentNode, ToNode: node}

					// Verifica se l'edge è già stato visitato
					if !visitedEdges[edge] {
						edgesToTarget = append(edgesToTarget, edge)
						visitedEdges[edge] = true
					}

					stack = append(stack, parentNode)
				}
			}
		}
	}

	return edgesToTarget
}