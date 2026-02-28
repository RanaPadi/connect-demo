package dspgFunc

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/config"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/utils"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

/* Various functions needed for the Graph Reduction. */

const minLinkToNode = 2

var debugFoldername = config.OutputPath + "/debug/"

var _ = os.MkdirAll(debugFoldername, os.ModePerm)

/*
The public function takes in a parameter nodeList represented as a map of strings to VertexProperties.
It calculates and returns the count of nodes in the nodeList with an InFlow value of 0.999.

The function initializes a variable count to 0 to keep track of the count.
It iterates through each node in the nodeList:

It checks if the InFlow value of the node is greater than 0.999.
If the condition is met, it increments the count.

The function returns the count, which represents the count of nodes in the nodeList with an InFlow value of 0.999.
*/
func countInFlow(nodeList map[string]dtos.VertexPropertiesDTO) int {
	var count int = 0
	for _, links := range nodeList {
		if links.InFlow > 0.9999 {
			count++
		}
	}
	return count
}

/*
The public function takes in three parameters: adj, listY, and nodeList.
It calculates the sum of outflows for a list of nodes listY within a graph represented by the adjacency map adj, using the information in the nodeList.

The function initializes a float64 variable inFlowNumber to zero, which will store the sum of outflows.
It iterates through each element el in the listY:

For each el, it accesses the corresponding node's properties in the nodeList.
It adds the OutFlow property of the node to the inFlowNumber.

After iterating through all elements in listY, the function returns the inFlowNumber, which represents the sum of outflows for the specified nodes.
*/
func sumOutFlow(listY map[string]struct{}, nodeList map[string]dtos.VertexPropertiesDTO) float64 {
	outFlowSum := 0.0

	for el := range listY {
		if nodeProps, exists := nodeList[el]; exists {
			outFlowSum += nodeProps.OutFlow
		}
	}

	return outFlowSum
}

/*
The public function takes in two parameters: adj represented as a map of strings to maps of strings to Edge[string], and x as a string.
It identifies and returns a list of nodes in the adj map that are connected to the node x via edges.

The function initializes an empty slice listY to store the connected nodes.
It then iterates through each node in the adj map:

For each node, it iterates through the edjes (edges) associated with that node.
For each edge in edjes, it checks if edge is equal to the input x.
If edge is equal to x, it means that node is connected to x via an edge, and it appends node to the listY.

After iterating through all nodes and edges, the function returns the listY, which contains the nodes connected to x via edges in the adj map.
*/
func findY(adj map[string]map[string]graph.Edge[string], x string) map[string]struct{} {
	adjacentNodes := make(map[string]struct{})

	for node, edges := range adj {
		if _, exists := edges[x]; exists {
			adjacentNodes[node] = struct{}{}
		}
	}

	return adjacentNodes
}

/*
The public function takes in a parameter adjacencyMap, represented as a map of strings to maps of strings to Edge[string].
It identifies and returns a list of source nodes within the graph represented by the adjacencyMap.
A source node is defined as a node with a minimum number of outgoing links, which is determined by the minLinkToNode threshold.

The function initializes an empty slice sources to store the source nodes.
It then iterates through each node in the adjacencyMap:

For each node, it checks the number of outgoing links associated with it by examining the length of the links map.
If the number of outgoing links is greater than or equal to the minLinkToNode threshold, it considers the node as a source node and appends it to the sources slice.

After iterating through all nodes, the function returns the sources slice, which contains the source nodes within the graph.
*/
func getSource(adjacencyMap map[string]map[string]graph.Edge[string]) []string {
	var sources []string
	for node, links := range adjacencyMap {
		if len(links) >= minLinkToNode {
			sources = append(sources, node)
		}
	}
	return sources
}

/*
The public function takes in a parameter adjacencyMap, represented as a map of strings to maps of strings to Edge[string].
It identifies and returns a list of target nodes within the graph represented by the adjacencyMap.
A target node is defined as a node with a minimum number of incoming links, which is determined by the minLinkToNode threshold.

The function initializes an empty slice targets to store the target nodes.
It then iterates through each node in the adjacencyMap:

For each node, it calculates the number of incoming links (edges) associated with it by calling the linksCounter function.
If the number of incoming links is greater than or equal to the minLinkToNode threshold, it considers the node as a target node and appends it to the targets slice.

After iterating through all nodes, the function returns the targets slice, which contains the target nodes within the graph.
*/
func getTarget(adjacencyMap map[string]map[string]graph.Edge[string]) []string {
	var targets []string
	for node := range adjacencyMap {
		link := linksCounter(adjacencyMap, node)
		if link >= minLinkToNode {
			targets = append(targets, node)
		}
	}
	return targets
}

/*
The private function takes in two parameters: adjacencyMap and node.
It is responsible for counting the number of incoming links (edges) to a specified node within a graph represented by the adjacencyMap.

The function initializes an integer variable linksCount to zero, which will store the count of incoming links.
It then iterates through each entry in the adjacencyMap:

For each entry, it represents a node (key) and its associated outgoing links (value).
It further iterates through the outgoing links for that node.
If a link in the outgoing links matches the specified node, it increments the linksCount by 1.

After iterating through all entries in the adjacencyMap, the function returns the linksCount, which represents the count of incoming links to the specified node.
*/
func linksCounter(adjacencyMap map[string]map[string]graph.Edge[string], node string) int {
	var linksCount = 0
	for _, links := range adjacencyMap {
		for link := range links {
			if link == node {
				linksCount += 1
			}
		}
	}
	return linksCount
}

/*
The public function takes in two parameters: dspgGraph and links.
It updates the nesting level of specified links within a directed graph represented by dspgGraph.

The function iterates through each link in the links slice:

For each link, it retrieves the corresponding edge from the dspgGraph using link.Source and link.Target.
It extracts the current nesting level from the edge's properties.
It increments the nesting level by 1 to update it.
It updates the edge's weight with the new nesting level using dspgGraph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(nestingLevel)).
*/
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

/*
The public function takes in three parameters: g represented as a graph with string nodes and string edges,
node representing the node to be added, and links representing a list of nodes to be connected to the node.
It adds the specified node and establishes edges between the node and the nodes in the links list.

The function starts by adding the node to the graph g using g.AddVertex(node).
It then iterates through each link in the links list:

It adds the link as a vertex to the graph g using g.AddVertex(link).
It establishes an edge between the node and the link using g.AddEdge(node, link).

The function completes the addition of the node and the establishment of edges between the node and the nodes in the links list.
*/
func AddVertexAndEdge(g graph.Graph[string, string], node string, links []string) graph.Graph[string, string] {
	g.AddVertex(node)
	for _, link := range links {
		g.AddVertex(link)
		g.AddEdge(node, link)
	}

	return g
}

/*
The public function takes in a parameter dspgGraph represented as a graph with string nodes and string edges.
It identifies and returns a list of isolated nodes in the graph, which are nodes that have no incoming or outgoing edges.

The function starts by obtaining the adjacency map of the dspgGraph using dspgGraph.AdjacencyMap() and handling any potential error using utils.Must(err).
It initializes an empty slice arrayIsoledNode to store the isolated nodes.
It then iterates through each node in the adjacency map:

It sets a boolean variable sem to false to track if the node is not connected to any other nodes.
It enters a nested loop structure to check for connections between nodes:

It iterates through links in the adjacency map.
For each link in links, it checks if the node is equal to link. If they are equal, it sets sem to true and breaks out of the loop.

If sem is still false after the nested loops, it means the node has no connections.

It further checks if the length of AdjacencyMap[node] is 0, indicating that the node has no outgoing edges.
If both conditions are met, it appends the node to the arrayIsoledNode.

The function returns the arrayIsoledNode, which contains the isolated nodes in the graph.
*/
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

/*
The public function takes in a parameter dspgGraph, represented as a directed graph with string nodes and string edges.
It identifies and returns the source node within the original graph represented by dspgGraph. The source node is defined as a node with no incoming edges.

The function starts by obtaining the adjacency map of the dspgGraph using dspgGraph.AdjacencyMap().
It initializes an empty string variable source to store the source node.
It also initializes a boolean variable sem to false.
The function then iterates through each node1 in the adjacency map:

For each node1, it sets sem to false.
It then iterates through all the links (outgoing links) associated with each node in the adjacency map.
If it finds any link that matches the current node1, it sets sem to true, indicating that the node1 has incoming edges.
If sem remains false after checking all links, it means that node1 has no incoming edges, and it sets the source variable to node1 and breaks out of the loop.

After the loop, the function returns the source string, which represents the source node within the original graph.
*/
func GetSourceOriginalGraph(dspgGraph graph.Graph[string, string]) string {

	adjacencyMap, err := dspgGraph.AdjacencyMap()
	utils.Must(err)

	var source string
	var sem bool

	for node1 := range adjacencyMap {
		sem = false
		for _, links := range adjacencyMap {
			for link := range links {
				if node1 == link {
					sem = true
				}
			}
		}
		if !sem {
			source = node1
			break
		}
	}

	return source

}

/*
The public function takes in a parameter dspgGraph, represented as a directed graph with string nodes and string edges.
It identifies and returns the target node within the original graph represented by dspgGraph. The target node is defined as a node with no outgoing edges.

The function starts by obtaining the adjacency map of the dspgGraph using dspgGraph.AdjacencyMap().
It initializes an empty string variable target to store the target node.
The function then iterates through each node in the adjacency map:

For each node, it checks the length of the associated links (outgoing links). If the length is zero, it indicates that the node has no outgoing edges.
In such cases, it sets the target variable to the current node.

After the loop, the function returns the target string, which represents the target node within the original graph.
*/
func GetTargetOriginalGraph(dspgGraph graph.Graph[string, string]) string {

	adjacencyMap, err := dspgGraph.AdjacencyMap()
	utils.Must(err)

	var target string

	for node, links := range adjacencyMap {
		if len(links) == 0 {
			target = node
		}
	}

	return target
}

/*
This function appears to construct a set of arcs (edges) starting from a set of source nodes in a graph represented by `adjacencyMap`.
The function uses the `wg` parameter of type `*sync.WaitGroup` to handle concurrent goroutine checking.
Source nodes are extracted from the graph, and for each source node, all arcs exiting the node are found using the `findAllEdgesFromSourceIterative` function.
The resulting arcs are then stored in `sourceSet`.
*/
func SourceSet(wg *sync.WaitGroup, adjacencyMap map[string]map[string]graph.Edge[string], sourceSet map[string][]dtos.EdgeDTO) {
	defer wg.Done()

	sources := getSource(adjacencyMap)
	for _, node := range sources {
		edgeFromSource := findAllEdgesFromSourceIterative(adjacencyMap, node)
		sourceSet[node] = append(sourceSet[node], edgeFromSource...)
	}
}

/*
This function is similar to the previous one, but instead of looking for arcs departing from source nodes,
it looks for arcs arriving at target nodes in a graph represented by adjacencyMap.
It uses the parameter wg to handle the control of concurrent goroutines.
Target nodes are extracted from the graph, and for each target node, all arcs arriving at the node are found using the findAllEdgesToTargetIterative function.
The resulting arcs are stored in targetSet.
*/
func TargetSet(wg *sync.WaitGroup, adjacencyMap map[string]map[string]graph.Edge[string], targetSet map[string][]dtos.EdgeDTO) {
	defer wg.Done()

	targets := getTarget(adjacencyMap)
	for _, node := range targets {
		edgeFromTarget := findAllEdgesToTargetIterative(adjacencyMap, node)
		targetSet[node] = append(targetSet[node], edgeFromTarget...)
	}
}

/*
This function iteratively searches for all arcs exiting from a specific source node (source) in a graph represented by dspgGraph.
It uses an iterative depth visit with a stack to explore the graph. The resulting arcs are stored as dtos.EdgeDTO objects in an array and returned as output.
*/
func findAllEdgesFromSourceIterative(dspgGraph map[string]map[string]graph.Edge[string], source string) []dtos.EdgeDTO {
	visited := make(map[string]bool)
	edgesFromSource := []dtos.EdgeDTO{}
	visitedEdges := make(map[dtos.EdgeDTO]bool)
	stack := []string{source}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visited[node] = true

		for neighbor := range dspgGraph[node] {
			edge := dtos.EdgeDTO{FromNode: node, ToNode: neighbor}

			if !visitedEdges[edge] {
				edgesFromSource = append(edgesFromSource, edge)
				visitedEdges[edge] = true
			}

			if !visited[neighbor] {
				stack = append(stack, neighbor)
			}
		}
	}

	return edgesFromSource
}

/*
This function is similar to the previous one but looks for all arcs coming to a specific
target node (target) in the graph represented by dspgGraph. It also uses an iterative depth visit with a stack to explore the graph.
The resulting arcs are stored as dtos.EdgeDTO objects in an array and returned as output.
*/
func findAllEdgesToTargetIterative(dspgGraph map[string]map[string]graph.Edge[string], target string) []dtos.EdgeDTO {
	edgesToTarget := []dtos.EdgeDTO{}
	stack := []string{target}
	visited := make(map[string]bool)
	visitedEdges := make(map[dtos.EdgeDTO]bool)

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[node] {
			continue
		}

		visited[node] = true

		for parentNode, edges := range dspgGraph {
			for neighbor := range edges {
				if neighbor == node {
					edge := dtos.EdgeDTO{FromNode: parentNode, ToNode: node}

					// Verifica se l'edge è già stato visitato
					if !visitedEdges[edge] {
						edgesToTarget = append(edgesToTarget, edge)
						visitedEdges[edge] = true
					}

					stack = append(stack, parentNode)
				}
			}
		}
	}

	return edgesToTarget
}

/*
This function computes the intersection of two slices of dtos.EdgeDTO objects.
It receives two slices, slice1 and slice2, and returns a new slice containing only the common elements between the two inputs.
The function uses a map to keep track of elements in slice1, then scrolls slice2 to find elements that also exist in slice1 and adds them to the output.
The output is an array containing the common elements between slice1 and slice2.
*/
func IntersectionSlices(slice1, slice2 []dtos.EdgeDTO) []dtos.EdgeDTO {
	set := make(map[dtos.EdgeDTO]struct{})
	for _, item := range slice1 {
		set[item] = struct{}{}
	}

	var intersection []dtos.EdgeDTO
	for _, item := range slice2 {
		if _, exists := set[item]; exists {
			intersection = append(intersection, item)
		}
	}

	return intersection
}

/*
The public function receives an all paths parameter represented as a slice of string slices.
It creates and returns a new directed graph with string nodes and string edges, based on the provided list of paths.

The function starts by creating a new directed graph graph using graph.New(graph.StringHash, graph.Directed(), graph.Acyclic()).
Then iterate through each slice of nodes in allPaths:

For each node in the allPaths slice, add the node as a vertex to the graph using graph.AddVertex(nodes.FromNode).
For each node in the allPaths slice, adds the node as a vertex to the graph using graph.AddVertex(nodes.ToNode).
For each node in the allPaths slice, it adds the edge as the link between the two nodes using graph.AddEdge(nodes.FromNode, nodes.ToNode).

The function completes the creation of the graph based on the list of paths and returns the resulting graph.
*/
func CreateGraphFromPath(allPaths []dtos.EdgeDTO) graph.Graph[string, string] {
	intersectGraph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())

	for _, nodes := range allPaths {
		intersectGraph.AddVertex(nodes.FromNode)
		intersectGraph.AddVertex(nodes.ToNode)
		intersectGraph.AddEdge(nodes.FromNode, nodes.ToNode)
	}

	return intersectGraph
}

/*
This public function accepts two parameters: subgraph and source .
It operates on a subgraph represented as a graph with nodes and string edges and checks the feasibility of a path from the source node within the subgraph.
The function calls another function calculateWeights that will return the computed flow.
The calculateWeights function starts by obtaining the adjacency map of the subgraph and stores it in SubGraphAdjacencyMap.
It initializes two maps, nodeList and weight, to keep track of vertex properties during execution.
For each node in the SubGraphAdjacencyMap, initialize the properties in weight with zero input and output flow.
Initialize the source node in weight with input flow equal to 1 and output flow equal to the division between 1 and the length of SubGraphAdjacencyMap[source] which would be the number of nodes exiting the source

We populate our string array nodeList where we save all the children of the source.
We create a loop that will run until nodeList is empty.
Iteratively process the nodes using a loop:
We initialize an empty nodeListSon array.
For each node in nodeList, compute in-flow and out-flow properties based on adjacent nodes and update weight using the updateWeights function.
Iteratively process the nodes using another loop:
For each node in the nodeList we go to find the children and if inside our nodeListSon array they do not exist, check done with the containsNode function we go to add them to nodeListSon
At this point we set our nodeList with the previously saved children in nodeListSon and repeat the steps from step 6. until the last one has no more children and nodeListSon is empty.
Once the execution is finished the array containing the stream in all nodes from the source to the target will be returned.
Through the countInFlow function we are going to see how many nodes have a flow greater than 0.999999, excluding source and target, and return true if they are <= 2 or false.
*/

func CheckPPS(subgraph graph.Graph[string, string], source string) bool {
	flow := calculateWeights(subgraph, source)
	return (len(flow) > 0 && countInFlow(flow) <= 2)
}

func calculateWeights(subgraph graph.Graph[string, string], source string) map[string]dtos.VertexPropertiesDTO {

	SubGraphAdjacencyMap, err := subgraph.AdjacencyMap()
	utils.Must(err)

	weight := make(map[string]dtos.VertexPropertiesDTO)

	for node := range SubGraphAdjacencyMap {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: 0, OutFlow: 0}
	}

	weight[source] = dtos.VertexPropertiesDTO{InFlow: 1, OutFlow: 1 / float64(len(SubGraphAdjacencyMap[source]))}

	var nodeList []string

	for node := range SubGraphAdjacencyMap[source] {
		nodeList = append(nodeList, node)
	}

	for len(nodeList) != 0 {
		var nodeListSon []string
		for _, node := range nodeList {
			weight = updateWeights(weight, SubGraphAdjacencyMap, node)
		}
		for _, nodes := range nodeList {
			for node := range SubGraphAdjacencyMap[nodes] {
				if !ContainsNode(nodeListSon, node) {
					nodeListSon = append(nodeListSon, node)
				}
			}
		}

		nodeList = nodeListSon
	}

	return weight
}

func ContainsNode[T string | dtos.EdgeDTO](slice []T, nodeToFind T) bool {
	for _, node := range slice {
		if node == nodeToFind {
			return true
		}
	}
	return false
}

func updateWeights(weight map[string]dtos.VertexPropertiesDTO, SubGraphAdjacencyMap map[string]map[string]graph.Edge[string], node string) map[string]dtos.VertexPropertiesDTO {
	findY := findY(SubGraphAdjacencyMap, node)
	vertexInWeight := sumOutFlow(findY, weight)
	outDegree := float64(len(SubGraphAdjacencyMap[node]))

	if outDegree == 0 {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: vertexInWeight, OutFlow: 0}
	} else {
		weight[node] = dtos.VertexPropertiesDTO{InFlow: vertexInWeight, OutFlow: vertexInWeight / outDegree}
	}

	return weight
}

func dfsAllPaths(graph map[string]map[string]graph.Edge[string], current string, end string, visited map[string]bool, path []string, paths *[][]string) {
	visited[current] = true
	path = append(path, current)

	if current == end {
		*paths = append(*paths, append([]string{}, path...))
	} else {
		for neighbor := range graph[current] {
			if !visited[neighbor] {
				dfsAllPaths(graph, neighbor, end, visited, path, paths)
			}
		}
	}

	visited[current] = false
}

func GetAllPaths(adj map[string]map[string]graph.Edge[string], start string, end string) [][]string {
	visited := make(map[string]bool)
	paths := [][]string{}
	path := []string{}
	dfsAllPaths(adj, start, end, visited, path, &paths)
	return paths
}

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

func ContainsAllNode(slice []dtos.EdgeDTO, nodeToFind dtos.EdgeDTO) bool {
	for _, node := range slice {
		if node.FromNode == nodeToFind.FromNode && node.ToNode == nodeToFind.ToNode {
			return true
		}
	}
	return false
}

func ToEdgeFromNodeAndLinks(input []dtos.EdgeDTO) []map[string]interface{} {
	nodes := make(map[string]map[string]bool)

	for _, edge := range input {
		fromNode := edge.FromNode
		toNode := edge.ToNode

		if _, ok := nodes[fromNode]; !ok {
			nodes[fromNode] = make(map[string]bool)
		}
		nodes[fromNode][toNode] = true
	}

	result := []map[string]interface{}{}
	for node, links := range nodes {
		nodeData := map[string]interface{}{
			"node":  node,
			"links": make([]string, 0, len(links)),
		}

		for link := range links {
			nodeData["links"] = append(nodeData["links"].([]string), link)
		}

		result = append(result, nodeData)
	}

	return result
}

func DspgEdgeCheck(synthesizingGraph dtos.SynthesizingGraph, nodeSource, nodeTarget string) []string {
	path := veryFindPath(synthesizingGraph, nodeSource, nodeTarget)
	if path == nil {
		return nil
	}
	if !containsTRY(synthesizingGraph, nodeSource, nodeTarget) {
		return nil
	}
	if !containsTRY(synthesizingGraph, nodeTarget, nodeSource) {
		return nil
	}
	return path
}

func veryFindPath(synthesizingGraph dtos.SynthesizingGraph, startNode, endNode string) []string {
	grafh := CreateGraphFromPath(synthesizingGraph.Edges)

	var path []string

	controller := func(node string) bool {
		path = append(path, node)
		return node == endNode
	}

	checked := graph.DFS(grafh, startNode, controller)
	if checked == nil {
		return path
	}

	return nil
}

func containsTRY(synthesizingGraph dtos.SynthesizingGraph, key string, node string) bool {
	if _, exists := synthesizingGraph.NodeToPPS[key]; !exists {
		return true
	}

	pps := synthesizingGraph.NodeToPPS[key]
	if node == pps.FromNode || node == pps.ToNode {
		return true
	}

	tmp := allPaths(synthesizingGraph, pps.FromNode, pps.ToNode)
	for _, tmpNode := range tmp {
		if ContainsNode(tmpNode, node) {
			return true
		}
	}

	return false

}

func DFS(synthesizingGraph dtos.SynthesizingGraph, s string, d string, visited map[string]bool, path []string, paths *[][]string) {
	adjMatrix := synthesizingGraph.Adj
	visited[s] = true
	path = append(path, s)

	if s == d {
		*paths = append(*paths, path)
	} else {
		for neighbor := range adjMatrix[s] {
			if !visited[neighbor] {
				DFS(synthesizingGraph, neighbor, d, visited, path, paths)
			}
		}
	}

	delete(visited, s)
	path = path[:len(path)-1]
}

func allPaths(synthesizingGraph dtos.SynthesizingGraph, s string, d string) [][]string {
	visited := make(map[string]bool)
	paths := [][]string{}
	path := []string{}

	DFS(synthesizingGraph, s, d, visited, path, &paths)

	return paths
}

func GetNodes(processingEdges []dtos.EdgeDTO) []string {
	n := len(processingEdges)
	nodes := make([]string, 0, n-1)

	for i := 1; i < n; i++ {
		nodes = append(nodes, processingEdges[i].FromNode)
	}

	return nodes
}

func VertexList(g graph.Graph[string, string]) []string {
	var vertexList []string

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	for vertex := range adj {
		vertexList = append(vertexList, vertex)
	}

	return vertexList
}

func VertexInDegreeCheck(g graph.Graph[string, string], vertex string) bool {
	adj, err := g.AdjacencyMap()
	utils.Must(err)

	for _, edges := range adj {
		for edge := range edges {
			if edge == vertex {
				return true
			}
		}
	}

	return false
}

func VertexInDegreeList(g graph.Graph[string, string], vertex string) []string {

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var vertexList []string

	for node, edges := range adj {
		for edge := range edges {
			if edge == vertex {
				vertexList = append(vertexList, node)
			}
		}
	}

	return vertexList
}

func VertexOutDegreeCheck(g graph.Graph[string, string], vertex string) bool {
	adj, err := g.AdjacencyMap()
	utils.Must(err)

	return len(adj[vertex]) > 0
}

func VertexOutDegreeList(g graph.Graph[string, string], vertex string) []string {

	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var vertexList []string

	for vertex := range adj[vertex] {
		vertexList = append(vertexList, vertex)
	}

	return vertexList
}

func DiffVertexDegree(g graph.Graph[string, string], subgraph graph.Graph[string, string]) int {
	var v []string

	for _, vertex := range VertexList(subgraph) {
		if VertexInDegreeCheck(subgraph, vertex) && VertexOutDegreeCheck(subgraph, vertex) {
			v = append(v, vertex)
		}
	}

	var diff []int

	var wg sync.WaitGroup

	var vi []string
	var vo []string

	for _, vertex := range v {
		wg.Add(2)

		go func() {
			defer wg.Done()
			vi = VertexInDegreeList(g, vertex)
		}()

		go func() {
			defer wg.Done()
			vo = VertexOutDegreeList(g, vertex)
		}()

		wg.Wait()

		diff = append(diff, len(vi)-len(vo))
	}

	var result int
	for _, value := range diff {
		result += value
	}

	return result
}

func LocalDSPGTest(g graph.Graph[string, string], selectedPps []graph.Graph[string, string]) bool {
	var results int

	for _, subgraph := range selectedPps {
		results += DiffVertexDegree(g, subgraph)
	}

	return len(selectedPps) > 0 && results == 0
}

func Intersect(sets ...[]string) []string {
	if len(sets) == 0 {
		return nil
	}

	intersection := make(map[string]bool)
	result := []string{}

	for _, element := range sets[0] {
		intersection[element] = true
	}

	for _, set := range sets[1:] {
		currentSet := make(map[string]bool)

		for _, element := range set {
			currentSet[element] = true
		}

		for element := range intersection {
			if !currentSet[element] {
				delete(intersection, element)
			}
		}
	}

	for element := range intersection {
		result = append(result, element)
	}

	return result
}

func LisPPSQ(g graph.Graph[string, string], paths [][]string) bool {
	intersectionResult := Intersect(paths...)
	return len(intersectionResult) == 2 && len(paths) >= 2
}

func LselectParallelPaths(g graph.Graph[string, string]) []graph.Graph[string, string] {
	vertexList := VertexList(g)

	var sourceSet []string
	var targetSet []string

	var list []graph.Graph[string, string]

	var wg sync.WaitGroup

	var resultsInDegree []string
	var resultsOutDegree []string

	for _, vertex := range vertexList {
		wg.Add(2)

		go func() {
			defer wg.Done()
			resultsInDegree = VertexInDegreeList(g, vertex)
		}()

		go func() {
			defer wg.Done()
			resultsOutDegree = VertexOutDegreeList(g, vertex)
		}()

		wg.Wait()
		if len(resultsInDegree) >= 2 {
			targetSet = append(targetSet, vertex)
		}
		if len(resultsOutDegree) >= 2 {
			sourceSet = append(sourceSet, vertex)
		}
	}

	tuples := utils.CartesianProduct(sourceSet, targetSet)

	for _, value := range tuples {
		adj, err := g.AdjacencyMap()
		utils.Must(err)
		values := GetAllPaths(adj, value[0], value[1])
		utils.Must(err)

		if len(values) > 0 && LisPPSQ(g, values) {

			list = append(list, Subgraph(g, values))
		}
	}

	return list
}

func VertexDegree(g graph.Graph[string, string]) []int {
	var vertexDegree []int

	v := VertexList(g)

	var wg sync.WaitGroup

	var vi []string
	var vo []string

	for _, vertex := range v {
		wg.Add(2)
		go func() {
			defer wg.Done()
			vi = VertexInDegreeList(g, vertex)
		}()

		go func() {
			defer wg.Done()
			vo = VertexOutDegreeList(g, vertex)
		}()
		wg.Wait()

		vertexDegree = append(vertexDegree, len(vi)+len(vo))
	}

	return vertexDegree
}

func MinVertexDegree(vertexList []int) int {
	min := vertexList[0]
	for _, value := range vertexList {
		if value < min {
			min = value
		}
	}

	return min
}

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

func Subgraph(g graph.Graph[string, string], vertices [][]string) graph.Graph[string, string] {
	subgraph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())
	adj, err := g.AdjacencyMap()
	utils.Must(err)

	var NodeMap map[string]string = make(map[string]string)

	for _, vertice := range vertices {
		for _, v := range vertice {
			NodeMap[v] = v
			subgraph.AddVertex(v)
		}
	}

	for _, v := range NodeMap {
		for _, edge := range adj[v] {
			subgraph.AddEdge(v, edge.Target)
		}

	}

	return subgraph
}
func DebugDrawGraph(g graph.Graph[string, string], name string) {
	graph, err := os.Create(debugFoldername + name + ".gv")
	utils.Must(err)
	err = draw.DOT(g, graph)
	utils.Must(err)

	exec.Command("dot", "-Tsvg", "-O", config.OutputPath+"/debug/"+name+".gv").Output()
}

func GenerateRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func RemoveSelfLoop(dspgGraph graph.Graph[string, string]) graph.Graph[string, string] {
	adj, _ := dspgGraph.AdjacencyMap()

	for node := range adj {
		dspgGraph.RemoveEdge(node, node)
	}

	return dspgGraph
}

func CreateGraphFromMatrix(adjMatrix [][]int) (graph.Graph[string, string], error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())

	numVertices := len(adjMatrix)
	for i := 1; i <= numVertices; i++ {
		_ = g.AddVertex(strconv.Itoa(i))
	}

	for i := 0; i < numVertices; i++ {
		for j := 0; j < numVertices; j++ {
			if adjMatrix[i][j] == 1 {
				_ = g.AddEdge(strconv.Itoa(i+1), strconv.Itoa(j+1))
			}
		}
	}

	return g, nil
}

func DrawGraph(g graph.Graph[string, string], name string) {
	graph, err := os.Create("./png/" + name + ".gv")
	utils.Must(err)
	err = draw.DOT(g, graph)
	utils.Must(err)

	exec.Command("dot", "-Tsvg", "-O", "./png/"+name+".gv").Output()
}
