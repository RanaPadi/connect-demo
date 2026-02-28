The public function takes in a parameter `nodeList` represented as a map of strings to `VertexProperties`. It calculates and returns the count of nodes in the `nodeList` with an `InFlow` value of 1.0.

### Parameters

- `nodeList`: A map of strings to `VertexProperties` representing nodes and their properties.

### Return Value

- `count`: An integer representing the count of nodes with an `InFlow` value of 0.999 in the `nodeList`.

### Flow of Execution

1. The function initializes a variable `count` to 0 to keep track of the count.
2. It iterates through each `node` in the `nodeList`:
   - It checks if the `InFlow` value of the node is greater than 0.999.
   - If the condition is met, it increments the `count`.
3. The function returns the `count`, which represents the count of nodes in the `nodeList` with an `InFlow` value of 0.999.

**Snippet**

```go
func countInFlow(nodeList map[string]dtos.VertexPropertiesDTO) int {
 var count int = 0
 for _, links := range nodeList {
  if links.InFlow > 0.9999 {
   count++
  }
 }
 return count
}
