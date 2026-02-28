package operation

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/config"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"errors"
	"github.com/vs-uulm/go-subjectivelogic/pkg/subjectivelogic"
	"github.com/vs-uulm/taf-tlee-interface/pkg/trustmodelstructure"
	"math"
)

/*
	The function GetFusionOperator called to retrieve a specific fusion operation based on the provided FusionOperator type.
	Once retrieved, the caller can use the returned function to combine or "fuse" two OpinionDTOValue instances.
	The error return allows the caller to handle cases where the fusion operator might be invalid, misconfigured, or unsupported.

	1. Takes a single argument `op` as input:
	The function expects one input parameter, op, of type trustmodelstructure.FusionOperator. This is a type that defines how a fusion operator behaves.

	2. Returns Two Values:

	First Return Value: A function (i.e., a higher-order function) that takes two parameters of type dtos.OpinionDTOValue and returns one dtos.OpinionDTOValue.
	This returned function is likely designed to perform a specific "fusion" operation (like combining or averaging) on two OpinionDTOValue values.

	Second Return Value: An error value, which indicates whether there was an issue with getting or creating the fusion operator.
	If the error is nil, then the operation was successful; otherwise, it contains information about the failure.
*/

func GetFusionOperator(op trustmodelstructure.FusionOperator) (func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, error) {
	if op == trustmodelstructure.AveragingFusion {
		return AveragingFusion, nil
	} else if op == trustmodelstructure.WeightedFusion {
		return WeightedFusion, nil
	} else if op == trustmodelstructure.CumulativeFusion {
		return CumulativeFusion, nil
	} else if op == trustmodelstructure.ConstraintFusion {
		return ConstraintFusion, nil
	} else if op == trustmodelstructure.ConsensusAndCompromiseFusion {
		return ConsCompFusion, nil
	} else {
		config.Logger.Error("Non supporting fusion operator", "value", op)
		return nil, errors.New("Non supporting fusion operator")
	}
}

/*
Same as for fusion operator, but here two discount operators. first one for functional and second one for referral trust.
*/
func GetDiscountOperator(op trustmodelstructure.DiscountOperator) (func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, func(dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue, error) {
	if op == trustmodelstructure.DefaultDiscount {
		return Discount, DiscountRef, nil
	} else if op == trustmodelstructure.OppositeBeliefDiscount {
		return DiscountingOppositeBelief, DiscountingOppositeBelief, nil
	} else {
		config.Logger.Error("Non supporting Discount operator", "value", op)
		return nil, nil, errors.New("Non supporting Discount operator")
	}
}

/*
	The following functions enable the use of operators from the SL Library
*/

// Discount operator
func Discount(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.TrustDiscounting(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())

}

// Our newly proposed DiscountRef operator (cf. FUSION paper)
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

// DiscountingOppositeBelief operator
func DiscountingOppositeBelief(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	config.Logger.Debug("APP OPPOSITE DISCOUNT")
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	//op, _ := subjectivelogic.TrustDiscountingOppositeBelief(&opx, &opy)
	op, _ := subjectivelogic.TrustDiscountingOppositeBelief(&opx, &opy)
	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

// AveragingFusion Operator
func AveragingFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.AveragingFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

// WeightedFusion Operator
func WeightedFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.WeightedFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

// CumulativeFusion Operator
func CumulativeFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.CumulativeFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

// ConstraintFusion Operator
func ConstraintFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.ConstraintFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

// Concensus&CompromiseFusion Operator
/* Note: As part of the `SL library` the operator for `Concensus&CompromiseFusion` is not implemented. Therefore, we use the logic of the `ConstraintFusion` operator instead. */
func ConsCompFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	opx, _ := subjectivelogic.NewOpinion(x.Belief(), x.Disbelief(),
		x.Uncertainty(), x.BaseRate())
	opy, _ := subjectivelogic.NewOpinion(y.Belief(), y.Disbelief(),
		y.Uncertainty(), y.BaseRate())
	op, _ := subjectivelogic.ConstraintFusion(&opx, &opy)

	return dtos.NewOpinionDTOValue(op.Belief(), op.Disbelief(), op.Uncertainty(), op.BaseRate())
}

/*
The private function takes in a slice of floating-point numbers representing opinions and returns a modified slice of opinions.
It ensures that each opinion value falls within the range [0, 1] and rounds them to two decimal places.

The function iterates through each op (opinion) in the opinion slice:

For each opinion, it checks if it is within the range [0, 1]:

If the opinion is greater than or equal to 1, it sets the opinion to 1.
If the opinion is less than 0, it sets the opinion to 0.

It rounds the opinion to two decimal places using math.Round(op*100) / 100.
It appends the modified opinion to the tempOpinion slice.

Finally, the function returns the tempOpinion slice containing the modified opinion values.
*/
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

/*
	Our custom-implemented SL operators. Discount, AveragingFusion and BeliefConstraintFusion. These operators are not used as part of the CONNECT projects.
*/

/*
Declaration of OperationFunc variable.
OperationFunc is a map where each **key is a string** and each **value is a function with the signature (dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue**.

Currently OperationFunc is initialized is initialized with two key-value pairs:
  - "FUSION": Associated with the AveragingFusion function.
  - "DISCOUNT": Associated with the Discount function.

Both AveragingFusion and Discount functions must match the signature (dtos.OpinionDTOValue, dtos.OpinionDTOValue) dtos.OpinionDTOValue.
*/

/*
	Discount operator

	The public function calculates a discounted opinion (dtos.Key) based on two input opinions (dtos.Key).
	It computes the discounted belief, disbelief, uncertainty, base rate,
	and projected probability using specific formulas and then rounds the values to two decimal places.

	The function extracts various values from the input opinions x and y, such as a (base rate of y), b (product of belief of y and projected
	probability of x), d (product of disbelief of y and projected probability of x),
	u (complement of the product of projected probability of x and the sum of belief and disbelief of y), and e (the sum of b and a multiplied by u).
	The extracted values are used to calculate a new opinion in the form of a float64 slice opinion, containing belief, disbelief, uncertainty, base rate, and projected probability.
	The calculate function is called to ensure that each value in the opinion slice falls within the range [0, 1] and is rounded to two decimal places.
	The calculated e is also rounded to two decimal places using math.Round.
	The function returns a dtos.Key representing the discounted opinion with the modified belief, disbelief, uncertainty, base rate, and projected probability.
*/

/* func Discount(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	//fmt.Println("APP DISCOUNT")
	config.Logger.Debug("APP DISCOUNT")
	a := y.BaseRate()
	b := x.ProjectedProbability() * y.Belief()
	d := x.ProjectedProbability() * y.Disbelief()
	u := (1 - x.ProjectedProbability()*(y.Belief()+y.Disbelief()))
	e := b + a*u

	opinion := []float64{b, d, u, a}
	opinion = calculate(opinion)
	e = math.Round(e*100) / 100

	return dtos.NewOpinionDTOValue(opinion[0], opinion[1], opinion[2], opinion[3])
} */

/*
	AveragingFusion

	The public function performs averaging fusion on two input opinions (dtos.Key) to calculate a fused opinion.
	It considers various cases based on the uncertainty values of the input opinions and
	computes the fused belief, disbelief, uncertainty, base rate, and projected probability, rounding the values to two decimal places.

	The function initializes variables b (belief), u (uncertainty), and a (base rate).
	It checks whether either x or y has non-zero uncertainty. If either of them does, it computes the fused belief (b), fused uncertainty (u), and fused base rate (a) using specific formulas that consider the uncertainty values of both opinions.
	If both x and y have zero uncertainty, it computes the fused belief (b) as the average of the beliefs of x and y, sets u to 0 (no uncertainty), and computes the fused base rate (a) as the average of the base rates of x and y.
	It calculates the fused disbelief (d) as the complement of the sum of fused belief (b) and fused uncertainty (u).
	The fused projected probability (e) is computed as the sum of fused belief (b) and fused base rate (a) multiplied by fused uncertainty (u).
	All numeric values (b, d, u, e, and a) are rounded to two decimal places using math.Round.
	The function returns a dtos.Key representing the fused opinion with the modified belief, disbelief, uncertainty, base rate, and projected probability.
*/

/* func AveragingFusion(x dtos.OpinionDTOValue, y dtos.OpinionDTOValue) dtos.OpinionDTOValue {
	var b float64
	var u float64
	var a float64
	if x.Uncertainty() != 0 || y.Uncertainty() != 0 {
		b = (x.Belief()*y.Uncertainty() + y.Belief()*x.Uncertainty()) /
			(x.Uncertainty() + y.Uncertainty())
		u = 2 * x.Uncertainty() * y.Uncertainty() / (x.Uncertainty() + y.Uncertainty())
		a = (x.BaseRate() + y.BaseRate()) / 2
	} else {
		b = 0.5 * (x.Belief() + y.Belief())
		u = 0
		a = 0.5 * (x.BaseRate() + y.BaseRate())
	}

	d := (1 - u - b)
	e := b + a*u

	b = math.Round(b*100) / 100
	d = math.Round(d*100) / 100
	u = math.Round(u*100) / 100
	e = math.Round(e*100) / 100
	a = math.Round(a*100) / 100
	return dtos.NewOpinionDTOValue(b, d, u, a)
} */

/*
	BeliefConstraintFusion

	This fusion operator takes two OpinionDTOValue and returns the result of the belief constraint fusion.

	BeliefConstraintFusion calculates the fused belief, disbelief, uncertainty, and base rate based on the given OpinionDTOs xA and xB.
	It uses harmony and conflict functions to derive the belief and uncertainty values, and adjusts the base rate depending on the base rates of xA and xB.

	The formula for belief (b) incorporates harmony and conflict between the two opinions, while the uncertainty (u) accounts for both opinions' uncertainties.
	The final base rate (a) is calculated based on the average base rates of xA and xB, with special handling for cases where both base rates are equal to 1.

	Returns a new OpinionDTO containing the fused belief, disbelief, uncertainty, and base rate.

*/

/* func BeliefConstraintFusion(xA, xB dtos.OpinionDTO) dtos.OpinionDTO {

	har := harmony(xA, xB)
	con := conflict(xA, xB)

	b := (har) / (1 - con)
	u := ((xA.Uncertainty*xB.Uncertainty) + har*(xA.Uncertainty+xB.Uncertainty)) / (2 - (xA.Uncertainty + xB.Uncertainty))
	d:= 1-(b+u)
	var a float64
	if xA.BaseRate + xB.BaseRate<2 {
		a = (xA.BaseRate*(1-xA.Uncertainty)+xB.BaseRate*(1-xB.Uncertainty)) /
		(2 - (xA.Uncertainty + xB.Uncertainty))
	} else { // When base rates are equal to 1
		a = (xA.BaseRate + xB.BaseRate) / 2
	}

	result := dtos.OpinionDTO{
		Belief:     b,
		Disbelief:  d,
		Uncertainty: u,
		BaseRate:   a,
		// ProjectedProbability() and other fields need to be calculated or assigned as needed
	}

	return result
}

	// harmony calculates the harmony between two opinions.
func harmony(xA, xB dtos.OpinionDTO) float64 {
	return (xA.BaseRate * xB.Uncertainty) + (xB.BaseRate * xA.Uncertainty)
}

	// conflict calculates the conflict between two opinions.
func conflict(xA, xB dtos.OpinionDTO) float64 {
	return math.Abs(xA.BaseRate - xB.BaseRate)
} */
