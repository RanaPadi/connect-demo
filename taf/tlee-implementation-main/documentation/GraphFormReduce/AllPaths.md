## Function `allPaths`

### Description
The function `allPaths` calculates all possible paths between two nodes (`s` and `d`) in a synthesizing graph represented by the `dtos.SynthesizingGraph`. It uses Depth-First Search (DFS) to explore all paths and store them in a slice of string slices (`[][]string`).

### Parameters
- `synthesizingGraph` (type: `dtos.SynthesizingGraph`, required): The synthesizing graph represented by the `dtos.SynthesizingGraph`.
- `s` (type: `string`, required): The starting node for path exploration.
- `d` (type: `string`, required): The destination node for path exploration.

### Return Type
- A slice of string slices (`[][]string`) containing all possible paths from node `s` to node `d` in the synthesizing graph.

### Execution Flow
1. Create a `visited` map to keep track of visited nodes during DFS.
2. Initialize an empty slice `paths` to store all possible paths.
3. Initialize an empty slice `path` to represent the current path being explored.
4. Call the DFS function (`DFS`) to perform depth-first search and populate the `paths` slice.
5. Return the calculated `paths` slice.


```go
func allPaths(synthesizingGraph dtos.SynthesizingGraph, s string, d string) [][]string {
	visited := make(map[string]bool)
	paths := [][]string{}
	path := []string{}

	DFS(synthesizingGraph, s, d, visited, path, &paths)

	return paths
}
