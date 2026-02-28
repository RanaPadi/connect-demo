
The `ContainsAllNode` function checks whether a specified `dtos.EdgeDTO` object is present within a given slice of `dtos.EdgeDTO` objects. It returns `true` if the specified object is found in the slice, indicating its presence.

## Parameters

- `slice`: A slice of `dtos.EdgeDTO` objects.
- `nodeToFind`: A `dtos.EdgeDTO` object representing the node to be found within the slice.

## Flow of Execution

1. The function iterates through each `dtos.EdgeDTO` object in the `slice.
2. For each object in the `slice, it compares the `FromNode` and `ToNode` properties of the object with the corresponding properties in the `nodeToFind` object.
3. If there is a match (i.e., both `FromNode` and `ToNode` properties match), the function returns `true`, indicating that the specified object is present in the `slice.
4. If the loop completes without finding a matching object, the function returns `false`, signifying that the specified object is not within the `slice.

This function is useful for determining the existence of a particular `dtos.EdgeDTO` object within a given slice of such objects.


```go
func ContainsNode[T string | dtos.EdgeDTO](slice []T, nodeToFind T) bool {
	for _, node := range slice {
		if node == nodeToFind {
			return true
		}
	}
	return false
}
