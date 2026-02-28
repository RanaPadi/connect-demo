### InputEvaluatorDTO
This dto, is similar to the [DataChildExpressionDTO](DataChildExpressionDTO.md) and will be used as input body to the `Evaluetor` endpoint, but contains the real expression definition, not the meta-tag, and in the field `opionionMode`, will be passed the path of the csv file, that will contains the opions values that will be used to evaluare the graph
``` json 
{
  "OpinionMode": "string", 
  "Expression": "DataChildDTO"
}
```
# Real example 

``` json

{
    "opinionMode": "conf/csv/opinions.csv",
    "expression":{
    "data": {
        "operation": "FUSION"
    },
    "child": [
        {
            "data": {
                "operation": "DISCOUNT"
            },
            "child": [
                {
                    "data": {
                        "fromNode": "a",
                        "toNode": "d"
                    }
                },
                {
                    "data": {
                        "operation": "FUSION"
                    },
                    "child": [
                        {
                            "data": {
                                "operation": "DISCOUNT"
                            },
                            "child": [
                                {
                                    "data": {
                                        "fromNode": "d",
                                        "toNode": "w"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "w",
                                        "toNode": "x"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "x",
                                        "toNode": "f"
                                    }
                                }
                            ]
                        },
                        {
                            "data": {
                                "operation": "DISCOUNT"
                            },
                            "child": [
                                {
                                    "data": {
                                        "fromNode": "d",
                                        "toNode": "z"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "z",
                                        "toNode": "f"
                                    }
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        {
            "data": {
                "operation": "DISCOUNT"
            },
            "child": [
                {
                    "data": {
                        "fromNode": "a",
                        "toNode": "g"
                    }
                },
                {
                    "data": {
                        "fromNode": "g",
                        "toNode": "f"
                    }
                }
            ]
        },
        {
            "data": {
                "operation": "DISCOUNT"
            },
            "child": [
                {
                    "data": {
                        "fromNode": "a",
                        "toNode": "s"
                    }
                },
                {
                    "data": {
                        "fromNode": "s",
                        "toNode": "f"
                    }
                }
            ]
        },
        {
            "data": {
                "operation": "DISCOUNT"
            },
            "child": [
                {
                    "data": {
                        "fromNode": "a",
                        "toNode": "b"
                    }
                },
                {
                    "data": {
                        "operation": "FUSION"
                    },
                    "child": [
                        {
                            "data": {
                                "operation": "DISCOUNT"
                            },
                            "child": [
                                {
                                    "data": {
                                        "fromNode": "b",
                                        "toNode": "h"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "h",
                                        "toNode": "f"
                                    }
                                }
                            ]
                        },
                        {
                            "data": {
                                "operation": "DISCOUNT"
                            },
                            "child": [
                                {
                                    "data": {
                                        "fromNode": "b",
                                        "toNode": "j"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "j",
                                        "toNode": "f"
                                    }
                                }
                            ]
                        },
                        {
                            "data": {
                                "operation": "DISCOUNT"
                            },
                            "child": [
                                {
                                    "data": {
                                        "fromNode": "b",
                                        "toNode": "l"
                                    }
                                },
                                {
                                    "data": {
                                        "fromNode": "l",
                                        "toNode": "f"
                                    }
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}
}