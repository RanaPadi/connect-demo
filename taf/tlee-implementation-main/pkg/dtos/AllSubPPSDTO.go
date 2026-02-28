package dtos

import "github.com/dominikbraun/graph"

type AllSubPPSDTO struct {
	Graph           graph.Graph[string, string]
	Source          string
	Target          string
	MinNestingLevel int
}
