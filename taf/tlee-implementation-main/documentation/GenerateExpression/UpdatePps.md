The public function takes in two parameters: `removedPps` represented as an object of type `AllSubPPS`, and `allSortedPps` represented as a slice of `AllSubPPS`. It updates the `allSortedPps` based on the removal of a specific `AllSubPPS` object.

### Parameters

- `removedPps`: An object of type `AllSubPPS` representing the subgraph to be removed.
- `allSortedPps`: A slice of objects of type `AllSubPPS` representing a list of sorted subgraphs.

### Return Values

- `ppsListUpdated`: A slice of objects of type `AllSubPPS` representing the updated list of sorted subgraphs.

### Flow of Execution

1. The function starts by obtaining the edges of the `removedPps.Graph` and stores them in the `linksRemovedPps` variable.
2. It initializes an empty slice `ppsListUpdated` to store the updated subgraphs.
3. It iterates through each `pps` object in `allSortedPps`:
   - It checks if the `pps` object has nodes corresponding to the `removedPps.Source` and `removedPps.Target`.
   - If the nodes are present, it removes the edges and isolated nodes from the `pps.Graph` based on the `linksRemovedPps`.
   - It adds the `removedPps.Source` and `removedPps.Target` as vertices and creates an edge between them in the `pps.Graph`.
   - The `pps` object is added to the `ppsListUpdated`.
4. The function returns `ppsListUpdated`, which represents the updated list of sorted subgraphs.

**Snippet**

```go
func UpdatePps(removedPps dtos.AllSubPPSDTO, allSortedPps []dtos.AllSubPPSDTO) []dtos.AllSubPPSDTO {
 linksRemovedPps, err := removedPps.Graph.Edges()
 utils.Must(err)

 var ppsListUpdated []dtos.AllSubPPSDTO

 for _, pps := range allSortedPps {
  if HasNode(pps, removedPps.Source) && HasNode(pps, removedPps.Target) {

   for _, link := range linksRemovedPps {
    pps.Graph.RemoveEdge(link.Source, link.Target)
   }
   arrayIsolateNode := dspgFunc.FindIsolatedVertex(pps.Graph)
   for _, node := range arrayIsolateNode {
    pps.Graph.RemoveVertex(node)
   }
   pps.Graph.AddVertex(removedPps.Source)
   pps.Graph.AddVertex(removedPps.Target)
   pps.Graph.AddEdge(removedPps.Source, removedPps.Target)
  }
  ppsListUpdated = append(ppsListUpdated, pps)
 }
 return ppsListUpdated
}
