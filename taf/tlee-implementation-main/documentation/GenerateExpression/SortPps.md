The public function takes in a slice of `AllSubPPS` objects represented as `allPps` and performs custom sorting based on specified criteria.

### Parameters

- `allPps`: A slice of objects of type `AllSubPPS` representing subgraphs with associated information.

### Flow of Execution

1. The function defines a custom sorting function `customSort` that compares two `AllSubPPS` objects `i` and `j` based on the following criteria:
   - If the `MinNestingLevel` of `i` is greater than the `MinNestingLevel` of `j`, it returns `true`.
   - If the `MinNestingLevel` of `i` is less than the `MinNestingLevel` of `j`, it returns `false`.
   - If the `MinNestingLevel` of `i` is equal to the `MinNestingLevel` of `j`, it further compares the number of edges in their respective graphs. It returns `true` if the number of edges in the graph of `i` is less than the number of edges in the graph of `j`, and `false` otherwise.
2. The `sort.SliceStable` function is called on the `allPps` slice, and the custom sorting function `customSort` is applied for stable sorting.
3. After sorting, the `allPps` slice is updated with the sorted order based on the specified criteria.

**Snippet**

```go
func SortPps(allPps []dtos.AllSubPPSDTO) {

 customSort := func(i, j int) bool {
  if allPps[i].MinNestingLevel > allPps[j].MinNestingLevel {
   return true
  } else if allPps[i].MinNestingLevel < allPps[j].MinNestingLevel {
   return false
  }

  a, err := allPps[i].Graph.Edges()
  utils.Must(err)
  b, err := allPps[j].Graph.Edges()
  utils.Must(err)

  return len(a) < len(b)
 }

 sort.SliceStable(allPps, customSort)
}
