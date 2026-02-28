The public function accepts two parameters: `slice1` and `slice2`, both represented as arrays of objects with type `EdgeDTO` containing `FromNode` and `ToNode`.

### Parameters:

- `slice1` and `slice2`: An object containing FromNode and ToNode

### Return value:

- `[]dtos.EdgeDTO`: An array containing the intersection of the two arrays.

**Snippet**

```go
func IntersectionSlices(slice1, slice2 []dtos.EdgeDTO) []dtos.EdgeDTO {
	// We create a map to keep track of unique objects in slice1
	set := make(map[dtos.EdgeDTO]struct{})
	for _, item := range slice1 {
		set[item] = struct{}{}
	}

	// We filter slice2 for objects that are also present in the map set
	var intersection []dtos.EdgeDTO
	for _, item := range slice2 {
		if _, exists := set[item]; exists {
			intersection = append(intersection, item)
		}
	}

	return intersection
}
