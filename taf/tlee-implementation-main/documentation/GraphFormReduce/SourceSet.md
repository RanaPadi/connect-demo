The `SourceSet` function populates a `sourceSet` map with information about edges originating from source nodes in a graph. It utilizes parallel processing with goroutines and `sync.WaitGroup` to enhance efficiency.

## Parameters

- `wg`: A pointer to a `sync.WaitGroup` to manage concurrent execution of goroutines.
- `adjacencyMap`: A map representing an adjacency map of a graph, where keys are node names, and values are maps of neighboring nodes and their corresponding edges.
- `sourceSet`: A map used to store information about edges originating from source nodes.

## Flow of Execution

1. The function starts by marking a goroutine as done using `defer wg.Done()` to indicate the completion of the goroutine.
2. It obtains a list of source nodes within the graph by calling the `getSource` function and stores them in the `sources` variable.
3. The function then iterates through each source node in the `sources` list.
4. For each source node, it calls the `findAllEdgesFromSourceIterative` function to find all edges originating from that source node.
5. The resulting edges are appended to the `sourceSet` map under the corresponding source node key.
6. The function repeats this process for each source node.
7. As goroutines may run concurrently to process different source nodes, this parallel processing speeds up the task of gathering information about edges from multiple sources.

This function is used to populate a `sourceSet` map with information about edges originating from source nodes within a graph, optimizing the process through concurrent goroutines managed by a `sync.WaitGroup`.

```go
func SourceSet(wg *sync.WaitGroup, adjacencyMap map[string]map[string]graph.Edge[string], sourceSet map[string][]dtos.EdgeDTO) {
	defer wg.Done()

	sources := getSource(adjacencyMap)
	for _, node := range sources {
		edgeFromSource := findAllEdgesFromSourceIterative(adjacencyMap, node)
		sourceSet[node] = append(sourceSet[node], edgeFromSource...)
	}
}