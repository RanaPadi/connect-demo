The `TargetSet` function populates a `targetSet` map with information about edges leading to target nodes in a graph. Similar to `SourceSet`, it utilizes parallel processing with goroutines and `sync.WaitGroup` for efficiency.

## Parameters

- `wg`: A pointer to a `sync.WaitGroup` to manage concurrent execution of goroutines.
- `adjacencyMap`: A map representing an adjacency map of a graph, where keys are node names, and values are maps of neighboring nodes and their corresponding edges.
- `targetSet`: A map used to store information about edges leading to target nodes.

## Flow of Execution

1. The function starts by marking a goroutine as done using `defer wg.Done()` to indicate the completion of the goroutine.
2. It obtains a list of target nodes within the graph by calling the `getTarget` function and stores them in the `targets` variable.
3. The function then iterates through each target node in the `targets` list.
4. For each target node, it calls the `findAllEdgesToTargetIterative` function to find all edges leading to that target node.
5. The resulting edges are appended to the `targetSet` map under the corresponding target node key.
6. The function repeats this process for each target node.
7. As goroutines may run concurrently to process different target nodes, this parallel processing speeds up the task of gathering information about edges leading to multiple targets.

This function is used to populate a `targetSet` map with information about edges leading to target nodes within a graph, optimizing the process through concurrent goroutines managed by a `sync.WaitGroup`.

```go
func TargetSet(wg *sync.WaitGroup, adjacencyMap map[string]map[string]graph.Edge[string], targetSet map[string][]dtos.EdgeDTO) {
	defer wg.Done()

	targets := getTarget(adjacencyMap)
	for _, node := range targets {
		edgeFromTarget := findAllEdgesToTargetIterative(adjacencyMap, node)
		targetSet[node] = append(targetSet[node], edgeFromTarget...)
	}
}