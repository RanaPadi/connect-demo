### OperationType
This is a constant that contains the meta tags of the operations

``` go 
  const (
    AGENT_OPERATOR_TRUST_FILTERING = "META_TRUST_DISCOUNT"
    AGENT_OPERATOR_TRUST_MERGE     = "META_TRUST_FUSION"
  )

  var LogicOperator = map[string]map[string]string{
    "subjectiveLogic": SubjectiveLogic,
  }

  var SubjectiveLogic = map[string]string{
    "META_TRUST_FUSION":   "FUSION",
    "META_TRUST_DISCOUNT": "DISCOUNT",
  }
```
