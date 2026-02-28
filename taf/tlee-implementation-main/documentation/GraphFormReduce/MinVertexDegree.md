## Function `MinVertexDegree`

### Description
The function `MinVertexDegree` calculates the minimum degree from a list of vertex degrees. It iterates through the provided list and finds the smallest value, representing the minimum vertex degree. The function then returns this minimum degree.

### Parameters
- `vertexList` (type: `[]int`, required): A list of integers representing vertex degrees.

### Return Type
- An integer representing the minimum vertex degree in the provided list.

### Execution Flow
1. Initialize the variable `min` with the first element of the `vertexList`.
2. Iterate through the `vertexList`, comparing each value with the current minimum (`min`).
3. If a value in the list is smaller than the current minimum, update the minimum value to the new value.
4. Return the final minimum vertex degree.

### Snippet
```go
func MinVertexDegree(vertexList []int) int {
	min := vertexList[0]
	for _, value := range vertexList {
		if value < min {
			min = value
		}
	}

	return min
}
