package dtos

import (
	"github.com/dominikbraun/graph"
)

type SynthesizingGraph struct {
	Adj       map[string]map[string]graph.Edge[string]
	Edges     []EdgeDTO
	Nodes     []string
	NodeToPPS map[string]EdgeDTO
}
