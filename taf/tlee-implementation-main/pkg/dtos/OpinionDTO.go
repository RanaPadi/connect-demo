package dtos

import "github.com/vs-uulm/go-subjectivelogic/pkg/subjectivelogic"

type OpinionDTO struct {
	Value    OpinionDTOValue
	FromNode string `csv:"fromNode" json:"fromNode,omitempty"`
	ToNode   string `csv:"toNode" json:"toNode,omitempty"`
}

// Implements interface TrustRelationship from trustmodelstructure
// Source() returns FromNode from OpinionDTO
func (o OpinionDTO) Source() string {
	return o.FromNode
}

// Implements interface TrustRelationship from trustmodelstructure
// Destination() returns ToNode from OpinionDTO
func (o OpinionDTO) Destination() string {
	return o.ToNode
}

// Implements interface TrustRelationship from trustmodelstructure
// Opinion() returns Value from OpinionDTO
func (o OpinionDTO) Opinion() subjectivelogic.QueryableOpinion {
	return o.Value
}
