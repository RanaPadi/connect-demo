The public function receives a `dspgGraph` represented as a graph with nodes and string edges. It performs a series of operations to find all possible paths between sources and destinations in the `dspgGraph` and computes the nesting levels for each valid path.

### Parameters

- `dspgGraph`: A graph with string nodes and string edges.

### Return values

- `allSubPPS`: A slice of `AllSubPPS` objects representing valid paths and their associated graphs, sources and targets.
- `dspgGraph`: The updated graph with all values and nestingLevel.

### Execution Flow

1. The function starts by initializing `sync.WaitGroup` which allows us to use go routines to make the code parallelizable.
   We also initialize `adjacencyMap` of the graph and 2 maps `sourceSet` and `targetSet`.
2. We initialize the counter of the go routine `wg.Add(2)` and 2 functions `SourceSet` and `TargetSet` are executed which will take all PPS paths from the source and target. In addition the `wg.Wait()` will wait for the `wg.Done()` to be done inside the functions which will communicate the end of the go routines.
3. Next, the two for will scroll through all the source and target nodes and check the PPS within a function initialized as a go routine that will allow parallelism.
4. Within the function the intersection of each path from the source and each path from the target will be done to find the common ones.
5. A graph will be created thanks to the `CreateGraphFromPath` function and it will be checked which intersection generated a result.
6. In the positive case the `checkPPS` (which possible to find in documentation) will be performed.
7. If this checkPPS was successful an array will be populated that will count the graph from the relative intersection, source and target of the PPS.
8. Finally given the PPS the nestingLevel will be added to the main graph.
9. Once this fuction and related for will be finished the return operation will be performed which will return array of PPS and the updated graph with related nestingLevel.

**Snippet**

```go

func FindPpsAndNestingLevel(dspgGraph graph.Graph[string, string]) ([]dtos.AllSubPPSDTO, graph.Graph[string, string]) {

 var wg sync.WaitGroup

 adjacencyMap, err := dspgGraph.AdjacencyMap()
 utils.Must(err)

 sourceSet := make(map[string][]dtos.EdgeDTO)
 targetSet := make(map[string][]dtos.EdgeDTO)

 wg.Add(2)

 dspgFunc.SourceSet(&wg, adjacencyMap, sourceSet)
 dspgFunc.TargetSet(&wg, adjacencyMap, targetSet)

 wg.Wait()

 var allSubPPS []dtos.AllSubPPSDTO

 for nodeSource, sourceLinks := range sourceSet {
  for nodeTarget, targetLinks := range targetSet {
   wg.Add(1)

   go func(nodeSource, nodeTarget string, sourceLinks, targetLinks []dtos.EdgeDTO) {
    defer wg.Done()

    intersection := dspgFunc.IntersectionSlices(sourceLinks, targetLinks)

    intersectGraph := dspgFunc.CreateGraphFromPath(intersection)

    graphLen, err := intersectGraph.Size()
    utils.Must(err)
    if graphLen == 0 {
     return
    }

    if dspgFunc.CheckPPS(intersectGraph, nodeSource) {
     allSubPPS = append(allSubPPS, dtos.AllSubPPSDTO{
      Graph:  intersectGraph,
      Source: nodeSource,
      Target: nodeTarget,
     })
     links, err := intersectGraph.Edges()
     utils.Must(err)
     dspgGraph = dspgFunc.UpsertNestingLevel(dspgGraph, links)
    }
   }(nodeSource, nodeTarget, sourceLinks, targetLinks)
  }
 }

 wg.Wait()

 return allSubPPS, dspgGraph
}
