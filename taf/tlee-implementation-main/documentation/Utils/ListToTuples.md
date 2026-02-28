The public function converts a list of strings into a list of string tuples, where each tuple contains two consecutive elements from the input list.

### Parameters:

- `aList`: A list of strings.

### Flow of Execution:
1. The function initializes an empty list `n` to store the resulting tuples.

2. It iterates through the input list `aList` up to the second-to-last element (excluding the last element) using a `for` loop.

3. In each iteration, it appends a tuple to the `n` list containing two consecutive elements from the input list.

4. Finally, the function returns the list of string tuples.

**Snippet**

```go
func ListToTuples(aList []string) [][]string {
	n := make([][]string, 0)
	for i := 0; i < len(aList)-1; i++ {
		n = append(n, []string{aList[i], aList[i+1]})
	}
	return n
}
