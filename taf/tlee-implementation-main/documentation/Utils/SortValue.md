## Function `SortValue`

### Description
The function `SortValue` takes a map with string keys and boolean values (`map[string]bool`) and returns a new map with the same key-value pairs sorted based on the keys in ascending order. It achieves this by converting the input map into a slice of structs, sorting the slice by the keys, and then constructing a new map from the sorted slice.

### Parameters
- `myMap` (type: `map[string]bool`, required): The input map to be sorted.

### Return Type
- A new map with the same key-value pairs as the input, but sorted based on the keys.

### Execution Flow
1. Create a slice of structs (`data`) to store key-value pairs from the input map.
2. Iterate through the input map (`myMap`) and append each key-value pair as a struct to the `data` slice.
3. Sort the `data` slice based on the keys in ascending order using the `sort.Slice` function and a custom sorting function.
4. Create a new map (`sortedMap`) and populate it with key-value pairs from the sorted `data` slice.
5. Return the `sortedMap`.

### Snippet 
```go
func SortValue(myMap map[string]bool) map[string]bool {
	var data []struct {
		Key   string
		Value bool
	}

	for k, v := range myMap {
		data = append(data, struct {
			Key   string
			Value bool
		}{k, v})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Key < data[j].Key
	})

	sortedMap := make(map[string]bool)
	for _, item := range data {
		sortedMap[item.Key] = item.Value
	}

	return sortedMap
}