## Function `AveragingFusion` from `SL library`

This function performs **averaging fusion** (consensus) on two subjective opinions (`x` and `y`) using operators from the **Subjective Logic (SL) Library**. It combines the opinions into a single opinion by averaging their beliefs, uncertainties, and base rates, following the principles of probabilistic reasoning under uncertainty.

### SL Library Integration
The function directly leverages the following SL Library components:
1. **`NewOpinion`**: Constructs SL-compatible opinions from input DTOs.
2. **`AveragingFusion`**: Core SL operator that computes the consensus opinion by averaging inputs.


### Parameters
- **`x`**: `dtos.OpinionDTOValue`  
  First input opinion (belief, disbelief, uncertainty, base rate).
- **`y`**: `dtos.OpinionDTOValue`  
  Second input opinion (belief, disbelief, uncertainty, base rate).

### Return Value
Returns a new `dtos.OpinionDTOValue` containing the fused:
- Belief (`b`), Disbelief (`d`), Uncertainty (`u`), and Base Rate (`a`).


### Flow of Execution
1. **Convert Inputs to SL Opinions**:
    - Uses `subjectivelogic.NewOpinion()` to convert `x` and `y` into SL-native opinion objects.
2. **Apply Fusion**:
    - Calls `subjectivelogic.AveragingFusion()` to compute the consensus opinion.
3. **Convert Back to DTO**:
    - Maps the fused SL opinion back to a `dtos.OpinionDTOValue`.

**Snippet**

```go
func AveragingFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.AveragingFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}
```