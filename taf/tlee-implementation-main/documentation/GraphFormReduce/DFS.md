The `dfs` (Depth-First Search) function performs a depth-first search in a given graph to find a path from the `start` node to the `end` node. It returns a boolean indicating whether a path exists and the path itself if found.

## Parameters

- `graph`: A map representing a graph with nodes and edges, where the keys are node names and the values are maps of neighboring nodes and their corresponding edges.
- `start`: A string representing the starting node for the search.
- `end`: A string representing the target node to reach.
- `visited`: A map used to track visited nodes during the search.
- `path`: A slice of strings representing the current path being explored.

## Flow of Execution

1. The function starts by marking the `start` node as visited by setting `visited[start]` to `true`.
2. It appends the `start` node to the `path` being explored.
3. If the `start` node is equal to the `end` node, the function returns `true`, indicating that a path has been found, along with the current `path`.
4. The function then iterates through the neighbors of the `start` node in the `graph`.
5. For each unvisited neighbor that is also connected to the `start` node, it recursively calls the `dfs` function to explore that neighbor.
6. If a path is found in the recursive call (`check` is `true`), the function immediately returns `true` and the path.
7. If no path is found in the current recursive call, the function continues to explore other neighbors.
8. If no path is found in the entire search, the function returns `false` and the current `path`.

This function is used for searching a graph to determine if a path exists between two nodes, and it returns the path if one is found.

```go
func dfs(graph map[string]map[string]graph.Edge[string], start, end string, visited map[string]bool, path []string) (bool, []string) {
	var check bool
	visited[start] = true
	path = append(path, start)
	if start == end {
		return true, path
	}

	for neighbor := range graph[start] {
		if !visited[neighbor] && ifConnected(graph, start, neighbor) {
			check, path = dfs(graph, neighbor, end, visited, path)
			if check {
				return true, path
			}
		}
	}

	return false, path
}