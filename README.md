# Rangetype

[![Build Status](https://travis-ci.org/xyproto/rangetype.svg?branch=master)](https://travis-ci.org/xyproto/rangetype) [![GoDoc](https://godoc.org/github.com/xyproto/rangetype?status.svg)](http://godoc.org/github.com/xyproto/rangetype) [![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/xyproto/rangetype/master/LICENSE) [![Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg?style=flat)](http://goreportcard.com/report/xyproto/rangetype)

A mini-language for defining numeric types by defining ranges.

The idea is to provide a DSL for defining and validating numeric types for implementations of programming languages.

It can also be used for iterating over ranges, generating lists of numbers or slicing a given slice (like slices in Python).

# Example 1 - defining a SmallInt type and checking if a given number is valid

```go
package main

import (
	"fmt"

	r "github.com/xyproto/rangetype"
)

func main() {
	// Define a new type that can hold numbers from 0 up to and including 99
	SmallInt := r.New("0..99")

	// Another way to define a number type from 0 up to and including 99
	//SmallInt := New("10**2~")

	// Another way to define a type from 0 up to and including 99
	//SmallInt := New("[0,100)")

	// Is 42 a valid SmallInt?
	fmt.Println("0 is a valid SmallInt value:", SmallInt.Valid(0))
	fmt.Println("2 is a valid SmallInt value:", SmallInt.Valid(2))
	fmt.Println("-1 is a valid SmallInt value:", SmallInt.Valid(-1))
	fmt.Println("99 is a valid SmallInt value:", SmallInt.Valid(99))
	fmt.Println("100 is a valid SmallInt value:", SmallInt.Valid(100))

	// How many integers are there room for?
	fmt.Printf("SmallInt can hold %d different numbers.\n", SmallInt.Len())
	fmt.Printf("Storage required for SmallInt: a %d-bit int\n", SmallInt.Bits())

	// All possible SmallInt values, comma separated:
	fmt.Println("All possible values for SmallInt:\n" + SmallInt.Join(",", 0))
}
```

# Example 2 - slicing slices and looping with the ForEach method

```go
package main

import (
	"fmt"

	r "github.com/xyproto/rangetype"
)

func main() {
	// Outputs 1 to 10, with 0 digits after "."
	fmt.Println(r.New("1..10").Join(", ", 0))

	// Outputs 2 and 4
	for _, x := range r.Slice([]float64{1.0, 2.0, 3.0, 4.0}, "1..3 step 2") {
		fmt.Print(x, " ")
	}
	fmt.Println()

	// Also outputs 2 and 4
	r.New("(0:6:2)").ForEach(func(x float64) {
		fmt.Print(x, " ")
	})
	fmt.Println()
}
```

# Syntax

Expressions can optionally start with:

* `[` for including the first value in the range, or
* `(` for excluding the first value in the range

And can end with:

* `]` for including the last value in the range, or
* `)` for excluding the last value in the range
* `~` for subtracting 1 from the preceeding number

An example of a range from 1 to 3 that includes both 1, 2 and 3 is:

`[1,3]`

The default is inclusive, so the above is the same as just:

`1,3`

Ranges inspired by Ruby also work:

`1..3`

These are also inclusive, unless a `(`, `)` or both are used. This will exclude `1` and `3` and keep only `2`:

`(1..3)`

Python style ranges are also supported, where the start value is inclusive and the end value is exclusive:

`1:3`

This is a range of only `1` and `2`, the `3` is excluded. However, this will include both `1` and `3`:

`[1:3]`

Adding an iteration step is also supported:

`1..5 step 2`

This is a range with the numbers `1`, `3` and `5`.

Using the Python-style syntax also supports steps:

`[3:1:-1]`

This is `3`, `2`, `1`.

Steps does not have to be integers:

`[3:1:-0.1]`

This steps from 3 (inclusive) down to 1 (inclusive) in step sizes of 0.1.

Looping over a range can be done by providing a function that takes a `float64`:

```
r.New("1..10").ForEach(func(x float64) {
  fmt.Println(int(x))
})
```

Collecting integers to a comma separated string can be done with `Join`:

    r.New("1..10").Join(", ", 0)

Or for floats, with 2 digits after the period, separated by semicolons:

    r.New("1..3 step 0.5").Join(";", 2)

---

The functions `New` and `Slice` will just return a value, and may panic, while `New2` and `Slice2` will return an error value as well, and not panic.

There are more examples in the `range_test.go` file.

---

License: MIT
Version: 0.1
