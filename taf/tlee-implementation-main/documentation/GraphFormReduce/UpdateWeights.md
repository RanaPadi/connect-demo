The `updateWeights` function recalculates and updates the properties of a vertex (node) in a graph, specifically its in-flow and out-flow, based on the properties of its neighboring nodes.

## Parameters

- `weight`: A map containing properties of vertices (nodes) represented as `dtos.VertexPropertiesDTO`.
- `SubGraphAdjacencyMap`: A map representing an adjacency map of a subgraph, where keys are node names, and values are maps of neighboring nodes and their corresponding edges.
- `node`: A string representing the node for which the properties are to be updated.

## Flow of Execution

1. The function starts by calling the `findY` function to find a set of neighboring nodes connected to the given `node`. This set of neighbors is stored in the `findY` variable.
2. It calculates the total in-flow from the neighboring nodes in the `findY` set by calling the `sumOutFlow` function and stores it in the `vertexInWeight` variable.
3. It calculates the out-degree of the current `node` by determining the number of neighbors in the `SubGraphAdjacencyMap`.
4. If the out-degree is equal to 0, the function updates the properties of the `node` in the `weight` map, setting its out-flow to 0 and in-flow to `vertexInWeight`.
5. If the out-degree is not 0, the function updates the properties of the `node` in the `weight` map, setting its out-flow to `vertexInWeight` divided by the out-degree, and in-flow to `vertexInWeight`.
6. The function returns the updated `weight` map.

This function is used to dynamically update the properties of a vertex in a graph, considering the in-flow and out-flow values based on its neighboring nodes. This is often used in flow analysis and network algorithms.

```go
func updateWeights(weight map[string]dtos.VertexPropertiesDTO, SubGraphAdjacencyMap map[string]map[string]graph.Edge[string], node string) map[string]dtos.VertexPropertiesDTO {
	findY := findY(SubGraphAdjacencyMap, node)
	vertexInWeight := sumOutFlow(findY, weight)
	outDegree := float64(len(SubGraphAdjacencyMap[node]))

	if outDegree == 0 {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: vertexInWeight, OutFlow: 0}
	} else {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: vertexInWeight, OutFlow: vertexInWeight / outDegree}
	}

	return weight
}