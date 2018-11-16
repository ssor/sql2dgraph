# Instruction

## How to use search

## Rules

### 1. Get Transactions

Input below words into search:
```
(tx=1-5&block=1)
```

the result is transactions from 1 to 5 in block 1

### 2. Get Blocks

Input below words into search:
```
(block=1-5)
```

the result is blocks from 1 to 5


### 3. Keyword Search

Just input some keyword, it will go.

The system will search block hash, transaction hash, transaction content

## Example

```
DDA6(block=2&tx=0-2)
```
This example means searching with keyword "DDA6"

```
(block=2&tx=0-2)
```
This example means query transactions index from 0 to 2 in block 2

```
(block=2-3)
```
This example means query  blocks from height 2 to 3

