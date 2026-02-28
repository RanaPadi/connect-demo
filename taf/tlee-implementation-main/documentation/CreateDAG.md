## Method `CreateDAG`

### Description
This function generates a random non DSPG Directed Acyclic Graph (DAG) represented by an adjacency matrix. It takes the number of nodes (`numNodes`) and the desired number of edges (`numEdges`) as input parameters. The generated DAG is then returned as a 2D slice representing the adjacency matrix.

### Parameters
- `numNodes` (type: `int`, required): The number of nodes in the graph.
- `numEdges` (type: `int`, required): The desired number of edges in the graph.

### Return Type
- `[][]int`: A 2D slice representing the adjacency matrix of the generated DAG.

### Execution Flow
1. The function initializes an empty adjacency matrix (`adjMatrix`) with dimensions `numNodes x numNodes`.
2. It seeds the random number generator using the current Unix timestamp to introduce randomness.
3. An adjacency list (`adjList`) is initialized to store the neighbors of each node.
4. A topological order of nodes is generated randomly using `rand.Perm(numNodes)`.
5. The function iterates over each node in the topological order and connects it to subsequent nodes based on the desired number of random edges (`numEdges`).
6. Random edges are generated for each node, and connections are made in the adjacency matrix and adjacency list.
7. The function returns the generated adjacency matrix representing the DAG.

### Method Snippet
```go
func CreateDAG(numNodes, numEdges int) [][]int {
	adjMatrix := make([][]int, numNodes)
	for i := range adjMatrix {
		adjMatrix[i] = make([]int, numNodes)
	}

	rand.Seed(time.Now().UnixNano())

	adjList := make([][]int, numNodes)

	topologicalOrder := rand.Perm(numNodes)

	for i := 0; i < numNodes; i++ {
		randomEdge := dspgFunc.GenerateRandomNumber(1, numEdges)
		for j := i + 1; j < numNodes && len(adjList[topologicalOrder[i]]) < randomEdge; j++ {
			if adjMatrix[topologicalOrder[i]][topologicalOrder[j]] == 0 {
				adjMatrix[topologicalOrder[i]][topologicalOrder[j]] = 1
				adjList[topologicalOrder[i]] = append(adjList[topologicalOrder[i]], topologicalOrder[j])
			}
		}
	}

	return adjMatrix
}
