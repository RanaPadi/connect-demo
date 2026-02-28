package main

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/core"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dspgFunc"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"

	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/vs-uulm/taf-tlee-interface/pkg/trustmodelstructure"
	"log/slog"
	"os"
	"strconv"
)

/* This line ensures that VertexEdgeDTO implements AdjacencyListEntry */
var _ trustmodelstructure.AdjacencyListEntry = (*dtos.VertexEdgeDTO)(nil)

/* calculateProbability calculates the Projected probability for an opinion. */
func calculateProbability(finalOp map[string]dtos.OpinionDTOValue) map[string]float64 {
	probMapPerScope := make(map[string]float64)

	// Iterate over the map and take the values
	for key, value := range finalOp {
		probMapPerScope[key] = value.ProjectedProbability()
	}
	return probMapPerScope
}

/* ----- Call the TLEE on different Trust Models with the following functions:  testTleeOneEdge(), testTlee(), testTlee_nonDSPG() and testTleeDebug() ----- */

func testTleeOneEdge() {
	my_graph := dtos.StructureGraphTAFMultiplePropConstr(
		trustmodelstructure.CumulativeFusion,
		trustmodelstructure.DefaultDiscount,
		[]trustmodelstructure.AdjacencyListEntry{

			//dtos.VertexEdgeDTO{Node: "V_RE", Links: []string{"V_PU"}},
			//dtos.VertexEdgeDTO{Node: "V_PU", Links: []string{"CAM1"}},
			dtos.VertexEdgeDTO{Node: "V_RE", Links: []string{"V_CE", "V_YE", "CAM2"}},
			dtos.VertexEdgeDTO{Node: "V_YE", Links: []string{"CAM1"}},
			dtos.VertexEdgeDTO{Node: "V_CE", Links: []string{"CAM1"}}, // => good
			dtos.VertexEdgeDTO{Node: "CAM1", Links: []string{"CAM2"}}, // => good
		},
	)
	fmt.Println("StructureGraphTAFMultipleProp: ", my_graph)

	/* Assigning the numerical values for trust opinions */

	/*	The values of the numerical values of the trust opinions (TOs) are stored in a a map data structure, with key type string, and the value type is a slice of dtos.OpinionDTO.
		The TOs are calculated per trust relationship.
		The key holds the concrete proposition/scope for which this TO is calculated. This helps us to differentiate cases when two trust relationships have the same trustor and trustee, but a different scope. */

	numValsPerTR := make(map[string][]trustmodelstructure.TrustRelationship)

	numValsPerTR["CAM2"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_CE",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_YE",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "CAM2",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_CE",
			ToNode:   "CAM1",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_YE",
			ToNode:   "CAM1",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "CAM1",
			ToNode:   "CAM2",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_PU",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "V_RE",
			ToNode:   "CAM2",
		},
	}

	fmt.Println("Numerical Trust Opinions: ", numValsPerTR)

	// Creates a new instance of the TLEE type. The & operator is used to get the address of a newly created TLEE object, essentially creating a pointer to TLEE.
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a logger that writes to the log file
	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	tlee := core.SpawnNewTLEE(logger, "./outp/", true)

	// Calls the RunTLEE method on the prevously created tlee instance. */
	finalOp, _ := tlee.RunTLEE("test", 0, 0, my_graph, numValsPerTR)
	fmt.Println("Final Opinion:", finalOp)

}
func testTlee() {

	my_graph := dtos.StructureGraphTAFMultiplePropConstr(
		trustmodelstructure.CumulativeFusion,
		trustmodelstructure.OppositeBeliefDiscount,
		[]trustmodelstructure.AdjacencyListEntry{

			dtos.VertexEdgeDTO{Node: "V_RE", Links: []string{"V_YE", "V_TEST"}},
			dtos.VertexEdgeDTO{Node: "V_YE", Links: []string{"CAM1"}},
			dtos.VertexEdgeDTO{Node: "V_TEST", Links: []string{"CAM1"}},
			//dtos.VertexEdgeDTO{Node: "V_CE", Links: []string{"CAM1"}}, // => good
			dtos.VertexEdgeDTO{Node: "CAM1", Links: []string{"CAM2"}}, // => good
		},
	)
	fmt.Println("StructureGraphTAFMultipleProp: ", my_graph)

	/* Assigning the numerical values for trust opinions */

	/*	The values of the numerical values of the trust opinions (TOs) are stored in a map data structure, with key type string, and the value type is a slice of dtos.OpinionDTO.
		The TOs are calculated per trust relationship.
		The key holds the concrete proposition/scope for which this TO is calculated. This helps us to differentiate cases when two trust relationships have the same trustor and trustee, but a different scope. */

	numValsPerTR := make(map[string][]trustmodelstructure.TrustRelationship)

	numValsPerTR["CAM2"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_CE",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_YE",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "CAM2",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_CE",
			ToNode:   "CAM1",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "V_YE",
			ToNode:   "CAM1",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "CAM1",
			ToNode:   "CAM2",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_PU",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "V_RE",
			ToNode:   "CAM2",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "V_RE",
			ToNode:   "V_TEST",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "V_TEST",
			ToNode:   "CAM1",
		},
	}

	fmt.Println("Numerical Trust Opinions: ", numValsPerTR)

	// Creates a new instance of the TLEE type. The & operator is used to get the address of a newly created TLEE object, essentially creating a pointer to TLEE.
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a logger that writes to the log file
	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	tlee := core.SpawnNewTLEE(logger, "./outp/", true)

	// Calls the RunTLEE method on the prevously created tlee instance.
	finalOp, _ := tlee.RunTLEE("test", 0, 0, my_graph, numValsPerTR)
	fmt.Println("Final Opinion:", finalOp)

}

/* Creating and testing a nonDSPG graph */
func testTlee_nonDSPG() {

	my_graph := dtos.StructureGraphTAFMultiplePropConstr(
		trustmodelstructure.CumulativeFusion,
		trustmodelstructure.DefaultDiscount,
		[]trustmodelstructure.AdjacencyListEntry{
			dtos.VertexEdgeDTO{Node: "Ver", Links: []string{"A", "B"}},
			dtos.VertexEdgeDTO{Node: "A", Links: []string{"E"}},
			dtos.VertexEdgeDTO{Node: "B", Links: []string{"D", "C"}},
			dtos.VertexEdgeDTO{Node: "D", Links: []string{"E", "User"}}, // => good
			dtos.VertexEdgeDTO{Node: "C", Links: []string{"User"}},
			dtos.VertexEdgeDTO{Node: "E", Links: []string{"User"}},
		},
	)

	fmt.Println("StructureGraphTAFMultipleProp: ", my_graph)

	numValsPerTR := make(map[string][]trustmodelstructure.TrustRelationship)

	numValsPerTR["User"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "Ver",
			ToNode:   "A",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "Ver",
			ToNode:   "B",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "A",
			ToNode:   "E",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "B",
			ToNode:   "D",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "B",
			ToNode:   "C",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "D",
			ToNode:   "E",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "E",
			ToNode:   "User",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "D",
			ToNode:   "User",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.7, 0.1, 0.2, 0.5),
			FromNode: "C",
			ToNode:   "User",
		},
	}

	fmt.Println("Numerical Trust Opinions: ", numValsPerTR)

	// Creates a new instance of the TLEE type. The & operator is used to get the address of a newly created TLEE object, essentially creating a pointer to TLEE.
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a logger that writes to the log file
	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	tlee := core.SpawnNewTLEE(logger, "./outp/", true)

	// Calls the RunTLEE method on the prevously created tlee instance.
	finalOp, _ := tlee.RunTLEE("test", 0, 0, my_graph, numValsPerTR)
	fmt.Println("Final Opinion:", finalOp)

}

/*
This function creates and visualizes a Directed Acyclic Graph (DAG) with `nnodes` nodes and `nedges` edges, then generates a debug visualization of the graph.
The generated graph is stored in /outp/debug/
*/
func testGenerator() {

	nnodes, nedges := 10, 20
	adjMat := core.CreateDAG(nnodes, nedges)
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic(), graph.Weighted())

	// Add nodes to the graph
	n := len(adjMat)
	for i := 0; i < n; i++ {
		_ = g.AddVertex(strconv.Itoa(i))
	}

	// Add edges based on the adjacency matrix
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adjMat[i][j] == 1 {
				_ = g.AddEdge(strconv.Itoa(i), strconv.Itoa(j))
			}
		}
	}

	// Exporting the graph
	dspgFunc.DebugDrawGraph(g, "NDSPG10.20")
}

/* testTleeDebug() for testing the TLEE with specifying different Discounting Operators */

func testTleeDebug() {

	my_graph := dtos.StructureGraphTAFMultiplePropConstr(
		trustmodelstructure.CumulativeFusion,
		trustmodelstructure.DefaultDiscount,
		[]trustmodelstructure.AdjacencyListEntry{
			dtos.VertexEdgeDTO{Node: "M_TCG", Links: []string{"V_x", "C_x", "C_y", "C_z"}},
			dtos.VertexEdgeDTO{Node: "V_x", Links: []string{"C_x", "C_y", "C_z"}},
		},
	)

	fmt.Println("StructureGraphTAFMultipleProp: ", my_graph)

	/* Numerical values for trust opinions */

	/*	The values of the numerical values of the trust opinions (TOs) are stored in a a map data structure, with key type string, and the value type is a slice of dtos.OpinionDTO.
		The TOs are calculated per trust relationship.
		The key holds the concrete proposition/scope for which this TO is calculated. This helps us to differentiate cases when two trust relationships have the same trustor and trustee, but a different scope. */

	numValsPerTR := make(map[string][]trustmodelstructure.TrustRelationship)

	numValsPerTR["C_x"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "V_x",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_x",
			ToNode:   "C_x",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "C_x",
		},
	}

	numValsPerTR["C_y"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "V_x",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "C_y",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_x",
			ToNode:   "C_y",
		},
	}

	numValsPerTR["C_z"] = []trustmodelstructure.TrustRelationship{
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(1, 0, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "V_x",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "M_TCG",
			ToNode:   "C_z",
		},
		dtos.OpinionDTO{
			Value:    dtos.NewOpinionDTOValue(0.5, 0.5, 0, 0.5),
			FromNode: "V_x",
			ToNode:   "C_z",
		},
	}

	fmt.Println("Numerical Trust Opinions: ", numValsPerTR)

	/* Creates a new instance of the TLEE type. The & operator is used to get the address of a newly created TLEE object, essentially creating a pointer to TLEE. */
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a logger that writes to the log file
	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	tlee := core.SpawnNewTLEE(logger, "./outp/", true)
	/* Calls the RunTLEE method on the prevously created tlee instance. */

	// retrieves a fusion operator based on my_graph.Operator(), converts it to the expected type, and assigns it to fusionOp, while discarding the secondary return value.
	finalOp, _ := tlee.RunTLEE("test", 0, 0, my_graph, numValsPerTR)
	fmt.Println("Final Opinion:", finalOp)

}

func main() {
	//testTleeOneEdge()
	testTlee()
	//testTlee_nonDSPG()
	//testTleeDebug()
	//testGenerator()
}
