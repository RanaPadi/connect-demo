The private function takes in two parameters: `adjacencyMap` and `node`. It is responsible for counting the number of incoming links (edges) to a specified `node` within a graph represented by the `adjacencyMap`.

### Parameters:

- `adjacencyMap`: A map of strings to maps of strings to `Edge[string]` representing the adjacency information of a graph.
- `node`: A string representing the target node for which the incoming links need to be counted.

### Return Value:

- `linksCount`: An integer representing the count of incoming links to the specified `node`.

### Flow of Execution:
1. The function initializes an integer variable `linksCount` to zero, which will store the count of incoming links.
2. It then iterates through each entry in the `adjacencyMap`:
   - For each entry, it represents a node (key) and its associated outgoing links (value).
   - It further iterates through the outgoing links for that node.
   - If a link in the outgoing links matches the specified `node`, it increments the `linksCount` by 1.
3. After iterating through all entries in the `adjacencyMap`, the function returns the `linksCount`, which represents the count of incoming links to the specified `node`.

**Snippet**

```go
func linksCounter(adjacencyMap map[string]map[string]graph.Edge[string], node string) int {
	var linksCount = 0
	for _, links := range adjacencyMap {
		for link := range links {
			if link == node {
				linksCount += 1
			}
		}
	}
	return linksCount
}

