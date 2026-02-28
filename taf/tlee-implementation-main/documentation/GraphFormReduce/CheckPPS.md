This public function accepts two parameters: `subgraph` and `source` . It operates on a `subgraph` represented as a graph with nodes and string edges and checks the feasibility of a path from the `source` node within the subgraph.

### Parameters

- `subgraph`: A graph with string nodes and string edges.
- `source`: The source node for path control.

### Return values

- `bool`: A boolean value indicating whether a feasible path exists from the `source` node within the `subgraph`.

### Execution flow

1. The function calls another function `calculateWeights` that will return the computed flow.

2. The `calculateWeights` function starts by obtaining the adjacency map of the `subgraph` and stores it in `SubGraphAdjacencyMap`.
2. It initializes two maps, `nodeList` and `weight`, to keep track of vertex properties during execution.
3. For each node in the `SubGraphAdjacencyMap`, initialize the properties in `weight` with zero input and output flow.
4. Initialize the `source` node in `weight` with input flow equal to 1 and output flow equal to the division between 1 and the length of `SubGraphAdjacencyMap[source]` which would be the number of nodes exiting the source
5. We populate our string array `nodeList` where we save all the children of the source.
6. We create a loop that will run until `nodeList` is empty.
7. Iteratively process the nodes using a loop:
   - We initialize an empty `nodeListSon` array.
   - For each node in `nodeList`, compute in-flow and out-flow properties based on adjacent nodes and update `weight` using the `updateWeights` function.
6. Iteratively process the nodes using another loop:
   - For each node in the `nodeList` we go to find the children and if inside our `nodeListSon` array they do not exist, check done with the `containsNode` function we go to add them to `nodeListSon`
7. At this point we set our `nodeList` with the previously saved children in `nodeListSon` and repeat the steps from step 6. until the last one has no more children and `nodeListSon` is empty.
8. Once the execution is finished the array containing the stream in all nodes from the source to the target will be returned.
9. Through the `countInFlow` function we are going to see how many nodes have a flow greater than 0.999999, excluding source and target, and return true if they are <= 2 or false.

**Snippet**

```go
func CheckPPS(subgraph graph.Graph[string, string], source string) bool {
 flow := calculateWeights(subgraph, source)
 return (len(flow) > 0 && countInFlow(flow) <= 2)
}


