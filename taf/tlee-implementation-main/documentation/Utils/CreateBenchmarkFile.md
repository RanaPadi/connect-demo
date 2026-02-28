## Function `CreateBenchmarkFile`

### Description
The function `CreateBenchmarkFile` is responsible for creating or updating a benchmark file with information about the processing time, the number of input edges, and the compression rate. The benchmark information is stored as a JSON array of `dtos.BenchMark` objects.

### Parameters
- `time` (type: `float64`, required): The processing time for the operation.
- `numEdgeInput` (type: `int64`, required): The number of input edges.
- `numEdgeOutput` (type: `int64`, required): The number of output edges after compression.

### Execution Flow
1. Set the default filename for the benchmark file as "benchmark.json".
2. Read the content of the existing benchmark file if it exists. If not, initialize an empty `benchMark` slice.
3. Unmarshal the existing file content into the `benchMark` slice.
4. Append a new `dtos.BenchMark` object to the `benchMark` slice, containing information about processing time, input edges, and the compression rate.
5. Open or create the benchmark file for writing.
6. Marshal the updated `benchMark` slice into JSON format.
7. Write the JSON content to the benchmark file, appending a newline for better readability.

### Snippet
```go
func CreateBenchmarkFile(time float64, numEdgeInput int64, numEdgeOutput int64) {

	var fileName string = "banchmark.json"
	var benchMark []dtos.BenchMark

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {

		benchMark = make([]dtos.BenchMark, 0)

	} else {

		err = json.Unmarshal(fileContent, &benchMark)
		Must(err)

	}

	benchMark = append(benchMark, dtos.BenchMark{Time: time, NumEdge: numEdgeInput, CompressionRate: float64(numEdgeOutput) / float64(numEdgeInput)})

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	Must(err)

	defer file.Close()

	benchmarkJSON, err := json.Marshal(benchMark)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n\n", benchmarkJSON))
	Must(err)
}
