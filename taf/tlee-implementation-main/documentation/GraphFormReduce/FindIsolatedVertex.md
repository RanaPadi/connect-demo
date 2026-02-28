The public function takes in a parameter `dspgGraph` represented as a graph with string nodes and string edges. It identifies and returns a list of isolated nodes in the graph, which are nodes that have no incoming or outgoing edges.

### Parameters

- `dspgGraph`: A graph with string nodes and string edges.

### Return Value

- `arrayIsoledNode`: A slice of strings representing the isolated nodes in the graph.

### Flow of Execution

1. The function starts by obtaining the adjacency map of the `dspgGraph` using `dspgGraph.AdjacencyMap()` and handling any potential error using `utils.Must(err)`.
2. It initializes an empty slice `arrayIsoledNode` to store the isolated nodes.
3. It then iterates through each `node` in the adjacency map:
   - It sets a boolean variable `sem` to `false` to track if the node is not connected to any other nodes.
   - It enters a nested loop structure to check for connections between nodes:
     - It iterates through `links` in the adjacency map.
     - For each `link` in `links`, it checks if the `node` is equal to `link`. If they are equal, it sets `sem` to `true` and breaks out of the loop.
   - If `sem` is still `false` after the nested loops, it means the `node` has no connections.
     - It further checks if the length of `AdjacencyMap[node]` is 0, indicating that the node has no outgoing edges.
     - If both conditions are met, it appends the `node` to the `arrayIsoledNode`.
4. The function returns the `arrayIsoledNode`, which contains the isolated nodes in the graph.

**Snippet**

```go

func FindIsolatedVertex(dspgGraph graph.Graph[string, string]) []string {
 AdjacencyMap, err := dspgGraph.AdjacencyMap()
 utils.Must(err)

 var arrayIsoledNode []string

 for node := range AdjacencyMap {
  sem := false
  for _, links := range AdjacencyMap {
   for link := range links {
    if node == link {
     sem = true
     break
    }
   }
  }
  if !sem {
   if len(AdjacencyMap[node]) == 0 {
    arrayIsoledNode = append(arrayIsoledNode, node)
   }
  }
 }

 return arrayIsoledNode
}
