
The `ContainsAllNode` function checks if a specified `dtos.EdgeDTO` object is contained within a given slice of `dtos.EdgeDTO` objects. It returns `true` if the object is found, indicating that the specified object exists in the slice.

## Parameters

- `slice`: A slice of `dtos.EdgeDTO` objects.
- `nodeToFind`: A `dtos.EdgeDTO` object representing the node to be found in the slice.

## Flow of Execution

1. The function iterates through each `dtos.EdgeDTO` object in the `slice`.
2. For each object in the `slice`, it checks if both the `FromNode` and `ToNode` properties match the corresponding properties in the `nodeToFind` object.
3. If a match is found (i.e., both `FromNode` and `ToNode` properties match), the function returns `true`, indicating that the specified object is present in the `slice.
4. If the function reaches the end of the loop without finding a matching object, it returns `false`, indicating that the specified object is not contained in the `slice.

This function is used to determine whether a specific `dtos.EdgeDTO` object exists within a given slice of such objects.

```go
func ContainsAllNode(slice []dtos.EdgeDTO, nodeToFind dtos.EdgeDTO) bool {
	for _, node := range slice {
		if node.FromNode == nodeToFind.FromNode && node.ToNode == nodeToFind.ToNode {
			return true
		}
	}
	return false
}