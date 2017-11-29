package main

import (
	"fmt"

	r "github.com/xyproto/rangetype"
)

func main() {

	// A byte type can be defined as for example "..2**8~"
	// or just "0..255"

	fmt.Println(" --- Byte ---")

	fmt.Println("42 is a valid byte value?", r.Byte.Valid(42))
	fmt.Println("256 is a valid byte value?", r.Byte.Valid(256))
	fmt.Println("-0.2 is a valid byte value?", r.Byte.Valid(-0.2))

	fmt.Printf("Storage required for a byte: %d-bit int\n", r.Byte.Bits())

	fmt.Println(" --- SmallInt ---")

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
	fmt.Println("All possible values in SmallInt:\n" + SmallInt.Join(",", 0))
}
