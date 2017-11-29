package rangetype

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// Unsigned integers
	U4   = New("..2**4~")
	U8   = New("..2**8~")
	U16  = New("..2**16~")
	U32  = New("..2**32~")
	U64  = New("..2**64~")
	U128 = New("..2**128~")

	// Aliases for unsigned integers
	Nibble = U4
	Char   = U8
	Byte   = U8
	Word   = U16
	Short  = U16
	Long   = U32
	Double = U32
	Quad   = U64

	// Signed integers
	I8   = New("-2**7..2**7~")
	I16  = New("-2**15..2**15~")
	I32  = New("-2**31..2**31~")
	I64  = New("-2**63..2**63~")
	I128 = New("-2**127..2**127~")
)

const (
	RANGE_MISSING       = 0
	RANGE_EXCLUDE_START = 1 << iota // Start at 1, then double
	RANGE_INCLUDE_START
	RANGE_EXCLUDE_STOP
	RANGE_INCLUDE_STOP
)

var (
	ErrRangeSyntax  = errors.New("INVALID RANGE SYNTAX")
	ErrMissingRange = errors.New("MISSING RANGE VALUES")
)

// Range can represent a number type in a programming language
// For example:
// An uint8 can be [0, 256) step 1
// A float between 0 and 1 can be [0, 1] step 0.01
type Range struct {
	rangeType int // inclusive or exclusive start and stop
	from      float64
	to        float64
	step      float64
}

// Valid is an alias for ValidFloat
func (r *Range) Valid(x float64) bool {
	return r.ValidFloat(x)
}

// ValidInt checks if the given integer is in the range
func (r *Range) ValidInt(i int) bool {
	return r.Integer() && r.ValidFloat(float64(i))
}

// ValidFloat checks if the given float is in the range,
// using half the range step size as the threshold for float equality
func (r *Range) ValidFloat(x float64) bool {
	halfStep := r.step / 2.0
	retval := r.Has(x, halfStep)
	//yn := map[bool]string{true: "YES", false: "NO"}
	//fmt.Printf("Is %v in the range %s? %v\n", x, r.String(), yn[retval])
	return retval
}

// Has checks if a given number is in the range.
// If the difference between the given float and the float in the range are
// less than the given threshold, they are counted as equal.
func (r *Range) Has(x, threshold float64) bool {
	a := min(r.from, r.to)
	b := max(r.from, r.to)
	// Check if the given number is out of the inclusive range
	if x < a || x > b {
		return false
	}
	// Check if the range has an exclusive start and if x is the same as the from value
	if ((r.rangeType & RANGE_EXCLUDE_START) != 0) && almostEqual(x, r.from, threshold) {
		return false
	}
	// Check if the range has an exclusive end and if x is the same as the to value
	if ((r.rangeType & RANGE_EXCLUDE_STOP) != 0) && almostEqual(x, r.to, threshold) {
		return false
	}
	// If the range type is an integer (step size is 1 or -1), check if x is an integer
	if r.Integer() {
		if float64(int(x)) != x {
			// Number differs when converting to an int and back to a float
			return false
		}
		// If both to and from are integers too, just check if x is between those two, and
		// then it is within range.
		if float64(int(r.from)) == r.from && float64(int(r.to)) == r.to {
			// Boundaries has already been checked, so just check between from and to
			if x > a && a < b {
				return true
			}
		}
	}

	if r.step > 0 && r.step < 1 {
		// If the step size is 0.1, extract the start value from x and check if it ends with 0.1
		translated := (x - r.from)
		fractionalPart := translated - float64(int(translated))
		if almostEqual(fractionalPart, r.step, threshold) {
			// OK, x is part of the range
			return true
		}
	}

	// TODO: Add quick checks for:
	// * step size -1
	// * step size 2
	// * step size -2

	// Now that the most important optimizations are covered,
	// fall back to actually iterating and see if x is there
	found, _ := r.Find(x, threshold)
	return found
}

// almostEqual checks if the difference between two floats are under the given threshold
func almostEqual(a, b, threshold float64) bool {
	return abs(a-b) < threshold
}

// Find searches a range for a given number
// The threshold is how close the float has to be a float in the range for it to be "equal"
// If the difference between the given float and the float in the range are
// less than the given threshold, they are counted as equal.
// The allowed difference could be 0.00001, for example. This is needed because of how floats are stored.
func (r *Range) Find(x, threshold float64) (bool, float64) {
	// Ok, loop through the range and see if is there
	found := false
	foundValue := -1.0
	r.ForEachWithBreak(func(xFromRange float64) (breakNow bool) {
		if almostEqual(x, xFromRange, threshold) {
			found = true
			foundValue = xFromRange
			breakNow = true
		}
		return
	})
	return found, foundValue
}

// Evaluate a simple expression
//
// An expression may be consists of
// floating point numbers, "**", "~" or "+".
//
// The operator presedence is undefined, and no paranthesis are supported yet.
//
// If the expression ends with "~", -1 is substracted from the result
//
// Example expression:
// > 10**2~
// 99
//
func eval(exp string) (retval float64, err error) {
	if exp == "" {
		// Return 0.0
		return retval, nil
	}
	if strings.HasSuffix(exp, "~") {
		// Start out with a value of -1
		retval = -1.0
		// Remove "~" from the expression
		exp = exp[:len(exp)-1]
	}
	if strings.Count(exp, "**") > 0 {
		elements := strings.SplitN(exp, "**", 2)
		var a, b float64
		if a, err = eval(elements[0]); err != nil {
			return retval, errors.New("INVALID VALUE: " + elements[0] + " IN " + err.Error())
		}
		if b, err = eval(elements[1]); err != nil {
			return retval, errors.New("INVALID VALUE: " + elements[1] + " IN " + err.Error())
		}
		retval += math.Pow(a, b)
		return
	} else if strings.Count(exp, "+") > 0 {
		elements := strings.SplitN(exp, "+", 2)
		var a, b float64
		if a, err = eval(elements[0]); err != nil {
			return retval, errors.New("INVALID VALUE: " + elements[0] + " IN " + err.Error())
		}
		if b, err = eval(elements[1]); err != nil {
			return retval, errors.New("INVALID VALUE: " + elements[1] + " IN " + err.Error())
		}
		retval += a + b
		return
	}
	var x float64
	if x, err = strconv.ParseFloat(exp, 64); err != nil {
		return retval, errors.New("INVALID VALUE: " + exp)
	}
	retval += x
	return
}

// New2 evaluates the given input string and returns a Range struct
func New2(rangeExpression string) (*Range, error) {
	var (
		r           = &Range{step: 1.0}
		contents    string
		err         error
		left, right string
		step        string
	)
	// If the input string contains (" step "), remove the last part
	if strings.Contains(rangeExpression, " step ") {
		elements := strings.SplitN(rangeExpression, " step ", 2)
		rangeExpression = elements[0]
		step = elements[1]
	}
	for _, c := range rangeExpression {
		switch c {
		case ' ':
			continue
		case '\t':
			continue
		case '\n':
			continue
		case '[':
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		case ']':
			r.rangeType |= RANGE_INCLUDE_STOP
			r.rangeType &= ^RANGE_EXCLUDE_STOP
		case '(':
			r.rangeType |= RANGE_EXCLUDE_START
			r.rangeType &= ^RANGE_INCLUDE_START
		case ')':
			r.rangeType |= RANGE_EXCLUDE_STOP
			r.rangeType &= ^RANGE_INCLUDE_STOP
		default:
			contents += string(c)
		}
	}
	if strings.Count(contents, "..") == 1 {
		// Ruby style range with ".."
		elements := strings.SplitN(contents, "..", 2)
		left = elements[0]
		right = elements[1]
		// Set both to inclusive, if not already set to exclusive in the switch above
		if (r.rangeType & RANGE_EXCLUDE_START) == 0 { // check if NOT set
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_EXCLUDE_STOP) == 0 { // check if NOT set
			r.rangeType |= RANGE_INCLUDE_STOP
			r.rangeType &= ^RANGE_EXCLUDE_STOP
		}
	} else if strings.Count(contents, ",") == 1 {
		elements := strings.SplitN(contents, ",", 2)
		left = elements[0]
		right = elements[1]
	} else if strings.Count(contents, ":") == 1 {
		// Python style range, as in x[0:5]
		elements := strings.SplitN(contents, ":", 2)
		left = elements[0]
		right = elements[1]
		// Set the first one to inclusive and the second one to exclusive, like in Python -
		// if not already set in the switch above.
		if (r.rangeType & RANGE_INCLUDE_START) == 0 { // check if NOT set
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 { // check if NOT set
			r.rangeType |= RANGE_EXCLUDE_STOP
			r.rangeType &= ^RANGE_INCLUDE_STOP
		}
	} else if strings.Count(contents, ":") == 2 {
		// Python style range with a step, as in x[0:5:-1]
		elements := strings.SplitN(contents, ":", 3)
		left = elements[0]
		right = elements[1]
		// Set the step, if not already set with a " step x" suffix
		if step == "" {
			step = elements[2]
		}
		// Set the first one to inclusive and the second one to exclusive, like in Python -
		// if not already set in the switch above.
		if (r.rangeType & RANGE_INCLUDE_START) == 0 { // check if NOT set
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 { // check if NOT set
			r.rangeType |= RANGE_EXCLUDE_STOP
			r.rangeType &= ^RANGE_INCLUDE_STOP
		}
	} else {
		return nil, ErrRangeSyntax
	}

	// Left side of the range expression
	if left == "" {
		// If the left side is missing, use 0
		r.from = 0.0
	} else if r.from, err = eval(left); err != nil {
		return nil, errors.New("INVALID RANGE VALUE: " + step + ", " + err.Error())
	}

	// Right side of the range expression
	if right == "" {
		return nil, ErrMissingRange
	} else if r.to, err = eval(right); err != nil {
		return nil, errors.New("INVALID RANGE VALUE: " + step + ", " + err.Error())
	}

	if step != "" {
		if r.step, err = eval(step); err != nil {
			return nil, errors.New("INVALID STEP SIZE: " + step + ", " + err.Error())
		}
	}
	return r, nil
}

// Integer checks if the range has a step of 1 or -1
func (r *Range) Integer() bool {
	return abs(r.step) == 1.0
}

// String returns the range as a string where "[" means inclusive and "(" means exclusive
func (r *Range) String() string {
	s := ""

	if (r.rangeType & RANGE_EXCLUDE_START) != 0 { // check if set
		s += "("
	} else {
		s += "["
	}

	s += fmt.Sprintf("%v, %v", r.from, r.to)

	if (r.rangeType & RANGE_EXCLUDE_STOP) != 0 { // check if set
		s += ")"
	} else {
		s += "]"
	}

	// Why "integer" instead of "step 1"?
	// The idea is to use a range to specify a number type in a future programming language.
	// By specifying a range with a step, all ints/floats/uints/bytes can be clearly defined in one single unified way.
	if r.Integer() {
		s += ", integer range"
	} else {
		s += fmt.Sprintf(", float range with step %v", r.step)
	}

	return s
}

// abs returns the absolute number
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// min returns the smallest number
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max returns the largest number
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// ForEach calls the given function for each iteration in the range
func (r *Range) ForEach(f func(float64)) {
	x := r.from
	if (r.rangeType & RANGE_INCLUDE_START) != 0 {
		if x == r.from {
			f(x)
		}
	}
	x += r.step
	if r.step > 0 {
		for x < r.to && x > r.from {
			f(x)
			x += r.step
		}
	} else if r.step < 0 {
		for x > r.to && x < r.from {
			f(x)
			x += r.step
		}
	}
	if (r.rangeType & RANGE_INCLUDE_STOP) != 0 {
		// But first check that it is within range
		f(r.to)
	}
}

// ForEachWithBreak calls the given function for each iteration in the range
// If the given function returns true, the remaining iterations are skipped
func (r *Range) ForEachWithBreak(f func(float64) bool) {
	x := r.from
	if (r.rangeType & RANGE_INCLUDE_START) != 0 {
		if x == r.from {
			if f(x) {
				// Break
				return
			}
		}
	}
	x += r.step
	if r.step > 0 {
		for x < r.to && x > r.from {
			if f(x) {
				// Break
				return
			}
			x += r.step
		}
	} else if r.step < 0 {
		for x > r.to && x < r.from {
			if f(x) {
				// Break
				return
			}
			x += r.step
		}
	}
	if (r.rangeType & RANGE_INCLUDE_STOP) != 0 {
		// But first check that it is within range
		f(r.to) // Nothing to break out of at this point
	}
}

// ForN runs the given function for the n first iterations
// If n is never reached, a smaller number of iterations will happen.
func (r *Range) ForN(n int, f func(float64)) {
	counter := 0
	x := r.from
	if (r.rangeType & RANGE_INCLUDE_START) != 0 {
		if x == r.from {
			f(x)
			counter++
			if counter >= n {
				return
			}
		}
	}
	x += r.step
	if r.step > 0 {
		for x < r.to && x > r.from {
			f(x)
			counter++
			if counter >= n {
				return
			}
			x += r.step
		}
	} else if r.step < 0 {
		for x > r.to && x < r.from {
			f(x)
			counter++
			if counter >= n {
				return
			}
			x += r.step
		}
	}
	if (r.rangeType & RANGE_INCLUDE_STOP) != 0 {
		f(r.to)
	}
}

// New is the same as New2, but panics if given an invalid input string
func New(rangeExpression string) *Range {
	r, err := New2(rangeExpression)
	if err != nil {
		panic(err)
	}
	return r
}

// All returns a slice of numbers, generated from the range
func (r *Range) All() []float64 {
	var xs []float64
	r.ForEach(func(x float64) {
		xs = append(xs, x)
	})
	return xs
}

// Slice2 can be used to slice a slice with a range
// Returns an error if the given expression is invalid
func Slice2(nums []float64, expression string) ([]float64, error) {
	var (
		selection []float64
		pos       int
	)

	r, err := New2(expression)
	if err != nil {
		return selection, err
	}

	r.ForEach(func(x float64) {
		pos = int(x)
		if pos < len(nums) {
			selection = append(selection, nums[pos])
		}
	})
	return selection, nil
}

// Slice can be used to slice a slice with a range
// Will panic if the given expresion is invalid
func Slice(nums []float64, expression string) []float64 {
	var (
		selection []float64
		pos       int
	)

	r := New(expression)
	r.ForEach(func(x float64) {
		pos = int(x)
		if pos < len(nums) {
			selection = append(selection, nums[pos])
		}
	})
	return selection
}

// Take returns a slice of n numbers, generated from the range.
// It not generate the entire slice first, but return numbers as it iterates.
func (r *Range) Take(n int) []float64 {
	var xs []float64
	r.ForN(n, func(x float64) {
		xs = append(xs, x)
	})
	return xs
}

// Join returns the output from the range as a string, where elements are separated by sep
func (r *Range) Join(sep string, digits int) string {
	numDigits := strconv.Itoa(digits) // Digits after "."

	var buf bytes.Buffer
	r.ForEach(func(x float64) {
		buf.WriteString(fmt.Sprintf("%."+numDigits+"f"+sep, x))
	})
	s := buf.String()
	lens := len(s)
	if lens > len(sep) {
		return s[:lens-len(sep)]
	}
	return s
}

// Sum adds all numbers in a range
func (r *Range) Sum() float64 {
	var sum float64
	r.ForEach(func(x float64) {
		sum += x
	})
	return sum
}

// Len returns the length of the range by iterating over it!
// May get stuck if the range is impossibly large.
func (r *Range) Len() uint64 {
	// TODO: Optimize for ranges where there is no need to actually iterate
	var counter uint64
	r.ForEach(func(_ float64) {
		counter++
	})
	return counter
}

// Bits returns the number of bits required to hold the range
func (r *Range) Bits() int {
	return int(math.Ceil(math.Log2(float64(r.Len()))))
}
