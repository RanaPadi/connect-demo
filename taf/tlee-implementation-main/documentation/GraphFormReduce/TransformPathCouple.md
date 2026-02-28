The `TransformPathToCouple` function converts a list of node names into a list of `dtos.EdgeDTO` objects, representing edges connecting those nodes.

## Parameters

- `slice`: A list of strings representing node names that form a path.

## Flow of Execution

1. The function initializes an empty list called `allEdge` to store `dtos.EdgeDTO` objects representing edges.
2. It calculates the length of the input `slice` and assigns it to `sliceLength`.
3. The function checks if the `slice` contains more than one element (i.e., `sliceLength` is greater than 1).
4. If the condition is met, the function enters a loop that iterates through the elements in the `slice` except the last element (i.e., `sliceLength - 1` iterations).
5. For each iteration, it creates a new `dtos.EdgeDTO` object to represent an edge.
6. It assigns the current element in the `slice` as the `FromNode` of the edge and the next element as the `ToNode`.
7. The edge is added to the `allEdge` list.
8. The loop continues until all elements in the `slice` have been processed.
9. The function returns the `allEdge` list, which contains `dtos.EdgeDTO` objects representing the edges between nodes.


```go
func TransformPathToCouple(slice []string) []dtos.EdgeDTO {

	var allEdge []dtos.EdgeDTO
	sliceLenght := len(slice)

	if sliceLenght > 1 {

		for i := 0; i < sliceLenght-1; i++ {
			edge := dtos.EdgeDTO{}

			edge.FromNode = slice[i]
			edge.ToNode = slice[i+1]
			allEdge = append(allEdge, edge)
		}
	}

	return allEdge
}