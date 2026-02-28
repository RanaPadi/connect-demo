The `MinNestingLevel` function calculates and returns the minimum nesting level among the edges in a given graph. Nesting level typically represents a depth or hierarchy associated with edges.

## Parameters

- `dspgGraph`: A graph represented as a `graph.Graph[string, string]`.

## Flow of Execution

1. The function starts by obtaining a list of edges from the `dspgGraph` using the `Edges` method and stores it in the `links` variable. 
2. It initializes the `minNestingLevel` variable with the nesting level of the first edge in the `links` list.
3. The function then iterates through each edge in the `links` list.
4. For each edge, it compares the nesting level (represented by the `Weight` property) with the current minimum nesting level stored in the `minNestingLevel` variable.
5. If the nesting level of the current edge is less than the current minimum nesting level, the `minNestingLevel` variable is updated to the nesting level of the current edge.
6. The function continues to iterate through all the edges, ensuring that `minNestingLevel` contains the minimum nesting level among all edges.
7. Finally, it returns the calculated minimum nesting level.

This function is used to find and return the minimum nesting level among the edges within a given graph, which is often used to understand the hierarchy or depth of relationships within the graph.


```go
func MinNestingLevel(dspgGraph graph.Graph[string, string]) int {
	links, err := dspgGraph.Edges()
	utils.Must(err)
	minNestingLevel := links[0].Properties.Weight
	for _, link := range links {
		if link.Properties.Weight < minNestingLevel {
			minNestingLevel = link.Properties.Weight
		}
	}
	return minNestingLevel
}
