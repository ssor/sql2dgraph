# sql2dgraph
demo of graphQL by convert mysql sample data to rdf sets and push into dgraph

## Must Know

- Constraint will be outside of system
- Indices with same name should be the same index



## Start Dgraph

```
dgraph zero

dgraph server --lru_mb 2048 --zero localhost:5080

dgraph-ratel
```
