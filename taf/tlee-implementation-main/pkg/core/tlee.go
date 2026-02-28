package core

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/config"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dspgFunc"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/file"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/findSubTM"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/generateExpression"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/operation"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/vs-uulm/go-subjectivelogic/pkg/subjectivelogic"
	"github.com/vs-uulm/taf-tlee-interface/pkg/trustmodelstructure"
	"log/slog"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type TLEE struct {
	//map the trust model instance to a mapping with maps a fingerprint to an expression already synthesized
	FingerprintToExpr map[string]map[string]dtos.DataChildDTO
	VersionToValue    map[int][]trustmodelstructure.TrustRelationship
}

/*
The public function is responsible for generating and simplifying an expression based on a given list of dtos.VertexEdgeDTO representing vertex edges.
It produces a dtos.DataChildExpression that represents the simplified expression.

The function begins by creating a new directed acyclic graph dspgGraph using the graph package, where each vertex corresponds to a node and edges represent links between nodes.
It iterates through the provided vertexEdges and adds vertices and edges to the dspgGraph using dspgFunc.AddVertexAndEdge.
The function generates all possible Partial Path Sets (PPS) and their nesting levels using generateExpression.FindPpsAndNestingLevel based on the dspgGraph.
Nesting levels are added to the PPS using generateExpression.AddNestingLevelToPps.
The minimum nesting level for each PPS is calculated using dspgFunc.MinNestingLevel, and these values are assigned to each PPS.
A map expression is created to store expressions for each edge in the dspgGraph. Each edge is represented by an dtos.Expression key with an empty dtos.DataChildExpression value.
The PPS are sorted based on their minimum nesting levels and the number of edges using generateExpression.SortPps.
The function initializes subgraphExpression to represent the reduced subgraph expression.
For each PPS, the function reduces the graph and updates the expression map using generateExpression.ReduceGraph. It also updates the list of PPS by removing the first PPS and adding the newly reduced PPS.
The final expression is synthesized using generateExpression.SynthesizeFinalPath based on the original dspgGraph and the expression map.
If the final expression is empty, the subgraphExpression is returned as the result.
The generated and simplified expression is returned as a dtos.DataChildExpression.
*/
func ExpressionSynthesizer(vertexEdges []dtos.VertexEdgeDTO, checkDSPG bool, debug *bool) (bool, dtos.DataChildExpressionDTO) {
	dspgGraph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())

	for _, vertexEdge := range vertexEdges {
		dspgGraph = dspgFunc.AddVertexAndEdge(dspgGraph, vertexEdge.Node, vertexEdge.Links)
	}

	var stepCounter int64

	if debug != nil && *debug {
		dspgFunc.DebugDrawGraph(dspgGraph, "dspgGraph")
		stepCounter++
	}

	var result bool = true

	var allPPS []dtos.AllSubPPSDTO

	if checkDSPG {

		var contInDegree int
		var contOutDegree int

		var vertexInDegree []string
		var vertexOutDegree []string

		var wg sync.WaitGroup

		vertexList := dspgFunc.VertexList(dspgGraph)
		for _, vertex := range vertexList {
			wg.Add(2)

			go func() {
				defer wg.Done()
				vertexInDegree = dspgFunc.VertexInDegreeList(dspgGraph, vertex)
			}()

			go func() {
				defer wg.Done()
				vertexOutDegree = dspgFunc.VertexOutDegreeList(dspgGraph, vertex)
			}()

			wg.Wait()
			if len(vertexInDegree) == 0 {
				contInDegree++
			}
			if len(vertexOutDegree) == 0 {
				contOutDegree++
			}
		}

		var contEdgeLevelOne int = 0
		var contEdgeLevelTwo int = 0
		EdgeList := dspgFunc.VertexDegree(dspgGraph)
		for _, edge := range EdgeList {
			if edge == 1 {
				contEdgeLevelOne++
			} else if edge == 2 {
				contEdgeLevelTwo++
			}
		}

		if contEdgeLevelOne == 2 && (len(EdgeList)-2) == contEdgeLevelTwo {
			result = true
			return result, dtos.DataChildExpressionDTO{}
		}

		if contInDegree > 1 || contOutDegree > 1 {
			return false, dtos.DataChildExpressionDTO{}
		}

		dspgGraph, allPPS = dspgFunc.LselectParallelPathSubgraphs(dspgGraph)
		if len(allPPS) <= 0 {
			result = false
			return result, dtos.DataChildExpressionDTO{}
		}

		generateExpression.AddNestingLevelToPps(dspgGraph, allPPS)

		for i := range allPPS {
			allPPS[i].MinNestingLevel = dspgFunc.MinNestingLevel(allPPS[i].Graph)
		}

	} else {
		allPPS, dspgGraph = generateExpression.FindPpsAndNestingLevel(dspgGraph)

		generateExpression.AddNestingLevelToPps(dspgGraph, allPPS)

		for i := range allPPS {
			allPPS[i].MinNestingLevel = dspgFunc.MinNestingLevel(allPPS[i].Graph)
		}
	}

	expression := make(map[dtos.ExpressionDTO]dtos.DataChildExpressionDTO)
	allEdge, err := dspgGraph.Edges()
	utils.Must(err)

	for _, edge := range allEdge {
		expression[dtos.ExpressionDTO{FromNode: edge.Source, ToNode: edge.Target}] = dtos.DataChildExpressionDTO{}
	}

	generateExpression.SortPps(allPPS)

	var subgraphExpression dtos.DataChildExpressionDTO

	for _, pps := range allPPS {

		if checkDSPG {

			edges, _ := dspgGraph.Edges()
			var maxNestingLevel int = 0
			for _, edge := range edges {
				if edge.Properties.Weight > maxNestingLevel {
					maxNestingLevel = edge.Properties.Weight
				}
			}

			var selectPPSNestedLevelEQTo []dtos.AllSubPPSDTO
			for _, pps := range allPPS {
				if pps.MinNestingLevel == maxNestingLevel {
					selectPPSNestedLevelEQTo = append(selectPPSNestedLevelEQTo, pps)
				}
			}

			if len(selectPPSNestedLevelEQTo) == 0 {
				result = false
				return result, dtos.DataChildExpressionDTO{}
			}

			var allPpsGraph []graph.Graph[string, string]
			for _, pps := range selectPPSNestedLevelEQTo {

				allPpsGraph = append(allPpsGraph, pps.Graph)

			}

			if !dspgFunc.LocalDSPGTest(dspgGraph, allPpsGraph) {
				result = false
				return result, dtos.DataChildExpressionDTO{}
			}

		}

		if debug != nil && *debug {
			stepGraph, err := dspgGraph.Clone()
			utils.Must(err)
			dspgFunc.DebugDrawGraph(stepGraph, fmt.Sprintf("STEP_%v", stepCounter))

			dspgFunc.DebugDrawGraph(pps.Graph, fmt.Sprintf("PPS_%v", stepCounter))

			stepCounter++
		}

		subgraphExpression = generateExpression.ReduceGraph(dspgGraph, pps, expression)

		removedPps := allPPS[0]
		allPPS = generateExpression.UpdatePps(removedPps, allPPS)

		allPPS = allPPS[1:]
	}

	finalExpression := generateExpression.SynthesizeFinalPath(dspgGraph, expression)

	if (finalExpression.Data == dtos.ExpressionDTO{} && finalExpression.Child == nil) {
		finalExpression = subgraphExpression
	}

	if debug != nil && *debug {
		stepGraph, err := dspgGraph.Clone()
		utils.Must(err)
		dspgFunc.DebugDrawGraph(stepGraph, fmt.Sprintf("STEP_%v", stepCounter))
		stepCounter++
	}

	return result, finalExpression
}

/*
The public function is responsible for retrieving the real logical expression based on the provided dtos.DataChildExpression.
It transforms the logical operators in the expression according to the specified mathematicalModel and returns the updated expression.

The function checks if the expression has child expressions (sub-expressions). If it does, it recursively calls RetrieveRealExpression on each child expression to transform them accordingly.
It checks if the expression.Data.Operation is not an empty string (i.e., it's an operation). If it is, the function transforms the operation using a mapping defined in the dtos.LogicOperator based on the specified mathematicalModel.
The updated expression is returned with transformed logical operators.
*/
func MetaToConcreteExpressionConverter(expression dtos.DataChildExpressionDTO, mathematicalModel string) dtos.DataChildExpressionDTO {
	if len(expression.Child) > 0 {
		for i := range expression.Child {
			expression.Child[i] = MetaToConcreteExpressionConverter(expression.Child[i], mathematicalModel)
		}
	}

	if expression.Data.Operation != "" {
		expression.Data.Operation = operation.LogicOperator[mathematicalModel][expression.Data.Operation]
	}

	return expression
}

func ToDSPGTransform(vertexEdges []dtos.VertexEdgeDTO, debug *bool) dtos.SynthesizingGraph {
	graphIn := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())

	var wg sync.WaitGroup

	// Sort the vertexEdges based on the Node field
	sort.Slice(vertexEdges, func(i, j int) bool {
		return vertexEdges[i].Node < vertexEdges[j].Node
	})

	// Iterate over vertexEdges in deterministic order
	for _, vertexEdge := range vertexEdges {
		graphIn = dspgFunc.AddVertexAndEdge(graphIn, vertexEdge.Node, vertexEdge.Links)
	}

	if debug != nil && *debug {
		dspgFunc.DebugDrawGraph(graphIn, "Graph")
	}
	var source string
	var target string

	wg.Add(2)

	go func() {
		defer wg.Done()
		source = dspgFunc.GetSourceOriginalGraph(graphIn)
	}()

	go func() {
		defer wg.Done()
		target = dspgFunc.GetTargetOriginalGraph(graphIn)
	}()

	wg.Wait()

	var processingPath dtos.ProcessingPath
	var tempSynthesizingGraph []string
	var tempProcessingPath []string

	adj, _ := graphIn.AdjacencyMap()

	pathToProcess := dspgFunc.GetAllPaths(adj, source, target)

	utils.SortByLength(pathToProcess)

	pathToProcess, tempSynthesizingGraph = utils.SlicePop(pathToProcess, len(pathToProcess)-1)

	synthesizingGraph := dtos.SynthesizingGraph{Nodes: tempSynthesizingGraph, NodeToPPS: make(map[string]dtos.EdgeDTO)}
	synthesizingGraph.Edges = append(synthesizingGraph.Edges, dspgFunc.TransformPathToCouple(tempSynthesizingGraph)...)

	synthesizingGraph.Adj, _ = dspgFunc.CreateGraphFromPath(synthesizingGraph.Edges).AdjacencyMap()

	for len(pathToProcess) != 0 {
		if len(pathToProcess)%1000 == 0 {
			config.Logger.Debug("pathToProcess", "len", len(pathToProcess))
		}
		pathToProcess, tempProcessingPath = utils.SlicePop(pathToProcess, len(pathToProcess)-1)
		processingPath = dtos.ProcessingPath{Nodes: tempProcessingPath, NodeToPPS: make(map[string]dtos.EdgeDTO), Index: 0}
		processingPath.Edges = append(processingPath.Edges, dspgFunc.TransformPathToCouple(tempProcessingPath)...)

		var processingEdges []dtos.EdgeDTO
		endAddingPath := false

		for !endAddingPath {
			edge := processingPath.Edges[processingPath.Index]
			processingPath.Index += 1

			if processingPath.Index == len(processingPath.Edges) {
				endAddingPath = true
			}
			if !dspgFunc.ContainsNode(synthesizingGraph.Nodes, edge.ToNode) {
				processingEdges = append(processingEdges, edge)
			} else if len(processingEdges) != 0 || !dspgFunc.ContainsAllNode(synthesizingGraph.Edges, edge) {
				processingEdges = append(processingEdges, edge)
				nodeSrc := processingEdges[0].FromNode
				ultimateEdge := len(processingEdges) - 1
				nodeTarget := processingEdges[ultimateEdge].ToNode
				intermediateNodes := dspgFunc.DspgEdgeCheck(synthesizingGraph, nodeSrc, nodeTarget)
				if len(intermediateNodes) != 0 {
					for _, edge := range processingEdges {
						if !dspgFunc.ContainsNode(synthesizingGraph.Edges, edge) {
							synthesizingGraph.Edges = append(synthesizingGraph.Edges, edge)
						}
					}
					synthesizingGraph.Adj, _ = dspgFunc.CreateGraphFromPath(synthesizingGraph.Edges).AdjacencyMap()
					for _, node := range processingEdges {
						if !dspgFunc.ContainsNode(synthesizingGraph.Nodes, node.FromNode) {
							synthesizingGraph.Nodes = append(synthesizingGraph.Nodes, node.FromNode)
						}
						if !dspgFunc.ContainsNode(synthesizingGraph.Nodes, node.ToNode) {
							synthesizingGraph.Nodes = append(synthesizingGraph.Nodes, node.ToNode)
						}
					}
					for _, node := range intermediateNodes {
						if node != nodeSrc && node != nodeTarget {
							synthesizingGraph.NodeToPPS[node] = dtos.EdgeDTO{FromNode: nodeSrc, ToNode: nodeTarget}
						}
					}
					for _, node := range dspgFunc.GetNodes(processingEdges) {
						synthesizingGraph.NodeToPPS[node] = dtos.EdgeDTO{FromNode: nodeSrc, ToNode: nodeTarget}
					}
					processingEdges = []dtos.EdgeDTO{}
				} else {
					endAddingPath = true
				}
			}
		}
	}

	return synthesizingGraph
}

func CreateDAG(numNodes, numEdges int) [][]int {
	adjMatrix := make([][]int, numNodes)
	for i := range adjMatrix {
		adjMatrix[i] = make([]int, numNodes)
	}

	rand.Seed(time.Now().UnixNano())

	adjList := make([][]int, numNodes)

	topologicalOrder := rand.Perm(numNodes)

	for i := 0; i < numNodes; i++ {
		randomEdge := dspgFunc.GenerateRandomNumber(1, numEdges)
		for j := i + 1; j < numNodes && len(adjList[topologicalOrder[i]]) < randomEdge; j++ {
			if adjMatrix[topologicalOrder[i]][topologicalOrder[j]] == 0 {
				adjMatrix[topologicalOrder[i]][topologicalOrder[j]] = 1
				adjList[topologicalOrder[i]] = append(adjList[topologicalOrder[i]], topologicalOrder[j])
			}
		}
	}

	return adjMatrix
}

func GenerateRandomDSPGGraph(n int) graph.Graph[string, string] {
	Graph := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())
	mapGraph := make(map[string]int)

	source := "s"
	target := "t"

	Graph.AddVertex(source)
	Graph.AddVertex(target)
	Graph.AddEdge(source, target, graph.EdgeWeight(1.0))
	mapGraph[source+target] = 1

	for i := 0; i < n; i++ {
		links, _ := Graph.Edges()
		randomNumber := dspgFunc.GenerateRandomNumber(1, len(links))

		selectedEdge := links[randomNumber-1]
		coin := dspgFunc.GenerateRandomNumber(0, 1)

		if coin == 0 {
			node := strconv.Itoa(i)
			Graph.AddVertex(node)
			Graph.AddEdge(selectedEdge.Source, node)
			Graph.AddEdge(node, selectedEdge.Target)
			mapGraph[selectedEdge.Source+node] = 1
			mapGraph[node+selectedEdge.Target] = 1

			if mapGraph[selectedEdge.Source+selectedEdge.Target] == 1 {
				Graph.RemoveEdge(selectedEdge.Source, selectedEdge.Target)
				delete(mapGraph, selectedEdge.Source+selectedEdge.Target)
			} else {
				mapGraph[selectedEdge.Source+selectedEdge.Target] -= 1
			}

		} else {
			Graph.AddEdge(selectedEdge.Source, selectedEdge.Target, graph.EdgeWeight(1.0))
			mapGraph[selectedEdge.Source+selectedEdge.Target] += 1
		}
	}

	Graph = dspgFunc.RemoveSelfLoop(Graph)

	return Graph
}

/*
The EvalExpr function evaluates a logical expression represented by dtos.DataChildDTO.
It recursively processes the expression and its child expressions based on the specified getOpinionMode and opinionMap.

If the expression.Data.Operation is an empty string, it directly retrieves the opinion from the opinionMap and returns it.
Otherwise, it checks the number of child expressions (n) and handles cases where n is 0 (error), 1 (single child), or greater than 1 (multiple children).

For multiple children, it recursively evaluates each child expression, then applies the logical operation specified in expression.Data.Operation to combine their opinions.
The resulting opinion is returned along with the associated nodes.
*/
func EvalExpr(expression dtos.DataChildDTO, getOpinionMode string,
	opinionMap map[dtos.KeyDTO]dtos.OpinionDTOValue,
	operationFunc map[string]func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, target string,
	debug bool) (dtos.KeyDTO, error) {
	// Assigns the OperationFunc map from the operation package to the local variable opinionFunc.
	if expression.Data.Operation == "" {
		res := expression.Data
		op := dtos.KeyDTO{Opinion: opinionMap[dtos.KeyDTO{FromNode: res.FromNode, ToNode: res.ToNode}]}
		return op, nil
	} else {
		//child should be greater than 1
		n := len(expression.Child)
		if debug {
			config.Logger.Debug("", "expression.Data.Operation", expression.Data.Operation)
			config.Logger.Debug("", "n", n)
		}
		if n == 0 {
			// This should not happen
			//fmt.Println("hererfsfssrefffveffe")
			return dtos.KeyDTO{}, errors.New("operator without child")
		}
		if n == 1 {
			res := expression.Child[0].Data
			op := dtos.KeyDTO{Opinion: opinionMap[dtos.KeyDTO{FromNode: res.FromNode, ToNode: res.ToNode}]}
			return op, nil

		}
		if n > 1 {
			res := []dtos.KeyDTO{}
			// To ensure that the first element of res correspond to the first children
			tmp, _ := EvalExpr(expression.Child[0], getOpinionMode, opinionMap, operationFunc, target, debug)
			res = append(res, tmp)
			for i := 1; i < len(expression.Child)-1; i++ {
				child := expression.Child[i]
				tmp, _ = EvalExpr(child, getOpinionMode, opinionMap, operationFunc, target, debug)
				res = append(res, tmp)
			}
			// To ensure that the last element of res correspond to the last children
			tmp, _ = EvalExpr(expression.Child[len(expression.Child)-1], getOpinionMode, opinionMap, operationFunc, target, debug)
			res = append(res, tmp)

			if debug {
				config.Logger.Debug("input", "expression.Data.Operation", expression.Data.Operation)
			}
			op1 := res[0]
			op2 := res[1]
			if debug {
				config.Logger.Debug("", "res[0].FromNode", res[0].FromNode, "res[0].ToNode",
					res[0].ToNode, "op1", op1, "res[1].FromNode", res[1].FromNode, "res[1].ToNode", res[1].ToNode)
			}

			if expression.Data.Operation == "DISCOUNT" && expression.Tag == dtos.ReferralTrust {
				expression.Data.Operation = "DISCOUNT_REF"
			}
			opValue := operationFunc[expression.Data.Operation](op1.Opinion, op2.Opinion)
			if debug {
				config.Logger.Debug("output", "opValue", opValue)
			}
			for i := 2; i < n; i++ {
				if debug {
					config.Logger.Debug("input", "expression.Data.Operation", expression.Data.Operation, "opValue",
						opValue, "res[i].FromNode", res[i].FromNode, "res[i].ToNode", res[i].ToNode)
				}
				opi := res[i]
				if debug {
					config.Logger.Debug("opi", "value", opi)
				}

				if expression.Data.Operation == "DISCOUNT" && expression.Tag == dtos.ReferralTrust {
					expression.Data.Operation = "DISCOUNT_REF"
				}
				opValue = operationFunc[expression.Data.Operation](opValue, opi.Opinion)
				if debug {
					config.Logger.Debug("output", "opValue", opValue)
				}
			}
			return dtos.KeyDTO{Opinion: opValue,
				FromNode: expression.Child[0].Data.FromNode, ToNode: expression.Child[n-1].Data.ToNode}, nil
		}

	}
	return dtos.KeyDTO{}, nil
}

func ReferralDiscountChecker(expression dtos.DataChildDTO, target string,
	debug bool) (dtos.DataChildDTO, error) {
	// Assigns the OperationFunc map from the operation package to the local variable opinionFunc.
	if expression.Data.Operation == "" {
		if expression.Data.ToNode == target {
			expression.Tag = dtos.FunctionalTrust
		}
		return expression, nil

	} else {
		//child should be greater than 1
		n := len(expression.Child)
		if debug {
			config.Logger.Debug("", "expression.Data.Operation", expression.Data.Operation, "n", n)
		}
		if n == 0 {
			// This should not happen
			//fmt.Println("hererfsfssrefffveffe")
			return expression, errors.New("operator without child")
		}
		if n == 1 {
			if expression.Child[0].Tag == dtos.FunctionalTrust {
				expression.Tag = dtos.FunctionalTrust
			}
			return expression, nil

		}
		if n > 1 {
			var children []dtos.DataChildDTO
			for _, child := range expression.Child {
				child, _ = ReferralDiscountChecker(child, target, debug)
				children = append(children, child)
				if child.Tag == dtos.FunctionalTrust {
					expression.Tag = dtos.FunctionalTrust
				}
			}
			expression.Child = children
		}

	}
	return expression, nil
}

/*
The Evaluator function initializes the evaluation of a logical expression represented by dtos.DataChildDTO.
It retrieves opinions either from a provided slice or by calling the csv2Opinion function based on the getOpinionMode.

After ensuring the opinions are valid, it constructs a map (opinionMap) to associate each opinion with a corresponding key (FromNode, ToNode).
The function then calls EvalExpr to evaluate the logical expression and returns the result.

This function acts as a setup and entry point for evaluating complex logical expressions with opinion data.
*/
func Evaluator(expression dtos.DataChildDTO, getOpinionMode string,
	values []trustmodelstructure.TrustRelationship,
	operationFunc map[string]func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, target string) (dtos.KeyDTO, error) {
	var allOpinion []trustmodelstructure.TrustRelationship
	var err error
	if values == nil {
		allOpinion, err = csv2Opinion(getOpinionMode)
		utils.Must(err)

	} else {
		allOpinion = values
	}

	// utils.Must(evaluator.CheckOpinion(allOpinion))
	opinionMap := make(map[dtos.KeyDTO]dtos.OpinionDTOValue)

	for _, opinionDTO := range allOpinion {
		opinion := opinionDTO.Opinion()
		opinionMap[dtos.KeyDTO{FromNode: opinionDTO.Source(), ToNode: opinionDTO.Destination()}] = dtos.NewOpinionDTOValue(opinion.Belief(), opinion.Disbelief(), opinion.Uncertainty(), opinion.BaseRate())

	}
	return EvalExpr(expression, getOpinionMode, opinionMap, operationFunc, target, true)
}

func ExprToDTO(expression dtos.DataChildExpressionDTO) dtos.DataChildDTO {
	data := dtos.KeyDTO{Operation: expression.Data.Operation,
		FromNode: expression.Data.FromNode,
		ToNode:   expression.Data.ToNode}
	res := dtos.DataChildDTO{Data: data}
	l := len(expression.Child)
	child := make([]dtos.DataChildDTO, l)

	if l > 0 {
		for i := 0; i < l; i++ {
			child[i] = ExprToDTO(expression.Child[i])
		}
	}
	res.Child = child
	return res

}

/*
The public function is responsible for retrieving the real logical expression based on the provided dtos.DataChildExpression.
It transforms the logical operators in the expression according to the specified mathematicalModel and returns the updated expression.

The difference between this function and MetaToConcreteExpressionConverter is that MetaToConcreteExpressionConverter
returns a dtos.DataChildExpressionDT while MetaToConcreteExpressionConverterUpd returns a dtos.DataChildDTO
*/
func MetaToConcreteExpressionConverterUpd(expression dtos.DataChildExpressionDTO,
	mathematicalModel string) dtos.DataChildDTO {
	tmp := MetaToConcreteExpressionConverter(expression,
		mathematicalModel)
	tmp2 := ExprToDTO(tmp)
	return tmp2

}
func Init(logger *slog.Logger, filePath string, debuggingMode bool) {
	config.Logger = logger
	config.Logger.Info("Starting the initialization process.")
	config.OutputPath = filePath
	config.DebuggingMode = debuggingMode

}

func SpawnNewTLEE(logger *slog.Logger, filePath string, debuggingMode bool) *TLEE {
	Init(logger, filePath, debuggingMode)
	return &TLEE{}
}

/*
	A pointer receiver menthod that operates on the instance of the TLEE type.

The modifications made to the receiver inside the method will directly affect the original instance.
*/
func (tlee *TLEE) RunTLEE(trustmodelID string, version int, fingerprint uint32, structure trustmodelstructure.TrustGraphStructure,
	values map[string][]trustmodelstructure.TrustRelationship) (map[string]subjectivelogic.QueryableOpinion, error) {
	fusionOp, _ := operation.GetFusionOperator(trustmodelstructure.FusionOperator(structure.Operator()))
	discountOp, discountOpRef, _ := operation.GetDiscountOperator(trustmodelstructure.DiscountOperator(structure.DiscountOperator()))
	operationFunc := map[string]func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue{
		"FUSION":       fusionOp,
		"DISCOUNT":     discountOp,
		"DISCOUNT_REF": discountOpRef,
	}
	//Case we get the functional belief as input
	config.Logger.Info("TLEE has been called")
	config.Logger.Debug("Structure to be processed", "value", structure)
	subTrustModelList, agent := findSubTM.FindSubTrustModels(structure)
	config.Logger.Debug("Agent of the TM", "value", agent)
	res := make(map[string]subjectivelogic.QueryableOpinion)
	for _, trustModelPP := range subTrustModelList {
		debugBool := &config.DebuggingMode
		config.Logger.Debug("Current Sub Trust Model", "value", trustModelPP)
		config.Logger.Debug("Current Sub Trust Model Value", "value", values[trustModelPP.Scope])
		if len(trustModelPP.AdjacencyList) == 1 {
			for _, opinionDTO := range values[trustModelPP.Scope] {
				if opinionDTO.Source() == agent && opinionDTO.Destination() == trustModelPP.Scope {
					res[trustModelPP.Scope] = opinionDTO.Opinion()
				}
			}

		} else {
			op, _ := tlee.runTLEEperSubTM(trustmodelID, version, fingerprint, trustModelPP, values[trustModelPP.Scope],
				operationFunc, agent, trustModelPP.Scope, debugBool)
			res[trustModelPP.Scope] = op.Opinion
		}

	}
	config.Logger.Info("TLEE has finished the calculation")
	return res, nil
}

func (tlee *TLEE) runTLEEperSubTM(tmID string, version int, fingerprint uint32, structure dtos.StructureGraphTAFSingleProp,
	values []trustmodelstructure.TrustRelationship,
	operationFunc map[string]func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, agent string, scope string, debugBool *bool) (*dtos.KeyDTO, bool) {
	var finalOpinion *dtos.KeyDTO
	finalOpinion = nil

	var concreteExpression dtos.DataChildDTO
	ok2 := false
	/*
		performing a map lookup and an assertion.
		The statement performs a lookup in the `FingerprintToExpr` map in the `tlee` instance, and it retrieves the value associated with the key `fingerprint`.
		The found value is associated to the `concreteExpression` variable.
		`ok` is assigned a boolean value indicating whether the key `fingerprint` was present in the map (`true` if the `fingerprint` key exists in the map `FingerprintToExpr`, `false` otherwise)
	*/
	fingerprintstr := strconv.FormatUint(uint64(fingerprint), 10)
	fingerprintstr = fingerprintstr + " - " + scope
	concreteExpressionMap, ok := tlee.FingerprintToExpr[tmID]
	if ok {
		concreteExpression, ok2 = concreteExpressionMap[fingerprintstr]
	}
	if !ok2 {

		/* running ToDSPGTransform */
		config.Logger.Info("Starting DSPG Transformer")
		synthesizedGraph := ToDSPGTransform(structure.AdjacencyList, debugBool)

		/* running ExpressionSynthesizer */
		config.Logger.Info("Starting Expression Synthesizer")
		var adjsynthesizedGraph []dtos.VertexEdgeDTO
		adjMat := synthesizedGraph.Adj
		for _, value := range synthesizedGraph.Nodes {
			var lEdge []string
			for _, edge := range adjMat[value] {
				lEdge = append(lEdge, edge.Target)
			}

			adjsynthesizedGraph = append(adjsynthesizedGraph, dtos.VertexEdgeDTO{Node: value, Links: lEdge})
		}
		_, metaExpression := ExpressionSynthesizer(adjsynthesizedGraph, false, debugBool)
		if *debugBool == true {
			utils.CreateExpDot(nil, &metaExpression, "meta")
		}

		/* running MetaToConcreteExpressionConverterUpd */
		config.Logger.Info(" Starting MetaToConcrete")
		concreteExpression = MetaToConcreteExpressionConverterUpd(metaExpression, "subjectiveLogic") // Outputs a dtos.DataChildDTO
		if *debugBool == true {
			utils.CreateExpDot(&concreteExpression, nil, "concrete")

		}
	}

	/* Running ReferralDiscountChecker*/
	concreteExpression, _ = ReferralDiscountChecker(concreteExpression, scope, *debugBool)
	/* running Evaluator */
	config.Logger.Info("Starting Evaluator")
	tmp, _ := Evaluator(concreteExpression, "", values, operationFunc, scope)
	finalOpinion = &tmp
	return finalOpinion, true

}

// Check if direct opinion from source to target and returns
func checkDirectRel(values []trustmodelstructure.TrustRelationship, agent string, scope string) (bool, trustmodelstructure.TrustRelationship) {

	for _, val := range values {
		if val.Source() == agent && val.Destination() == scope {
			return true, val
		}
	}
	return false, nil
}

func json2vTe(graphJsonFileName string) ([]dtos.VertexEdgeDTO, error) {
	fileContent, err := os.ReadFile(graphJsonFileName)
	if err != nil {
		config.Logger.Error("json2vTe cannot read file", "value", err)

	}
	var vTe []dtos.VertexEdgeDTO
	err = json.Unmarshal([]byte(fileContent), &vTe)
	if err != nil {
		config.Logger.Error("json2vTe cannot parse file content", "value", err)
	}
	return vTe, nil
}

func csv2Opinion(opininionValueCSVFileName string) ([]trustmodelstructure.TrustRelationship, error) {
	val, err := file.GetOpinion(opininionValueCSVFileName)
	if err != nil {
		config.Logger.Error("Cannot parse Csv into Opinion", "value", err)
	}
	var res []trustmodelstructure.TrustRelationship
	for _, opinionVal := range val {
		res = append(res, dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(opinionVal.Belief, opinionVal.Disbelief, opinionVal.Uncertainty, opinionVal.BaseRate),
			FromNode: opinionVal.FromNode,
			ToNode:   opinionVal.ToNode,
		})
	}
	return res, nil
}

func calculateProbability(finalOp map[string]dtos.OpinionDTOValue) map[string]float64 {
	probMapPerScope := make(map[string]float64)

	// Iterate over the map and take the values
	for key, value := range finalOp {
		probMapPerScope[key] = value.ProjectedProbability()
	}

	return probMapPerScope
}

/* --------------------- run TLEE ----------------------- */

// This line ensures that VertexEdgeDTO implements AdjacencyListEntry
var _ trustmodelstructure.AdjacencyListEntry = (*dtos.VertexEdgeDTO)(nil)
