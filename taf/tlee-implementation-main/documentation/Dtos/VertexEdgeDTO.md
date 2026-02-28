### VertexEdgeDTO
This Dto is used to perform an `ExpresisonSyntesizer` endpoint, is an array of json object that represent a graph which elements are `node` and the list of the correlated edges `link`


``` json 
{
  "Node": "string", 
  "Links": "string[]"
}
```

## Real Example
``` json
[
    {
        "node": "a",
        "links": [
            "b",
            "d",
            "g",
            "s"
        ]
    },
    {
        "node": "b",
        "links": [
            "h",
            "j",
            "l"
        ]
    }
]
