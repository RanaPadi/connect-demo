package dtos

type ProcessingPath struct {
	Edges     []EdgeDTO
	Nodes     []string
	NodeToPPS map[string]EdgeDTO
	Index     int
}
