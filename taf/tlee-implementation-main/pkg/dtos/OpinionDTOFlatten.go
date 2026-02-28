package dtos

type OpinionDTOFlatten struct {
	Belief      float64 `csv:"belief" json:"belief"`
	Disbelief   float64 `csv:"disbelief" json:"disbelief"`
	Uncertainty float64 `csv:"uncertainty" json:"uncertainty"`
	BaseRate    float64 `csv:"baseRate" json:"baseRate"`
	FromNode    string  `csv:"fromNode" json:"fromNode,omitempty"`
	ToNode      string  `csv:"toNode" json:"toNode,omitempty"`
}
