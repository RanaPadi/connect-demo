## Function `RemoveSelfLoop`

### Description
The function `RemoveSelfLoop` removes self-loop edges from a directed acyclic graph represented by a graph object (`dspgGraph`). A self-loop occurs when an edge connects a vertex to itself. The function iterates through each vertex in the graph, identified by its adjacency map, and removes any self-loop edges associated with that vertex. The modified graph, with self-loops removed, is then returned.

### Parameters
- `dspgGraph` (type: `graph.Graph[string, string]`, required): The directed acyclic graph from which self-loop edges should be removed.

### Return Type
- A modified graph (`graph.Graph[string, string]`) with self-loop edges removed.

### Execution Flow
1. Obtain the adjacency map of the directed acyclic graph (`dspgGraph`) using the `AdjacencyMap` function.
2. Iterate through each node in the adjacency map.
3. For each node, remove any self-loop edges by calling the `RemoveEdge` function with the node as both the source and target.
4. Return the modified graph with self-loops removed.

### Snippet
```go
func RemoveSelfLoop(dspgGraph graph.Graph[string, string]) graph.Graph[string, string] {
	adj, _ := dspgGraph.AdjacencyMap()

	for node := range adj {
		dspgGraph.RemoveEdge(node, node)
	}

	return dspgGraph
}
