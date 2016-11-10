# harsh
[![GoDoc](https://godoc.org/github.com/twitchyliquid64/harsh?status.svg)](http://godoc.org/github.com/twitchyliquid64/harsh) [![Build Status](https://travis-ci.org/twitchyliquid64/harsh.svg?branch=master)](https://travis-ci.org/twitchyliquid64/harsh) [![Go Report Card](https://goreportcard.com/badge/github.com/twitchyliquid64/harsh)](https://goreportcard.com/report/github.com/twitchyliquid64/harsh)

Harsh is a simplified and architecture-independent code graph representation / toolchain.

```shell
go build github.com/twitchyliquid64/harsh/debugprint
./debugprint src/github.com/twitchyliquid64/harsh/compiler/translate.go
```

Planned packages, in various states are:

 * ast - represents most features of a code graph in an abstract way
  * exec - executes the code graph
  * print - prints the code graph to stdout for debugging

* compiler - translate a subset of Go to the representation used in `ast`.
 * typecheck - validate the typing of the code graph.

* visualiser - generate an image / SVG of the code graph. [planned]
* mutate - methods to mutate an existing graph or swap nodes from graphs/subgraphs (breed) [planned]
* generate - randomly generate graphs constrained by a complexity score [planned]


### Immediate TODO

#### Discrete items

 - [ ] Implement struct access
 - [ ] Implement function calls
 - [ ] Implement function types
 - [ ] Implement loops
 - [ ] Implement end-to-end tests
 - [ ] Refactor print to arbitrary output and colours
 - [ ] SVG visualiser of basic nodes

#### Longer term plan

 * Get something that is useful for local evaluation
  * Loops
  * Function calls
  * More binary operations

 * Write a visualiser
 * Write a backend for some language / javascript / assembly
  * Generalise backends to some interface that allows them to hot-plug?

 ##### If I want to do ML / random generation

 * Write a costing algorithm (min/max bounds)
 * Write a random generator
 * Write a breeder / mutator
