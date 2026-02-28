The public function takes in two parameters: `dspgGraph` and `links`. It updates the nesting level of specified links within a directed graph represented by `dspgGraph`.

### Parameters

- `dspgGraph`: A directed graph with string nodes and string edges.
- `links`: A slice of `graph.Edge[string]` representing the links whose nesting levels need to be updated.

### Flow of Execution

1. The function iterates through each `link` in the `links` slice:
   - For each `link`, it retrieves the corresponding edge from the `dspgGraph` using `link.Source` and `link.Target`.
   - It extracts the current nesting level from the edge's properties.
   - It increments the nesting level by 1 to update it.
   - It updates the edge's weight with the new nesting level using `dspgGraph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(nestingLevel))`.

**Snippet**

```go
func UpsertNestingLevel(dspgGraph graph.Graph[string, string], links []graph.Edge[string]) graph.Graph[string, string] {
 for _, link := range links { 
  node, err := dspgGraph.Edge(link.Source, link.Target)
  utils.Must(err)
  nestingLevel := node.Properties.Weight
  nestingLevel += 1
  dspgGraph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(nestingLevel))
 }

 return dspgGraph
}
