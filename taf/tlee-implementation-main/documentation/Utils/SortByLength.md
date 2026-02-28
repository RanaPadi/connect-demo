## Function `SortByLength`

### Description
The function `SortByLength` sorts a slice of string slices (`[][]string`) based on the length of the inner string slices in ascending order. It utilizes the `sort.Slice` function with a custom sorting function that compares the lengths of the inner slices.

### Parameters
- `arrays` (type: `[][]string`, required): A slice of string slices to be sorted based on the length of the inner slices.

### Execution Flow
1. Use the `sort.Slice` function to sort the `arrays` slice based on the lengths of the inner string slices.
2. The sorting function compares the length of the string slice at index `i` with the length of the string slice at index `j`.
3. The result of the comparison determines the order in which the slices should be sorted.

### Snippet
```go
func SortByLength(arrays [][]string) {
	sort.Slice(arrays, func(i, j int) bool {
		return len(arrays[i]) < len(arrays[j])
	})
}