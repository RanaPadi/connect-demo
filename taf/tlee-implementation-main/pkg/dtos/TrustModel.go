package dtos

import (
	"github.com/vs-uulm/taf-tlee-interface/pkg/trustmodelstructure"
)

type ScopeVertexEdgeDTO struct {
	Node  string   `json:"node"`
	Links []string `json:"links"`
	Scope string   `json:"scope"`
}

type StructureGraphTAFSingleProp struct {
	AdjacencyList []VertexEdgeDTO //AdjacencyList is a slice of `VertexEdgeDTO` structs.
	Scope         string          // target variable
}

/*type StructureGraphTAFMultipleProp struct {
	Operator      string //Assuming same fusion op for all nodes
	AdjacencyList []VertexEdgeDTO
}*/

type StructureGraphTAFMultipleProp struct {
	fusionOperator   trustmodelstructure.FusionOperator //Assuming same fusion op for all nodes
	discountOperator trustmodelstructure.DiscountOperator
	adjacencyList    []trustmodelstructure.AdjacencyListEntry
}

// StructureGraphTAFMultiplePropConstr is a constructor function for StructureGraphTAFMultipleProp.
func StructureGraphTAFMultiplePropConstr(fusionOperator trustmodelstructure.FusionOperator, discountOperator trustmodelstructure.DiscountOperator, adjacencyList []trustmodelstructure.AdjacencyListEntry) StructureGraphTAFMultipleProp {
	return StructureGraphTAFMultipleProp{
		fusionOperator:   fusionOperator,
		discountOperator: discountOperator,
		adjacencyList:    adjacencyList,
	}
}

// Implements interface TrustGraphStructure from trustmodelstructure
// Operator returns the fusion operator of the StructureGraphTAFMultipleProp.
func (t StructureGraphTAFMultipleProp) Operator() trustmodelstructure.FusionOperator {
	return t.fusionOperator
}

// Implements interface TrustGraphStructure from trustmodelstructure
// DiscountOperator returns the discount operator of the StructureGraphTAFMultipleProp.
func (t StructureGraphTAFMultipleProp) DiscountOperator() trustmodelstructure.DiscountOperator {
	return t.discountOperator
}

// Implements interface TrustGraphStructure from trustmodelstructure
// AdjacencyList returns the adjacency list of the StructureGraphTAFMultipleProp.
// func (t StructureGraphTAFMultipleProp) AdjacencyList() []AdjacencyListEntry {
func (t StructureGraphTAFMultipleProp) AdjacencyList() []trustmodelstructure.AdjacencyListEntry {
	return t.adjacencyList
}
