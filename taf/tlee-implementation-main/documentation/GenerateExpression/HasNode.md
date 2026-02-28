The public function takes in two parameters: `graph` of type `AllSubPPS` and `node` of type string. It checks whether a specified `node` exists in the `graph.Graph` represented by the `graph` parameter.

### Parameters

- `graph`: An object of type `AllSubPPS` containing a graph.
- `node`: A string representing the node to be checked for existence.

### Return Values

- `bool`: A boolean value indicating whether the specified `node` exists in the `graph.Graph`.

### Flow of Execution

1. The function starts by obtaining the adjacency map of `graph.Graph` and stores it in the `adj` variable.
2. It iterates through the keys in the `adj` map, which represent nodes in the graph.
3. For each node in the map, it checks if it matches the specified `node`.
4. If a match is found, the function returns `true` to indicate that the `node` exists in the graph.
5. If no match is found after iterating through all nodes, the function returns `false` to indicate that the `node` does not exist in the graph.

**Snippet**

```go

func HasNode(graph dtos.AllSubPPSDTO, node string) bool {
 adj, err := graph.Graph.AdjacencyMap()
 utils.Must(err)
 for nodeList := range adj {
  if nodeList == node {
   return true
  }
 }
 return false
}
