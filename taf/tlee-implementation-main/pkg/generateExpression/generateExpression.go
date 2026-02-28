package generateExpression

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/operation"
	"sort"
	"sync"

	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dspgFunc"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/utils"

	"github.com/dominikbraun/graph"
)

/*
The public function takes in two parameters:
dspgGraph represented as a graph with string nodes and string edges, and expression represented as a map of Expression to DataChildExpression.
It synthesizes a final data structure DataChildExpression based on the input data and expressions.
The function starts by obtaining the source and target nodes from the original dspgGraph using the GetSourceOriginalGraph and GetTargetOriginalGraph functions.
It calculates all paths between the source and target nodes in dspgGraph and stores them in the path variable.
It checks if the length of the first path in path is greater than 2, indicating a valid path:

If it's a valid path, the function calls the PathToExpression function to construct a final expression based on the path and expressions.
It removes all edges from dspgGraph.
It adds an edge from the first two nodes of the path to dspgGraph.
It identifies and removes isolated nodes from dspgGraph using the FindIsolatedVertex function.
The function returns the finalExpression.
If the path length is not greater than 2, it returns an empty DataChildExpression.
*/
func SynthesizeFinalPath(dspgGraph graph.Graph[string, string], expression map[dtos.ExpressionDTO]dtos.DataChildExpressionDTO) dtos.DataChildExpressionDTO {

	var wg sync.WaitGroup

	var source string
	var target string

	wg.Add(2)

	go func() {
		defer wg.Done()
		source = dspgFunc.GetSourceOriginalGraph(dspgGraph)
	}()

	go func() {
		defer wg.Done()
		target = dspgFunc.GetTargetOriginalGraph(dspgGraph)
	}()

	wg.Wait()

	adj, err := dspgGraph.AdjacencyMap()
	utils.Must(err)
	path := dspgFunc.GetAllPaths(adj, source, target)
	if len(path[0]) > 2 {
		finalExpression := PathToExpression(dspgGraph, expression, path)

		linksDspgGraph, err := dspgGraph.Edges()
		utils.Must(err)

		for _, link := range linksDspgGraph {
			dspgGraph.RemoveEdge(link.Source, link.Target)
		}

		dspgGraph.AddEdge(path[0][0], path[0][1])

		arrayIsolateNode := dspgFunc.FindIsolatedVertex(dspgGraph)
		for _, node := range arrayIsolateNode {
			dspgGraph.RemoveVertex(node)
		}

		return finalExpression
	} else {
		return dtos.DataChildExpressionDTO{}
	}
}

/*
The public function takes in three parameters:
dspgGraph represented as a graph with string nodes and string edges, expression represented as a map of Expression to DataChildExpression,
and path represented as a slice of slices of strings.
It constructs a data structure DataChildExpression based on the input data.
The function starts by initializing finalPath using the first path in the path slice.
It initializes finalExpression with an operation type AGENT_OPERATOR_TRUST_FILTERING.
It initializes tempTree and tempList to empty values.
It iterates through each path in finalPath:

It checks if there is an existing expression for the path in the expression map.
If an expression exists, it populates tempTree and appends its children to tempList.
If no expression exists, it creates a new expression and appends it to tempList.

It sets the Child field of finalExpression to tempList.
The function returns finalExpression, which represents the constructed data structure.
*/
func PathToExpression(dspgGraph graph.Graph[string, string], expression map[dtos.ExpressionDTO]dtos.DataChildExpressionDTO, path [][]string) dtos.DataChildExpressionDTO {

	finalPath := utils.ListToTuples(path[0])

	var finalExpression dtos.DataChildExpressionDTO
	finalExpression.Data.Operation = operation.AGENT_OPERATOR_TRUST_FILTERING
	var tempTree dtos.DataChildExpressionDTO
	var tempList []dtos.DataChildExpressionDTO

	for _, path := range finalPath {

		check := expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Data != dtos.ExpressionDTO{}
		if check {
			tempTree.Data = expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Data
			tempTree.Child = append(tempTree.Child, expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Child...)
			tempList = append(tempList, tempTree)
		} else {
			expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}] = dtos.DataChildExpressionDTO{Data: dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}}
			tempList = append(tempList, expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}])
		}
	}

	finalExpression.Child = tempList

	return finalExpression
}

/*
The private function takes in three parameters:
dspgGraph represented as a graph with string nodes and string edges,
pps represented as an object of type AllSubPPS, and expression represented as a map of Expression to DataChildExpression.
It constructs a data structure DataChildExpression based on the input data and expressions.
The function starts by initializing head and discount as DataChildExpression objects.
It sets the Operation of head to AGENT_OPERATOR_TRUST_MERGE.
It initializes tempTree and tempList to empty values.
It calculates all paths between pps.Source and pps.Target in pps.Graph and stores them in AllPaths.
It iterates through each path in AllPaths and performs the following operations:

It converts the path into tuples using utils.ListToTuples and stores it in paths.
It iterates through each tuple in paths:

It checks if there is an existing expression for the tuple in the expression map.
If an expression exists, it populates tempTree and appends its children to tempList.
If no expression exists, it creates a new expression and appends it to tempList.

It sets the Operation of discount to AGENT_OPERATOR_TRUST_FILTERING and assigns tempList as its children.
It appends discount to the Child field of head.

The function returns head, which represents the constructed data structure.
*/
func subgraphWithExpression(dspgGraph graph.Graph[string, string], pps dtos.AllSubPPSDTO, expression map[dtos.ExpressionDTO]dtos.DataChildExpressionDTO) dtos.DataChildExpressionDTO {
	var head dtos.DataChildExpressionDTO
	var discount dtos.DataChildExpressionDTO

	head.Data.Operation = operation.AGENT_OPERATOR_TRUST_MERGE

	var tempTree dtos.DataChildExpressionDTO

	var tempList []dtos.DataChildExpressionDTO

	adj, err := pps.Graph.AdjacencyMap()
	utils.Must(err)
	AllPaths := dspgFunc.GetAllPaths(adj, pps.Source, pps.Target)

	for _, pathPps := range AllPaths {
		paths := utils.ListToTuples(pathPps)
		tempList = []dtos.DataChildExpressionDTO{}
		for _, path := range paths {
			check := expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Data != dtos.ExpressionDTO{}
			if check {
				tempTree.Data = expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Data
				tempTree.Child = expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}].Child
				tempList = append(tempList, tempTree)
			} else {
				expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}] = dtos.DataChildExpressionDTO{Data: dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}}
				tempList = append(tempList, expression[dtos.ExpressionDTO{FromNode: path[0], ToNode: path[1]}])
			}
		}

		discount.Data.Operation = operation.AGENT_OPERATOR_TRUST_FILTERING
		discount.Child = tempList

		head.Child = append(head.Child, discount)

	}

	return head
}

/*
The public function takes in two parameters: removedPps represented as an object of type AllSubPPS,
and allSortedPps represented as a slice of AllSubPPS.
It updates the allSortedPps based on the removal of a specific AllSubPPS object.
The function starts by obtaining the edges of the removedPps.Graph and stores them in the linksRemovedPps variable.
It initializes an empty slice ppsListUpdated to store the updated subgraphs.
It iterates through each pps object in allSortedPps:

It checks if the pps object has nodes corresponding to the removedPps.Source and removedPps.Target.
If the nodes are present, it removes the edges and isolated nodes from the pps.Graph based on the linksRemovedPps.
It adds the removedPps.Source and removedPps.Target as vertices and creates an edge between them in the pps.Graph.
The pps object is added to the ppsListUpdated.

The function returns ppsListUpdated, which represents the updated list of sorted subgraphs.
*/
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

/*
The public function takes in two parameters:
graph of type AllSubPPS and node of type string.
It checks whether a specified node exists in the graph.Graph represented by the graph parameter.

The function starts by obtaining the adjacency map of graph.Graph and stores it in the adj variable.
It iterates through the keys in the adj map, which represent nodes in the graph.
For each node in the map, it checks if it matches the specified node.
If a match is found, the function returns true to indicate that the node exists in the graph.
If no match is found after iterating through all nodes, the function returns false to indicate that the node does not exist in the graph.
*/
func HasNode(graph dtos.AllSubPPSDTO, node string) bool {
	adj, err := graph.Graph.AdjacencyMap()
	utils.Must(err)
	for nodeList := range adj {
		if nodeList == node {
			return true
		}
	}
	return false
}

/*
This private function takes in two parameters: dspgGraph and allPPS. It operates on a graph
represented by dspgGraph and a list of AllSubPPS objects represented by allPPS.
This function propagates edge weights from the dspgGraph to the pps.Graph for all combinations of edges and pps objects,
effectively updating the weights in the pps.Graph based on the data from the dspgGraph.

It starts by obtaining the edges of the dspgGraph and stores them in the links variable.
It then enters a nested loop structure:

The outer loop iterates through each link in the links slice, representing edges in the dspgGraph.
The inner loop iterates through each pps object in the allPPS list.

For each combination of link and pps, it updates the weight of the edge in the pps.Graph by calling
pps.Graph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(link.Properties.Weight)).
This essentially transfers the weight information from the dspgGraph to the corresponding edge in pps.Graph.
*/
func AddNestingLevelToPps(dspgGraph graph.Graph[string, string], allPPS []dtos.AllSubPPSDTO) {
	links, err := dspgGraph.Edges()
	utils.Must(err)
	for _, link := range links {
		for _, pps := range allPPS {
			pps.Graph.UpdateEdge(link.Source, link.Target, graph.EdgeWeight(link.Properties.Weight))
		}
	}
}

/*
The public function takes in three parameters:
dspgGraph represented as a graph with string nodes and string edges, pps represented as an object of type AllSubPPS,
and expression represented as a map of Expression to DataChildExpression.
It performs operations to reduce the dspgGraph and update the expression map.

The function starts by creating subgraphWithExpression by calling the subgraphWithExpression function with dspgGraph, pps, and expression as parameters.
It retrieves the edges of pps.Graph and stores them in the edges variable.
It iterates through each edge in edges and removes the corresponding edge from dspgGraph.
It adds an edge from pps.Source to pps.Target in dspgGraph.
It updates the expression map with the subgraphWithExpression using the pps.Source and pps.Target as the key.
It identifies and removes isolated nodes from dspgGraph using the FindIsolatedVertex function.
The function returns subgraphWithExpression, which represents the reduced subgraph with expressions.
*/

func ReduceGraph(dspgGraph graph.Graph[string, string], pps dtos.AllSubPPSDTO, expression map[dtos.ExpressionDTO]dtos.DataChildExpressionDTO) dtos.DataChildExpressionDTO {
	subgraphWithExpression := subgraphWithExpression(dspgGraph, pps, expression)

	edges, err := pps.Graph.Edges()
	utils.Must(err)

	for _, edge := range edges {
		dspgGraph.RemoveEdge(edge.Source, edge.Target)
	}

	dspgGraph.AddEdge(pps.Source, pps.Target)

	expression[dtos.ExpressionDTO{FromNode: pps.Source, ToNode: pps.Target}] = subgraphWithExpression

	arrayIsolateNode := dspgFunc.FindIsolatedVertex(dspgGraph)
	for _, node := range arrayIsolateNode {
		dspgGraph.RemoveVertex(node)
	}

	return subgraphWithExpression
}

/*
The public function takes in a slice of AllSubPPS objects represented as allPps and performs custom sorting based on specified criteria.

The function defines a custom sorting function customSort that compares two AllSubPPS objects i and j based on the following criteria:

If the MinNestingLevel of i is greater than the MinNestingLevel of j, it returns true.
If the MinNestingLevel of i is less than the MinNestingLevel of j, it returns false.
If the MinNestingLevel of i is equal to the MinNestingLevel of j, it further compares the number of edges in their respective graphs. It returns true if the number of edges in the graph of i is less than the number of edges in the graph of j, and false otherwise.

The sort.SliceStable function is called on the allPps slice, and the custom sorting function customSort is applied for stable sorting.
After sorting, the allPps slice is updated with the sorted order based on the specified criteria.
*/
func SortPps(allPps []dtos.AllSubPPSDTO) {

	customSort := func(i, j int) bool {
		if allPps[i].MinNestingLevel > allPps[j].MinNestingLevel {
			return true
		} else if allPps[i].MinNestingLevel < allPps[j].MinNestingLevel {
			return false
		}

		a, err := allPps[i].Graph.Edges()
		utils.Must(err)
		b, err := allPps[j].Graph.Edges()
		utils.Must(err)

		return len(a) < len(b)
	}

	sort.SliceStable(allPps, customSort)
}

/*
The public function receives a dspgGraph represented as a graph with nodes and string edges.
It performs a series of operations to find all possible paths between sources and
destinations in the dspgGraph and computes the nesting levels for each valid path.

The function starts by initializing sync.WaitGroup which allows us to use go routines to make the code parallelizable.
We also initialize adjacencyMap of the graph and 2 maps sourceSet and targetSet.
We initialize the counter of the go routine wg.Add(2) and 2 functions SourceSet and TargetSet are executed which will take all PPS paths from the source and target. In addition the wg.Wait() will wait for the wg.Done() to be done inside the functions which will communicate the end of the go routines.
Next, the two for will scroll through all the source and target nodes and check the PPS within a function initialized as a go routine that will allow parallelism.
Within the function the intersection of each path from the source and each path from the target will be done to find the common ones.
A graph will be created thanks to the CreateGraphFromPath function and it will be checked which intersection generated a result.
In the positive case the checkPPS (which possible to find in documentation) will be performed.
If this checkPPS was successful an array will be populated that will count the graph from the relative intersection, source and target of the PPS.
Finally given the PPS the nestingLevel will be added to the main graph.
Once this fuction and related for will be finished the return operation will be performed which will return array of PPS and the updated graph with related nestingLevel.
*/
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
