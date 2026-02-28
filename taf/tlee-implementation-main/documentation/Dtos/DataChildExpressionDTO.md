### DataChildExpressionDTO
This Dto, maybe the most complex one,is used in the `MetaToConcreteExpressionConverter` endpoint as request is a recursive dto, containing the Meta-operation `operation` and the relatives paths to evaluate, described as an object containing the source node `fromNode` and the target node `toNode`
``` json 
{
  "Data": "ExpressionDTO", 
  "Child": "DataChildExpressionDTO[]",
}
```
# Real example

```json
{
    "data": {
        "operation": "META_TRUST_FUSION"
    },
    "child": [
        {
            "data": {
                "operation": "META_TRUST_DISCOUNT"
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
                        "operation": "META_TRUST_FUSION"
                    },
                    "child": [
                        {
                            "data": {
                                "operation": "META_TRUST_DISCOUNT"
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
                                "operation": "META_TRUST_DISCOUNT"
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
                "operation": "META_TRUST_DISCOUNT"
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
                "operation": "META_TRUST_DISCOUNT"
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
                "operation": "META_TRUST_DISCOUNT"
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
                        "operation": "META_TRUST_FUSION"
                    },
                    "child": [
                        {
                            "data": {
                                "operation": "META_TRUST_DISCOUNT"
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
                                "operation": "META_TRUST_DISCOUNT"
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
                                "operation": "META_TRUST_DISCOUNT"
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