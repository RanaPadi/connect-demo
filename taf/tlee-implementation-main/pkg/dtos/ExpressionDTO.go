package dtos

type ExpressionDTO struct {
	Operation string `json:"operation,omitempty"`
	FromNode  string `json:"fromNode,omitempty"`
	ToNode    string `json:"toNode,omitempty"`
}
