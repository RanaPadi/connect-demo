package dtos

type InputEvaluatorDTO struct {
	OpinionMode string       `json:"opinionMode"`
	Expression  DataChildDTO `json:"expression"`
}
