
The `containsTRY` function checks whether a specified `node` is found within a given `dtos.SynthesizingGraph`. It returns `true` if the specified node is present, indicating its existence.

## Parameters

- `synthesizingGraph`: A `dtos.SynthesizingGraph` representing a graph structure.
- `key`: A string key used to identify a node within the `synthesizingGraph`.
- `node`: A string representing the node to be searched for within the graph.

## Flow of Execution

1. The function starts by checking if the given `key` exists within the `synthesizingGraph.NodeToPPS` map. If the `key` does not exist, the function immediately returns `true`, indicating that the node is not found.
2. If the `key` exists in the map, the function proceeds to check if the specified `node` matches either the `FromNode` or `ToNode` properties of the corresponding `pps` object associated with the `key`. If a match is found, the function returns `true`, indicating the presence of the node.
3. The function calculates a list of nodes between `pps.FromNode` and `pps.ToNode` using the `allPaths` function.
4. It then iterates through the nodes in the list (`tmp`) and checks if the specified `node` matches any of the nodes in the list. If a match is found, the function returns `true`, indicating the presence of the node.
5. If none of the previous conditions are met, the function returns `false`, indicating that the specified `node` is not found within the graph.

This function is used to determine whether a specific `node` exists within a given `dtos.SynthesizingGraph`, considering its association with a particular `key` and its relationships within the graph structure.

```go
func containsTRY(synthesizingGraph dtos.SynthesizingGraph, key string, node string) bool {
	if _, exists := synthesizingGraph.NodeToPPS[key]; !exists {
		return true
	}

	pps := synthesizingGraph.NodeToPPS[key]
	if node == pps.FromNode || node == pps.ToNode {
		return true
	}

	tmp := allPaths(synthesizingGraph.Adj, pps.FromNode, pps.ToNode)
	for _, tmpNode := range tmp {
		if ContainsNode(tmpNode, node) {
			return true
		}
	}

	return false

}