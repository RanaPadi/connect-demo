The `ifConnected` function checks if two nodes are connected in a graph. It returns `true` if there is a connection between the `start` and `neighbor` nodes, otherwise, it returns `false`.

## Parameters

- `graph`: A map representing a graph with nodes and edges, where the keys are node names, and the values are maps of neighboring nodes and their corresponding edges.
- `start`: A string representing the starting node.
- `neighbor`: A string representing the neighboring node to check for connectivity.

## Flow of Execution

1. The function iterates through the neighbors of the `start` node by looping through the nodes within the map `graph[start]`.
2. For each neighbor, it checks if it is equal to the specified `neighbor` node.
3. If a neighbor is found that matches the `neighbor` node, the function returns `true`, indicating that there is a connection between the `start` and `neighbor` nodes.
4. If the loop completes without finding a matching neighbor, the function returns `false`, indicating that there is no direct connection between the two nodes.

This function is used to determine whether there is a direct connection between two nodes in a given graph.

```go
func ifConnected(graph map[string]map[string]graph.Edge[string], start string, neighbor string) bool {
	for tempNeighbor := range graph[start] {
		if tempNeighbor == neighbor {
			return true
		}
	}

	return false
}
