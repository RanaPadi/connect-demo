## Function `SlicePop`

### Description
The function `SlicePop` is a generic function designed to remove an element at a specific index from a slice and return the modified slice along with the removed element. It utilizes the `append` function to achieve the removal while maintaining the order of the remaining elements in the slice.

### Parameters
- `s` (type: `[]T`, required): A slice of any type (`T`) from which an element needs to be removed.
- `i` (type: `int`, required): The index of the element to be removed from the slice.

### Return Type
- The modified slice after removing the element.
- The removed element of type `T`.

### Execution Flow
1. Store the element at index `i` in the input slice (`s`) in the variable `elem`.
2. Use the `append` function to create a new slice that excludes the element at index `i`. This is achieved by combining the elements before index `i` and those after index `i+1`.
3. Return the modified slice and the removed element.

### Snippet
```go
func SlicePop[T any](s []T, i int) ([]T, T) {
	elem := s[i]
	s = append(s[:i], s[i+1:]...)
	return s, elem
}