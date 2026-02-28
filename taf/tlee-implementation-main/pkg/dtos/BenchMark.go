package dtos

type BenchMark struct {
	NumEdge         int64   `json:"n_edge"`
	CompressionRate float64 `json:"compressionRate"`
	Time            float64 `json:"time"`
}
