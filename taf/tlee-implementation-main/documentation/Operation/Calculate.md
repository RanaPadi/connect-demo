## Function `Calculate`

The private function takes in a slice of floating-point numbers representing opinions and returns a modified slice of opinions. It ensures that each opinion value falls within the range [0, 1] and rounds them to two decimal places.

### Parameters:

- `opinion`: A slice of floating-point numbers representing opinions.

### Flow of Execution:
1. The function iterates through each `op` (opinion) in the `opinion` slice:
   - For each `opinion`, it checks if it is within the range [0, 1]:
     - If the opinion is greater than or equal to 1, it sets the opinion to 1.
     - If the opinion is less than 0, it sets the opinion to 0.
   - It rounds the opinion to two decimal places using `math.Round(op*100) / 100`.
   - It appends the modified opinion to the `tempOpinion` slice.

2. Finally, the function returns the `tempOpinion` slice containing the modified opinion values.

**Snippet**

```go
func calculate(opinion []float64) []float64 {

	var tempOpinion []float64

	for _, op := range opinion {
		if !((op < 1) && (op > 0)) {
			if op >= 1 {
				op = 1
			} else {
				op = 0
			}
		}
		op = math.Round(op*100) / 100
		tempOpinion = append(tempOpinion, op)
	}
	return tempOpinion
}
