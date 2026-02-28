### **Method `FindSubTrustModels`**

#### **Description**
This function processes a trust graph structure to extract all hierarchical subgraphs (organized by scope) and identifies the root node of the overall structure. It converts the input into an adjacency list, discovers all paths via DFS, groups them by terminal scope, and transforms them into a standardized output format.

#### **Parameters**
- `structure` (type: `trustmodelstructure.TrustGraphStructure`, required):  
  A trust graph structure containing an adjacency list of nodes and their connections. Expected to provide:
    - `AdjacencyList()`: Method returning a slice of graph edges with `SourceNode()` and `TargetNodes()` methods.

#### **Return Type**
- `([]dtos.StructureGraphTAFSingleProp, string)`:
    - **First value**: Slice of `StructureGraphTAFSingleProp` structs, each representing a subgraph with:
        - `AdjacencyList`: Transformed graph edges (`[]VertexEdgeDTO`).
        - `Scope`: Terminal node of the subgraph.
    - **Second value**: Root node (string) of the entire graph (unreferenced by any other node).

#### **Execution Flow**
1. **Graph Conversion**:
    - Converts the input `TrustGraphStructure` into an adjacency list (`map[string][]string`).

2. **Path Discovery**:
    - Uses `FindSubPaths` (DFS-based) to extract all paths in the graph.

3. **Scope-Based Grouping**:
    - Groups paths by their terminal node ("scope") into `ScopeVertexEdgeDTO` structs.

4. **Transformation**:
    - Converts grouped paths into `StructureGraphTAFSingleProp` via `transformPathsToGraph`:
        - Merges duplicate edges.
        - Identifies the root node (unreferenced node).

5. **Output**:
    - Returns the list of scope-specific subgraphs and the root node.

#### **Method Snippet**
```go
func FindSubTrustModels(structure trustmodelstructure.TrustGraphStructure) ([]dtos.StructureGraphTAFSingleProp, string) {
    graph := make(map[string][]string)
    for _, elem := range structure.AdjacencyList() {
        graph[elem.SourceNode()] = elem.TargetNodes()
    }
    subpaths := FindSubPaths(graph)
    var subgraph dtos.StructureGraphTAFSingleProp

    // Create a list to hold ScopeVertexEdgeDTO structs
    var scopeList []dtos.ScopeVertexEdgeDTO
    var subgraphsList []dtos.StructureGraphTAFSingleProp
    
    // Initialize the previous scope to an empty string
    prevScope := ""
    
    // Create currentScope variable
    var currentScope string

    for _, paths := range subpaths {
        for _, path := range paths {
            scopeVertexEdgeDTO := dtos.ScopeVertexEdgeDTO {
                Node:  path[0],
                Links: path[1:],
                Scope: path[len(path)-1],
            }
            // Check if the current scope differs from the previous scope
            if scopeVertexEdgeDTO.Scope != prevScope {
                 // Clear the scopeList if the scope changes
                scopeList = nil
                // Update the previous scope to the current scope
                prevScope = scopeVertexEdgeDTO.Scope
			}
        
            // Append the newly created struct to the scopeList
            scopeList = append(scopeList, scopeVertexEdgeDTO)
            currentScope = scopeVertexEdgeDTO.Scope
        }

        // Get the first return value from transformPathsToGraph
        var firstElementFromtransformPathsToGraph, _ = transformPathsToGraph(scopeList)
        
        // Getting the transformed values from func transformPathsToGraph
        // creating the desired output as for StructureGraphTAFSingleProp
        subgraph = dtos.StructureGraphTAFSingleProp{
            AdjacencyList: firstElementFromtransformPathsToGraph,
            Scope:         currentScope,
        }
        subgraphsList = append(subgraphsList, subgraph)
    }
    // Get the second return value from transformPathsToGraph
    var _, secondElementFromtransformPathsToGraph = transformPathsToGraph(scopeList)
    
    return subgraphsList, secondElementFromtransformPathsToGraph

}

```

#### **Example Output**
For a graph `A -> B -> C` and `A -> D`:
```go
{
	"C": [["A", "B", "C"]],
	"D": [["A", "D"]]
}
```
