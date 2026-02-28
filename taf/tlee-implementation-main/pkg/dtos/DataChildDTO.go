package dtos

import (
	"fmt"
)

type TrustType int

// Define the enum values using iota
const (
	ReferralTrust TrustType = iota
	FunctionalTrust
)

type DataChildDTO struct {
	Data  KeyDTO         `json:"data,omitempty"`
	Child []DataChildDTO `json:"child,omitempty"`
	Tag   TrustType
}

/*
NEW
*/
func (node DataChildDTO) GenerateDotD(nodeID int) (string, int) {
	// Label for the current node
	var nodeLabel string
	if node.Data.Operation != "" {
		nodeLabel = node.Data.Operation
	} else {
		nodeLabel = fmt.Sprintf("%s -> %s",
			node.Data.FromNode, node.Data.ToNode)
	}
	// nodeLabel := fmt.Sprintf("%s\\n%s -> %s", node.Data.Operation, node.Data.FromNode, node.Data.ToNode)
	// Create the current node's DOT representation
	dot := fmt.Sprintf("    %d [label=\"%s\"];\n", nodeID, nodeLabel)

	// Initial ID for children of this node
	nextID := nodeID + 1
	var childID int

	// Iterate over child nodes, generating DOT code for each
	for _, child := range node.Child {
		// Connect the current node to its child
		dot += fmt.Sprintf("    %d -> %d;\n", nodeID, nextID)
		// Generate DOT for the child, including all its descendants
		var childDot string
		childDot, childID = child.GenerateDotD(nextID)
		dot += childDot
		// Update nextID based on the IDs used by this child and its descendants
		nextID = childID + 1
	}

	return dot, nextID - 1 // Return the DOT code and the last used ID
}

/*
END NEW
*/
