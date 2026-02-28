package dtos

type VertexEdgeDTO struct {
	Node  string   `json:"node"`
	Links []string `json:"links"`
}

// Implements interface AdjacencyListEntry from trustmodelstructure
// SourceNode returns the source node of the VertexEdgeDTO.
func (v VertexEdgeDTO) SourceNode() string {
	return v.Node
}

// Implements interface AdjacencyListEntry from trustmodelstructure
// TargetNodes returns the target nodes of the VertexEdgeDTO.
func (v VertexEdgeDTO) TargetNodes() []string {
	return v.Links
}
