## Function `Intersect`

### Description
The function `Intersect` calculates the intersection of multiple string slices (sets). It takes an arbitrary number of sets as input and returns a new string slice containing the common elements present in all input sets.

### Parameters
- `sets` (type: `...[]string`, variadic): Zero or more string slices representing sets.

### Return Type
- A string slice containing the common elements present in all input sets.

### Execution Flow
1. Check if the number of input sets is zero. If so, return `nil` since there are no sets to intersect.
2. Initialize an `intersection` map to keep track of elements that exist in the first set.
3. Iterate through the elements of the first set and mark them in the `intersection` map.
4. For each subsequent set, create a new map (`currentSet`) to mark its elements.
5. Compare the elements in the `intersection` map with the elements in the current set. Remove any elements in the `intersection` map that do not exist in the current set.
6. After processing all sets, append the remaining elements in the `intersection` map to the result slice.
7. Return the result slice containing the common elements.

### Snippet
```go
func Intersect(sets ...[]string) []string {
	if len(sets) == 0 {
		return nil
	}

	intersection := make(map[string]bool)
	result := []string{}

	for _, element := range sets[0] {
		intersection[element] = true
	}

	for _, set := range sets[1:] {
		currentSet := make(map[string]bool)

		for _, element := range set {
			currentSet[element] = true
		}

		for element := range intersection {
			if !currentSet[element] {
				delete(intersection, element)
			}
		}
	}

	for element := range intersection {
		result = append(result, element)
	}

	return result
}
