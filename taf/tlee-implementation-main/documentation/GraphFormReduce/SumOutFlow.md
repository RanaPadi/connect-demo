The public function takes in three parameters: `adj`, `listY`, and `nodeList`. It calculates the sum of outflows for a list of nodes `listY` within a graph represented by the adjacency map `adj`, using the information in the `nodeList`.

### Parameters

- `adj`: A map of strings to maps of strings to `Edge[string]` representing the adjacency information of a graph.
- `listY`: A slice of strings representing a list of nodes for which the outflows need to be summed.
- `nodeList`: A map of strings to `VertexProperties` representing properties associated with nodes in the graph.

### Return Value

- `inFlowNumber`: A float64 representing the sum of outflows for the nodes in the `listY`.

### Flow of Execution

1. The function initializes a float64 variable `inFlowNumber` to zero, which will store the sum of outflows.
2. It iterates through each element `el` in the `listY`:
   - For each `el`, it accesses the corresponding node's properties in the `nodeList`.
   - It adds the `OutFlow` property of the node to the `inFlowNumber`.
3. After iterating through all elements in `listY`, the function returns the `inFlowNumber`, which represents the sum of outflows for the specified nodes.

**Snippet**

```go
func sumOutFlow(listY map[string]struct{}, nodeList map[string]dtos.VertexPropertiesDTO) float64 {
 outFlowSum := 0.0

 for el := range listY {
  if nodeProps, exists := nodeList[el]; exists {
   outFlowSum += nodeProps.OutFlow
  }
 }

 return outFlowSum
}
