
The `calculateWeights` function takes in two parameters: `subgraph` represented as a graph with string nodes and string edges, and `source` representing the source node. It calculates and returns a map of vertex properties for the nodes in the subgraph with respect to the given source node.

## Parameters

- `subgraph`: A graph with string nodes and string edges.
- `source`: A string representing the source node for weight calculation.

## Flow of Execution

1. The function starts by obtaining the adjacency map of the `subgraph` using `subgraph.AdjacencyMap()`. This map contains information about the connections between nodes.
2. The function initializes an empty map called `weight` to store the calculated vertex properties. This map will have node names as keys and `VertexPropertiesDTO` as values.
3. For each node in the `SubGraphAdjacencyMap`, it initializes the `weight` map with default values for `InFlow` and `OutFlow` set to 0.
4. It then sets the `InFlow` and `OutFlow` properties for the `source` node in the `weight` map. `InFlow` is set to 1, and `OutFlow` is calculated as 1 divided by the number of outgoing edges from the `source` node.
5. The function initializes a list called `nodeList` to store nodes for which weights need to be calculated.
6. For each node adjacent to the `source` node (found in the `SubGraphAdjacencyMap`), it adds them to the `nodeList`.
7. The function enters a loop that continues until `nodeList` is empty.
   - Inside the loop, it initializes an empty list called `nodeListSon` to store the next set of nodes for weight calculation.
   - For each node in the `nodeList`, it updates the weights of nodes in the `weight` map using the `updateWeights` function.
   - It iterates through the nodes in `nodeList` to find their adjacent nodes and appends them to `nodeListSon` if they are not already present.
   - The `nodeList` is then updated with the new set of nodes in `nodeListSon`.
8. The function returns the `weight` map, which contains the calculated weights for each node in the subgraph with respect to the given source node.

This function is responsible for determining the weights of nodes in a subgraph based on their connectivity to the source node.


```go
func calculateWeights(subgraph graph.Graph[string, string], source string) map[string]dtos.VertexPropertiesDTO {

	SubGraphAdjacencyMap, err := subgraph.AdjacencyMap()
	utils.Must(err)

	weight := make(map[string]dtos.VertexPropertiesDTO)

	for node := range SubGraphAdjacencyMap {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: 0, OutFlow: 0}
	}

	weight[source] = dtos.VertexPropertiesDTO{InFlow: 1, OutFlow: 1 / float64(len(SubGraphAdjacencyMap[source]))}

	var nodeList []string

	for node := range SubGraphAdjacencyMap[source] {
		nodeList = append(nodeList, node)
	}

	for len(nodeList) != 0 {
		var nodeListSon []string
		for _, node := range nodeList {
			weight = updateWeights(weight, SubGraphAdjacencyMap, node)
		}
		for _, nodes := range nodeList {
			for node := range SubGraphAdjacencyMap[nodes] {
				if !ContainsNode(nodeListSon, node) {
					nodeListSon = append(nodeListSon, node)
				}
			}
		}

		nodeList = nodeListSon
	}

	return weight
}