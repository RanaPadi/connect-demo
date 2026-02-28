The `ToEdgeFromNodeAndLinks` function processes a list of `dtos.EdgeDTO` objects to create a structured representation of nodes and their connections (links). This structured representation is returned as a list of maps.

## Parameters

- `input`: A list of `dtos.EdgeDTO` objects representing edges in a graph.

## Flow of Execution

1. The function initializes a map called `nodes` to store information about nodes and their connections.
2. It then iterates through each `dtos.EdgeDTO` object in the `input` list.
3. For each edge, it extracts the `fromNode` and `toNode`.
4. If the `fromNode` does not exist in the `nodes` map, it creates an entry for it and initializes an inner map for its connections.
5. It records the connection from the `fromNode` to the `toNode` in the inner map.
6. The function constructs a list called `result` to store structured node data.
7. It iterates through the `nodes` map to create structured data for each node and its connections.
8. For each node, it creates a map containing:
   - `"node"`: The node's name.
   - `"links"`: A list of strings representing the names of nodes it is connected to.
9. The connections (links) are appended to the `"links"` list within each node's map.
10. The node data map is appended to the `result` list.
11. Finally, the function returns the `result`, which contains structured data about nodes and their connections.

This function is used to transform a list of edge data into a structured representation of nodes and their connections, which can be useful for various graph-related analyses and visualization.

```go
func ToEdgeFromNodeAndLinks(input []dtos.EdgeDTO) []map[string]interface{} {
	nodes := make(map[string]map[string]bool)

	for _, edge := range input {
		fromNode := edge.FromNode
		toNode := edge.ToNode

		if _, ok := nodes[fromNode]; !ok {
			nodes[fromNode] = make(map[string]bool)
		}
		nodes[fromNode][toNode] = true
	}

	result := []map[string]interface{}{}
	for node, links := range nodes {
		nodeData := map[string]interface{}{
			"node":  node,
			"links": make([]string, 0, len(links)),
		}

		for link := range links {
			nodeData["links"] = append(nodeData["links"].([]string), link)
		}

		result = append(result, nodeData)
	}

	return result