This private function takes in two parameters: `dspgGraph` and `allPPS`. It operates on a graph represented by `dspgGraph` and a list of `AllSubPPS` objects represented by `allPPS`. This function propagates edge weights from the `dspgGraph` to the `pps.Graph` for all combinations of edges and `pps` objects, effectively updating the weights in the `pps.Graph` based on the data from the `dspgGraph`.

### Parameters

- `dspgGraph`: A graph with string nodes and string edges.
- `allPPS`: A list of objects of type `AllSubPPS`.

### Flow of Execution

1. It starts by obtaining the edges of the `dspgGraph` and stores them in the `links` variable.
2. It then enters a nested loop structure:
   - The outer loop iterates through each `link` in the `links` slice, representing edges in the `dspgGraph`.
   - The inner loop iterates through each `pps` object in the `allPPS` list.
3. For each combination of `link` and `pps`, it updates the weight of the edge in the `pps.Graph` by calling `pps.Graph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(link.Properties.Weight))`. This essentially transfers the weight information from the `dspgGraph` to the corresponding edge in `pps.Graph`.

**Snippet**

```go
func AddNestingLevelToPps(dspgGraph graph.Graph[string, string], allPPS []dtos.AllSubPPSDTO) {
 links, err := dspgGraph.Edges()
 utils.Must(err)
 for _, link := range links {
  for _, pps := range allPPS {
   pps.Graph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(link.Properties.Weight))
  }
 }
}
