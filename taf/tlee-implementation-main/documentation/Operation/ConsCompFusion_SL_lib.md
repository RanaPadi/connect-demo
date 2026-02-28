## Function `Concensus&CompromiseFusion` from `SL library`

As part of the `SL library` the operator for `Concensus&CompromiseFusion` is not implemented. Therefore, we use the logic of the `ConstraintFusion` operator instead. More details about the `ConstraintFusion` operator can be found in `ConstraintFusion_SL_lib.md`.

**Snippet**

```go
func ConsCompFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.ConstraintFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}
```
