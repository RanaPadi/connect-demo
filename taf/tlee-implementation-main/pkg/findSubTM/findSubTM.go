package findSubTM

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"github.com/vs-uulm/taf-tlee-interface/pkg/trustmodelstructure"
)

/*
Package provides functionality for graph-based trust model analysis.

- The `Graph` type represents an adjacency list where nodes (strings) are mapped to lists of connected nodes.
- The `DFS` function performs a depth-first search to identify all subgraphs by exploring paths from each node.
- The `FindSubPaths` function iterates through the graph, using DFS to extract and group subpaths by their final node (scope).
- The `transformPathsToGraph` function converts a list of paths into a structured graph representation using `VertexEdgeDTO`, ensuring unique links and identifying the root node.
- The `FindSubTrustModels` function processes a `TrustGraphStructure` by converting its adjacency list into subpaths and then transforming these subpaths into structured subgraphs for trust model analysis.

Overall, this code extracts, transforms, and structures subgraphs from a given trust model representation.
*/

// Graph represents the adjacency list of a graph
type Graph map[string][]string

/*
		Purpose: Performs depth-first search (DFS) to explore all paths in the graph starting from node.

		Behavior:
	 	1. Appends the current node to the path.
		2. If the node is a leaf (no outgoing edges), adds the current path to paths under the leaf node's scope (last node in the path)
	 	3. Recursively visits all neighbors of the current node.
*/
func DFS(graph Graph, node string, path []string, paths map[string][][]string) {
	path = append(path, node)

	// If the node has no outgoing edges, it's a leaf node, so add the current path to the paths slice
	if len(graph[node]) == 0 {
		scope := path[len(path)-1]
		paths[scope] = append(paths[scope], append([]string{}, path...))
		return
	}

	// Recursively explore each neighbor
	for _, neighbor := range graph[node] {
		DFS(graph, neighbor, append([]string{}, path...), paths)
	}
}

//

/*
	Purpose: Finds all subpaths in the graph by initiating DFS from each unvisited node.
	Steps:
		1. Initializes paths (to store subpaths) and visited (to track visited nodes).
		2. Iterates over each node in the graph, running DFS if the node hasn't been visited.
	Returns: A map of subpaths grouped by their scope */

func FindSubPaths(graph Graph) map[string][][]string {
	paths := make(map[string][][]string)
	visited := make(map[string]bool)

	// Perform DFS from each unvisited node
	for node := range graph {
		if !visited[node] {
			DFS(graph, node, []string{}, paths)
		}
		visited[node] = true
	}
	return paths
}

/*
	Purpose: Transforms a list of ScopeVertexEdgeDTO (scoped paths) into a simplified adjacency list (VertexEdgeDTO) and identifies the root node.
	Behavior:
		1. Skips nodes with no links.
		2. For each node in scopeList:
			* If the node is new, initializes its entry in transformedGraph.
			* Otherwise, merges links while avoiding duplicates.
		3. Identifies the root node (a node not referenced by any other node).
	Returns:
		1. A slice of VertexEdgeDTO representing the transformed graph.
		2. The root node (string) of the graph.*/

func transformPathsToGraph(scopeList []dtos.ScopeVertexEdgeDTO) ([]dtos.VertexEdgeDTO, string) {

	transformedGraph := dtos.VertexEdgeDTO{}
	transformedGraphComplete := []dtos.VertexEdgeDTO{}
	linkedNodes := make(map[string]bool)

	for i := range scopeList {

		if len(scopeList[i].Links) == 0 {
			continue
		}
		// Populate the set with linked nodes
		for _, link := range scopeList[i].Links {
			linkedNodes[link] = true
		}

		// Check if the first element has a concrete value
		if transformedGraph.Node != scopeList[i].Node {
			transformedGraph.Node = scopeList[i].Node
			transformedGraph.Links = []string{scopeList[i].Links[0]}
			transformedGraphComplete = append(transformedGraphComplete, transformedGraph)

		} else {

			// Assuming scopeList[i].Links[0] is of type Link
			linkToAdd := scopeList[i].Links[0]
			linkExists := false

			// Check if linkToAdd already exists in transformedGraph.Links
			for _, existingLink := range transformedGraph.Links {
				if existingLink == linkToAdd {
					linkExists = true
					break
				}
			}

			// If linkToAdd doesn't exist in transformedGraph.Links, append it
			if !linkExists {
				transformedGraph.Links = append(transformedGraph.Links, linkToAdd)
			}

			for index, value := range transformedGraphComplete {
				if value.Node == scopeList[i].Node {
					transformedGraphComplete[index].Links = transformedGraph.Links
				}
			}

		}
	}
	// Identify the root node
	rootNode := ""
	for _, node := range transformedGraphComplete {
		if _, exists := linkedNodes[node.Node]; !exists {
			rootNode = node.Node
			break
		}
	}

	return transformedGraphComplete, rootNode

}

/*
	Purpose: Processes a trust graph structure into subgraphs grouped by scope and extracts the root node.
	Workflow:
		1. Converts the input TrustGraphStructure into an adjacency list (graph).
		2. Uses FindSubPaths to discover all subpaths in the graph.
		3. Groups paths by their terminal scope (ScopeVertexEdgeDTO).
		4. Transforms each scoped group into a StructureGraphTAFSingleProp using transformPathsToGraph.
	Returns:
		* A slice of StructureGraphTAFSingleProp (subgraphs per scope).
		* The root node of the entire graph.*/

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
			scopeVertexEdgeDTO := dtos.ScopeVertexEdgeDTO{
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
