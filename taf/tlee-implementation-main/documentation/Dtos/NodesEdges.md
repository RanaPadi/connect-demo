### NodesEdges
This dto,is used as request body in the `DspgGenerator` POST endpoint to specify the size of the grpah, essentially defining the number of nodes and the number of links each node has.
``` json 
{
  "NumNodes": "int", 
  "NumEdges": "int"
}
```
# Real example 
``` json 
{
  "NumNodes":10, 
  "NumEdges": 50
}
```