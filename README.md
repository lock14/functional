# Functional

[![Go Reference](https://pkg.go.dev/badge/github.com/lock14/functional.svg)](https://pkg.go.dev/github.com/lock14/functional)
[![Go CI](https://github.com/lock14/functional/actions/workflows/go.yml/badge.svg)](https://github.com/lock14/functional/actions/workflows/go.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A type-safe functional programming library for modern Go (1.23+), providing eager slice transformations, lazy standard library iterators, concurrent streaming channels, and composable predicates.

## Features

- **`slice`**: Eager operations on Go slices (`Map`, `Filter`, `FoldLeft`, `FoldRight`, `Reduce`, `Zip`, `UnZip`, `Partition`, `Sum`, `Concat`).
- **`iterator`**: Lazy sequence transformations built directly on Go 1.23+ `iter.Seq` and `iter.Seq2` (`Map`, `Filter`, `Zip`, `Distinct`, `Sorted`, `Range`, `Limit`, `Skip`, `AllMatch`, `AnyMatch`).
- **`channel`**: Concurrent, streaming pipelines over Go channels with full `context.Context` cancellation, worker pools (`ParallelMap`, `ParallelFilter`, `ParallelFlatMap`), and error channel splitting (`MapWithErr`, `FilterWithErr`).
- **`predicate`**: Generic higher-order predicate helpers and reflection-based nil checks (`IsNil`, `NotNil`, `Not`, `True`, `False`).

---

## Installation

```bash
go get github.com/lock14/functional
```

---

## Quickstart

### 1. Eager Slice Transformations (`slice`)

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/lock14/functional/slice"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6}

	// Filter even numbers
	evens := slice.Filter(numbers, func(n int) bool { return n%2 == 0 })

	// Map to strings
	strings := slice.Map(evens, strconv.Itoa)

	// Sum numbers
	total := slice.Sum(numbers)

	fmt.Println("Evens:", strings) // ["2", "4", "6"]
	fmt.Println("Sum:", total)     // 21
}
```

### 2. Lazy Iterators (`iterator`)

Leverage standard Go 1.23+ range-over-func iterators (`iter.Seq`):

```go
package main

import (
	"fmt"

	"github.com/lock14/functional/iterator"
)

func main() {
	// Generate an infinite stream, take the first 5 even squares, and collect
	naturals := iterator.Range(1, 100)
	evens := iterator.Filter(naturals, func(n int) bool { return n%2 == 0 })
	squares := iterator.Map(evens, func(n int) int { return n * n })
	first5 := iterator.Limit(squares, 5)

	for val := range first5 {
		fmt.Println(val) // 4, 16, 36, 64, 100
	}
}
```

### 3. Concurrent Pipelines (`channel`)

```go
package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lock14/functional/channel"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := channel.Of(1, 2, 3, 4, 5, 6, 7, 8)

	// Process with 4 parallel worker goroutines
	results := channel.ParallelMap(ctx, 4, input, func(n int) string {
		return "Item-" + strconv.Itoa(n)
	})

	for item := range results {
		fmt.Println(item)
	}
}
```

---

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
