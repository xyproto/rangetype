# Rangetype

A mini-language for ranges

# Example


```go
package main

import (
	"fmt"

	r "github.com/xyproto/rangetype"
)

func main() {
	// Outputs 1 to 10, with 0 digits after "."
	fmt.Println(r.MustRange("1..10").Join(", ", 0))

	// Outputs 2 and 4
	for _, x := range r.MustSlice([]float64{1.0, 2.0, 3.0, 4.0}, "1..3 step 2") {
		fmt.Print(x, " ")
	}
	fmt.Println()

	// Also outputs 2 and 4
	r.MustRange("(0:6:2)").ForEach(func(x float64) {
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
rangetype.MustRange("1..10").ForEach(func(x float64) {
  fmt.Println(int(x))
})
```

Collecting integers to a comma separated string can be done with `Join`:

    rangetype.MustRange("1..10").Join(", ", 0)

Or for floats, with 2 digits after the period, separated by semicolons:

    rangetype.MustRange("1..3 step 0.5").Join(";", 2)

---

The functions starting with `Must` returns no error and can panic, while the corresponding functions may return an error and does not panic.

There are more examples in the `*_test.go` file.

---

License: MIT
Version: 0.1
