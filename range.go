package rangetype

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// NewRange evaluates the given input string and returns a Range struct
func NewRange(inputString string) (*Range, error) {
	var (
		r           = &Range{step: 1.0}
		contents    string
		err         error
		left, right string
		step        string
	)
	// If the input string contains (" step "), remove the last part
	if strings.Contains(inputString, " step ") {
		elements := strings.SplitN(inputString, " step ", 2)
		inputString = elements[0]
		step = elements[1]
	}
	for _, c := range inputString {
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
		if (r.rangeType & RANGE_EXCLUDE_START) == 0 {
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_EXCLUDE_STOP) == 0 {
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
		if (r.rangeType & RANGE_INCLUDE_START) == 0 { // no inclusive start defined
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 { // no inclusive stop defined
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
		if (r.rangeType & RANGE_INCLUDE_START) == 0 { // no inclusive start defined
			r.rangeType |= RANGE_INCLUDE_START
			r.rangeType &= ^RANGE_EXCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 { // no inclusive stop defined
			r.rangeType |= RANGE_EXCLUDE_STOP
			r.rangeType &= ^RANGE_INCLUDE_STOP
		}
	} else {
		return nil, ErrRangeSyntax
	}

	if left == "" || right == "" {
		return nil, ErrMissingRange
	}

	if r.from, err = strconv.ParseFloat(left, 64); err != nil {
		return nil, errors.New("INVALID RANGE VALUE: " + left)
	}

	if r.to, err = strconv.ParseFloat(right, 64); err != nil {
		return nil, errors.New("INVALID RANGE VALUE: " + right)
	}

	if step != "" {
		if r.step, err = strconv.ParseFloat(step, 64); err != nil {
			return nil, errors.New("INVALID STEP SIZE: " + step)
		}
	}

	return r, nil
}

// IsInteger checks if the range has a step of 1
func (r *Range) IsInteger() bool {
	return r.step == 1.0
}

// String returns the range as a string where "[" means inclusive and "(" means exclusive
func (r *Range) String() string {
	s := ""

	if (r.rangeType & RANGE_EXCLUDE_START) != 0 {
		s += "("
	} else {
		s += "["
	}

	s += fmt.Sprintf("%v, %v", r.from, r.to)

	if (r.rangeType & RANGE_EXCLUDE_STOP) != 0 {
		s += ")"
	} else {
		s += "]"
	}

	// Why "integer" instead of "step 1"?
	// The idea is to use a range to specify a number type in a future programming language.
	// By specifying a range with a step, all ints/floats/uints/bytes can be clearly defined in one single unified way.
	if r.IsInteger() {
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

// MustRange ris the same as NewRange, but panics if given an invalid input string
func MustRange(inputString string) *Range {
	r, err := NewRange(inputString)
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

// Slice can be used to slice a slice with a range
func Slice(nums []float64, expression string) ([]float64, error) {
	r, err := NewRange(expression)
	if err != nil {
		return []float64{}, err
	}
	var selection []float64
	var pos int
	r.ForEach(func(x float64) {
		pos = int(x)
		if pos < len(nums) {
			selection = append(selection, nums[pos])
		}
	})
	return selection, nil
}

// MustSlice can be used to slice a slice with a range
// Will panic if the given expresion is invalid
func MustSlice(nums []float64, expression string) []float64 {
	r := MustRange(expression)
	var selection []float64
	var pos int
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
