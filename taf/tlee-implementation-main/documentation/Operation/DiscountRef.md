## Function `DiscountRef`

This function implements **reference discounting** between two subjective opinions (`x` and `y`) using core subjective logic operations. It calculates how opinion `x` should be discounted based on opinion `y`'s reliability, producing a new opinion that reflects this trust relationship.

### Parameters
- **`x`**: `dtos.OpinionDTOValue`  
  The opinion being discounted
- **`y`**: `dtos.OpinionDTOValue`  
  The reference opinion providing the discount factors


### Return Value
Returns a new `dtos.OpinionDTOValue` with:
- Discounted belief (`b`)
- Adjusted disbelief (`d`)
- Calculated uncertainty (`u`)
- Projected expectation (`e`)

### Execution Flow

**Step 1: Initial Discounting Calculation**
```go
b := x.Belief() * y.Belief()       // Discounted belief
d := x.Belief() * y.Disbelief()    // Discounted disbelief
u := 1 - (b + d)                   // Residual uncertainty
```
- Multiplies the belief components of both opinions
- Calculates the discounted disbelief
- Derives remaining uncertainty

**Step 2: Base Rate Adjustment**
```go
numerator := (x.Belief()+x.Uncertainty()*x.BaseRate())*
             (y.Belief()+y.Uncertainty()*y.BaseRate()) - 
             (x.Belief() * y.Belief())

denominator := 1 - x.Belief()*(y.Belief()+y.Disbelief())

a := numerator / denominator
```
1. Combines beliefs with their uncertainty-weighted base rates
2. Subtracts the simple belief product
3. Normalizes by the complement of x's belief over y's positive/negative components

**Step 3: Expectation Projection**
```go
e := b + a*u  // Projected expectation value
```
- Combines discounted belief with adjusted uncertainty
- Represents the final probability expectation

**Step 4: Opinion Normalization**
```go
opinion := []float64{b, d, u, a}
opinion = calculate(opinion)  // Normalizes and validates values
```
- Packages components into a temporary array
- Applies normalization via helper function to ensure:
    - All values ∈ [0,1]
    - b + d + u = 1

**Step 5: Precision Control**
```go
e = math.Round(e*100)/100  // Rounds to 2 decimal places
```
- Ensures consistent numerical representation
- Matches financial/conventional reporting standards

**Step 6: Result Packaging**
```go
return dtos.NewOpinionDTOValue(opinion[0], opinion[1], opinion[2], opinion[3])
```
- Constructs final opinion DTO with:
    - Position 0: Discounted belief
    - Position 1: Discounted disbelief
    - Position 2: Normalized uncertainty
    - Position 3: Adjusted base rate


**Snippet**

```go
func DiscountRef(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	config.Logger.Debug("APP DISCOUNT REF")
	b := x.Belief() * y.Belief()
	d := x.Belief() * y.Disbelief()
	u := 1 - (b + d)

	a := ((x.Belief()+x.Uncertainty()*x.BaseRate())*
		(y.Belief()+y.Uncertainty()*y.BaseRate()) -
		(x.Belief() * y.Belief())) / (1 - x.Belief()*(y.Belief()+y.Disbelief()))
	e := b + a*u

	opinion := []float64{b, d, u, a}

	opinion = calculate(opinion)

	e = math.Round(e*100) / 100

	return dtos.NewOpinionDTOValue(opinion[0], opinion[1], opinion[2], opinion[3])
}
```
