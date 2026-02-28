package dtos

type KeyDTO struct {
	Operation string          `json:"operation,omitempty"`
	FromNode  string          `json:"fromNode,omitempty"`
	ToNode    string          `json:"toNode,omitempty"`
	Opinion   OpinionDTOValue `json:"opinion,omitempty"`
}
