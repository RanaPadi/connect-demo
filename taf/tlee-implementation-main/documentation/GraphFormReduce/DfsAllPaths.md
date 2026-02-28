# `dfsAllPaths` Function

The `dfsAllPaths` (Depth-First Search for All Paths) function performs a depth-first search in a given graph to find all possible paths from the `current` node to the `end` node. It populates a slice of paths that represent all such paths.

## Parameters

- `graph`: A map representing a graph with nodes and edges, where the keys are node names and the values are maps of neighboring nodes and their corresponding edges.
- `current`: A string representing the current node being explored.
- `end`: A string representing the target node to reach.
- `visited`: A map used to track visited nodes during the search.
- `path`: A slice of strings representing the current path being explored.
- `paths`: A pointer to a slice of slices of strings, where the result paths are stored.

## Flow of Execution

1. The function starts by marking the `current` node as visited by setting `visited[current]` to `true`.
2. It appends the `current` node to the `path` being explored.
3. If the `current` node is equal to the `end` node, the function appends a copy of the current `path` to the `paths` slice to store this found path.
4. If the `current` node is not the `end` node, the function iterates through the neighbors of the `current` node in the `graph`.
5. For each unvisited neighbor, it recursively calls the `dfsAllPaths` function to explore that neighbor.
6. The function continues to explore all possible paths from the `current` node, and if it encounters the `end` node, it stores the path in the `paths` slice.
7. After exploring all paths from the `current` node, the `visited` status of the `current` node is reset to `false` to allow it to be revisited during other paths.

```go
func dfsAllPaths(graph map[string]map[string]graph.Edge[string], current string, end string, visited map[string]bool, path []string, paths *[][]string) {
	visited[current] = true
	path = append(path, current)

	if current == end {
		*paths = append(*paths, append([]string{}, path...))
	} else {
		for neighbor := range graph[current] {
			if !visited[neighbor] {
				dfsAllPaths(graph, neighbor, end, visited, path, paths)
			}
		}
	}

	visited[current] = false
}