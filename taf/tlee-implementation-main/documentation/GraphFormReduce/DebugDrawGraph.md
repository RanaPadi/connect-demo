## Function `DebugDrawGraph`

### Description
The function `DebugDrawGraph` generates a Graphviz DOT file and an SVG file from a graph using the `graph` package. The DOT file is created in the "./debug/" directory with the specified name, and the SVG file is generated using the Graphviz `dot` command. This function is useful for visually debugging the structure of a graph.

### Parameters
- `g` (type: `graph.Graph[string, string]`, required): The graph to be visualized using Graphviz.
- `name` (type: `string`, required): The name used for creating the DOT and SVG files.

### Execution Flow
1. Create a new DOT file with the specified name in the "./debug/" directory.
2. Use the `draw.DOT` function from the `graph` package to generate the DOT representation of the graph and write it to the DOT file.
3. Execute the Graphviz `dot` command to convert the DOT file to an SVG file in the same directory.

### Snippet
```go
func DebugDrawGraph(g graph.Graph[string, string], name string) {
	graph, err := os.Create("./debug/" + name + ".gv")
	utils.Must(err)
	err = draw.DOT(g, graph)
	utils.Must(err)

	exec.Command("dot", "-Tsvg", "-O", "./debug/"+name+".gv").Output()
}
