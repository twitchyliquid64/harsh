# harsh
[![GoDoc](https://godoc.org/github.com/twitchyliquid64/harsh?status.svg)](http://godoc.org/github.com/twitchyliquid64/harsh) [![Build Status](https://travis-ci.org/twitchyliquid64/harsh.svg?branch=master)](https://travis-ci.org/twitchyliquid64/harsh) 

A AST library intended for code generation, mutation, and visualisation.

```shell
go build github.com/twitchyliquid64/harsh/debugprint
./debugprint src/github.com/twitchyliquid64/harsh/compiler/translate.go
```

Planned features, in various states are:

 * ast - represents most features of a code graph in an abstract way
  * exec - executes the code graph
  * print - prints the code graph to stdout
  * vis - Generate a graphical representation of the code graph [planned]

* compiler - translate a subset of Go to the representation used in `ast`.
* mutate - methods to mutate an existing graph or swap nodes from graphs/subgraphs (breed) [planned]
* generate - randomly generate graphs constrained by a complexity score [planned]


### Immediate TODO

#### AST

 - [x] Testing for new type system - arrays, array declaration, array subtype etc
 - [ ] Testing for type structure of parameters
 - [x] e2e test for new bool operators & arrays
 - [ ] Testing for exec errors
 - [ ] Testing for translate errors
 - [x] Implement array use - move Exec to use pointer to Variant instead of Variant. Write tests.
 - [x] Boolean operators
 - [ ] Unary operators (!)
