# Rangetype

[![Build Status](https://travis-ci.org/xyproto/rangetype.svg?branch=master)](https://travis-ci.org/xyproto/rangetype) [![GoDoc](https://godoc.org/github.com/xyproto/rangetype?status.svg)](http://godoc.org/github.com/xyproto/rangetype) [![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/xyproto/rangetype/master/LICENSE) [![Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg?style=flat)](http://goreportcard.com/report/xyproto/rangetype)

A mini-language for defining numeric types by defining ranges.

The idea is to provide a DSL for defining and validating numeric types for implementations of programming languages.

It can also be used for iterating over ranges, generating lists of numbers or slicing a given slice (like slices in Python).

## Sample syntax

An int with a range from 1 to 3 that includes both 1, 2 and 3:

`1,3`

A float with a range from 1 to 3 that includes 1.0, 1.1 etc up to 3.0:

`1,3 step 0.1`

Inclusivity is specified with square brackets:

`[1,3]`

Exclusivity is specified with parenthesis:

`(1,3)`

Ranges inspired by Ruby also work:

`1..3`

These are inclusive, unless parenthesis are used.

Ruby-style range which will exclude `1` and `3` and only keep `2`:

`(1..3)`

Python style ranges are also supported, where the start value is inclusive and the end value is exclusive:

`1:3`

Adding square brackets makes the range inclusive:

`[1:3]`

Brackets and parenthesis does not have to be balanced. This works too:

`1:3]`

Adding an iteration step is also possible:

`1..5 step 2`

This is a range with the numbers `1`, `3` and `5`.

The Python-style syntax also supports steps:

`[3:1:-1]`

This is `3`, `2`, `1`.

Steps does not have to be integers:

`[3:1:-0.1]`

This steps from 3 (inclusive) down to 1 (inclusive) in step sizes of 0.1.

## ForEach

Looping over a range can be done by providing a function that takes a `float64`:

```go
r.New("1..10").ForEach(func(x float64) {
  fmt.Println(int(x))
})
```

## Join

Collecting integers to a comma separated string can be done with `Join`:

```go
r.New("1..10").Join(", ", 0)
```

Or for floats, with 2 digits after the period, separated by semicolons:

```go
r.New("1..3 step 0.5").Join(";", 2)
```

## Features and limitations

* Can handle very large ranges without storing the actual numbers in the ranges, but iterating over them may be slow.
* Only `**` and `~` are supported for manipulating numbers in the range expressions. It can not handle addition, subtraction, parenthesis etc. It's not a general language, it's only a DSL for expressing ranges of integers or floating point numbers, with an optional step size.

## Syntax

Expressions can optionally start with:

* `[` for including the first value in the range, or
* `(` for excluding the first value in the range

And can end with:

* `]` for including the last value in the range, or
* `)` for excluding the last value in the range

Numbers can be suffixed by:

* `~` for subtracting 1 from the preceding number. Any number of `~` is possible.

The ranges can be Python-style:

`[0:10]`

Python-style with a step:

`[1:20:-1]`

Ruby-style:

`1..10`

Ruby-style with a step:

`1..10 step 2`

Math-style:

`[1,5)`

Math-style with a step:

`(5,1] step -0.1`

Or with powers. Here's an expression for specifying the range for a 16-bit unsigned integer:

`..2**16~`

This can be used for validating if a number fits the type:

```go
IntType := r.New("..2**16~")     // from 0 up to and including 65536-1
IntType.Valid(42)                // true
```

## More examples

### Defining a SmallInt type and checking if a given number is valid

```go
package main

import (
	"fmt"

	r "github.com/xyproto/rangetype"
)

func main() {
	// Define a new type that can hold numbers from 0 up to and including 99
	SmallInt := r.New("0..99")

	// Another way to define a number type from 0 up to and excluding 100
	//SmallInt := New("[0,100)")

	// Another way to define a number type from 0 up to and excluding 100
	//SmallInt := New("10**2~")

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

## Example 2

### Slicing slices and looping with the ForEach method

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

There are more examples in the `range_test.go` file.

## Notes

* `New2` and `Slice2` will return both a value and an error (if the expression failed to evaluate) and are the recommended functions to use.
* `New` and `Slice` are fine to use for expressions that are already known to evaluate, but may panic if there are errors in the given expression.


## General info

* License: MIT
* Version: 0.1
* Author: Alexander F Rødseth &lt;xyproto@archlinux.org&gt;
