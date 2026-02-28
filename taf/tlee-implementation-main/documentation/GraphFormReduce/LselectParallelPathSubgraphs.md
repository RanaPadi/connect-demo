## Function `LselectParallelPathSubgraphs`

### Description
The function `LselectParallelPathSubgraphs` selects parallel paths in a directed graph (`g`) and identifies subgraphs that meet specific conditions. For each selected parallel path, it retrieves information such as the source and target vertices, vertex degrees, and the total number of vertices. Subgraphs meeting certain criteria (minimum vertex degree and minimum number of vertices) are stored along with their associated information in a map (`allSubPPS`). The function returns the original graph (`g`) with updated nesting levels and a slice of structures (`dtos.AllSubPPSDTO`) representing the selected parallel path subgraphs.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The directed graph in which parallel paths are selected.

### Return Type
- The original graph (`g`) with updated nesting levels.
- A slice of structures (`dtos.AllSubPPSDTO`) representing the selected parallel path subgraphs, including the subgraph, source, target, etc.

### Execution Flow
1. Obtain a list of selected parallel paths using the `LselectParallelPaths` function.
2. Initialize a map (`allSubPPS`) to store selected subgraphs along with associated information.
3. Use goroutines to concurrently fetch source and target vertices, vertex degrees, and the number of vertices for each selected subgraph.
4. Check specific conditions (minimum vertex degree and minimum number of vertices) for each subgraph.
5. Store information for subgraphs meeting the conditions in the map (`allSubPPS`), and update nesting levels in the original graph.
6. Convert the map values to a slice (`allSubPPSSet`) for the final result.
7. Return the original graph with updated nesting levels and the slice of selected parallel path subgraphs.

### Snippet
```go
func LselectParallelPathSubgraphs(g graph.Graph[string, string]) (graph.Graph[string, string], []dtos.AllSubPPSDTO) {
	list := LselectParallelPaths(g)

	allSubPPS := make(map[string]dtos.AllSubPPSDTO)

	var wg sync.WaitGroup

	for _, graph := range list {

		var source string
		var target string

		var vio []int
		var lenght int

		wg.Add(4)

		go func() {
			defer wg.Done()
			source = GetSourceOriginalGraph(graph)
		}()

		go func() {
			defer wg.Done()
			target = GetTargetOriginalGraph(graph)
		}()

		go func() {
			defer wg.Done()
			vio = VertexDegree(graph)
		}()

		go func() {
			defer wg.Done()
			lenght = len(VertexList(graph))
		}()

		wg.Wait()

		if MinVertexDegree(vio) >= 1 && lenght >= 2 {
			key := fmt.Sprint(source + "" + target)
			allSubPPS[key] = dtos.AllSubPPSDTO{
				Graph:  graph,
				Source: source,
				Target: target,
			}
			links, err := graph.Edges()
			utils.Must(err)
			g = UpsertNestingLevel(g, links)
		}
	}

	var allSubPPSSet []dtos.AllSubPPSDTO
	for _, value := range allSubPPS {
		allSubPPSSet = append(allSubPPSSet, value)
	}

	return g, allSubPPSSet
}

