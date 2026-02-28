The `nodesBetween` function determines and returns a list of nodes that lie on the path between a `startNode` and an `endNode` within a graph, using depth-first search (DFS).

## Parameters

- `adjMatrix`: A map representing an adjacency matrix of a graph, where keys are node names, and values are maps of neighboring nodes and their corresponding edges.
- `startNode`: A string representing the starting node.
- `endNode`: A string representing the ending node.

## Flow of Execution

1. The function defines a nested function `dfs`, which is a depth-first search (DFS) traversal function. It recursively explores nodes in the graph.
2. The outer function initializes an empty list `nodesBetweenList` to store nodes between the `startNode` and `endNode`.
3. It creates a `visited` map to track visited nodes during the DFS traversal.
4. The DFS traversal starts from the `startNode`. The `dfs` function is called with the `startNode` and `visited` as arguments.
5. Within the `dfs` function, the current node is marked as visited.
6. It iterates through the neighbors of the current node in the `adjMatrix`.
7. If an unvisited neighbor is encountered, it is added to the `nodesBetweenList` to indicate that it lies between the `startNode` and `endNode`.
8. The DFS continues recursively for unvisited neighbors until all paths are explored.
9. After the DFS traversal, the function checks if the `endNode` was visited. If it was, the `nodesBetweenList` is returned as the list of nodes between the two nodes.
10. If the `endNode` was not visited, an empty list is returned, indicating that there is no path between the `startNode` and `endNode`.

This function is used to identify and return a list of nodes that exist on the path between a specified `startNode` and `endNode` within a given graph, utilizing a depth-first search (DFS) algorithm.

```go
func MinNestingLevel(dspgGraph graph.Graph[string, string]) int {
	links, err := dspgGraph.Edges()
	utils.Must(err)
	minNestingLevel := links[0].Properties.Weight
	for _, link := range links {
		if link.Properties.Weight < minNestingLevel {
			minNestingLevel = link.Properties.Weight
		}
	}
	return minNestingLevel
}