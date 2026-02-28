package dtos

import "fmt"

type OpinionDTOValue struct {
	belief      float64 `csv:"belief" json:"belief"`
	disbelief   float64 `csv:"disbelief" json:"disbelief"`
	uncertainty float64 `csv:"uncertainty" json:"uncertainty"`
	baseRate    float64 `csv:"baseRate" json:"baseRate"`
}

func NewOpinionDTOValue(belief float64, disbelief float64, uncertainty float64, baseRate float64) OpinionDTOValue {
	return OpinionDTOValue{
		belief:      belief,
		disbelief:   disbelief,
		uncertainty: uncertainty,
		baseRate:    baseRate,
	}
}

func (o OpinionDTOValue) ProjectedProbability() float64 {
	return o.belief + o.uncertainty*o.baseRate
}

// Implements interface QueryableOpinion from subjectivelogic
// GetBelief() returns Belief from OpinionDTOValue
func (o OpinionDTOValue) Belief() float64 {
	return o.belief
}

// Implements interface QueryableOpinion from subjectivelogic
// GetDisbelief() returns Disbelief from OpinionDTOValue
func (o OpinionDTOValue) Disbelief() float64 {
	return o.disbelief
}

// Implements interface QueryableOpinion from subjectivelogic
// GetUncertainty() returns Uncertainty from OpinionDTOValue
func (o OpinionDTOValue) Uncertainty() float64 {
	return o.uncertainty
}

// Implements interface QueryableOpinion from subjectivelogic
// GetBaseRate() returns BaseRate from OpinionDTOValue
func (o OpinionDTOValue) BaseRate() float64 {
	return o.baseRate
}

func (o OpinionDTOValue) String() string {
	return fmt.Sprintf("Belief: %v, Disbelief: %v, Uncertainty: %v, Base Rate: %v", o.belief, o.disbelief, o.uncertainty, o.baseRate)
}
