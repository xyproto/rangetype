package range2

import (
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

	RANGE_INCLUDE_INCLUDE = RANGE_INCLUDE_START | RANGE_INCLUDE_STOP
	RANGE_INCLUDE_EXCLUDE = RANGE_INCLUDE_START | RANGE_EXCLUDE_STOP
	RANGE_EXCLUDE_INCLUDE = RANGE_EXCLUDE_START | RANGE_INCLUDE_STOP
	RANGE_EXCLUDE_EXCLUDE = RANGE_EXCLUDE_START | RANGE_EXCLUDE_STOP
)

var (
	ErrRangeSyntax  = errors.New("INVALID RANGE SYNTAX")
	ErrMissingRange = errors.New("MISSING RANGE VALUES")
)

type Range struct {
	rangeType uint8 // inclusive or exclusive start and stop
	from      float64
	to        float64
	step      float64
}

// IsInteger checks if the range has a step of 1
func (r *Range) IsInteger() bool {
	return r.step == 1.0
}

// String returns the range as a string where "[" means inclusive and "(" means exclusive
func (r *Range) String() string {
	s := ""
	switch r.rangeType {
	case RANGE_EXCLUDE_EXCLUDE, RANGE_EXCLUDE_INCLUDE:
		s += "("
	default:
		//case RANGE_INCLUDE_EXCLUDE, RANGE_INCLUDE_INCLUDE:
		s += "["
	}

	s += fmt.Sprintf("%v %v", r.from, r.to)

	switch r.rangeType {
	case RANGE_INCLUDE_EXCLUDE, RANGE_EXCLUDE_EXCLUDE:
		s += ")"
	default:
		//case RANGE_INCLUDE_INCLUDE, RANGE_EXCLUDE_INCLUDE:
		s += "]"
	}

	// Why "integer" instead of "step 1"?
	// The idea is to use a range to specify a number type in a future programming language.
	// By specifying a range with a step, all ints/floats/uints/bytes can be clearly defined in one single unified way.
	if r.IsInteger() {
		s += ", integer"
	} else {
		s += fmt.Sprintf(", float with step %v", r.step)
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
		diff := x - r.to
		// If the remaining distance to the goal is smaller than half a step, use the goal value
		if abs(diff) < abs(r.step)/2.0 {
			f(r.to)
		}
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
		diff := x - r.to
		// If the remaining distance to the goal is smaller than half a step, use the goal value
		if abs(diff) < abs(r.step)/2.0 {
			f(r.to)
		}
	}
}

// NewRange evaluates the given input string and returns a Range struct
func NewRange(inputString string) (*Range, error) {
	var (
		r           = &Range{step: 1.0}
		contents    string
		err         error
		left, right string
	)
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
		case ']':
			r.rangeType |= RANGE_INCLUDE_STOP
		case '(':
			r.rangeType |= RANGE_EXCLUDE_START
		case ')':
			r.rangeType |= RANGE_EXCLUDE_STOP
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
		}
		if (r.rangeType & RANGE_EXCLUDE_STOP) == 0 {
			r.rangeType |= RANGE_INCLUDE_STOP
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
		// Set the first one to inclusive and the second one to exclusive, if not already set otherwise
		if (r.rangeType & RANGE_EXCLUDE_START) == 0 {
			r.rangeType |= RANGE_INCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 {
			r.rangeType |= RANGE_EXCLUDE_STOP
		}
	} else if strings.Count(contents, ":") == 2 {
		// Python style range with a step, as in x[0:5:-1]
		elements := strings.SplitN(contents, ":", 3)
		left = elements[0]
		right = elements[1]
		step := elements[2]
		if r.step, err = strconv.ParseFloat(step, 64); err != nil {
			return nil, errors.New("INVALID STEP SIZE: " + step)
		}
		// Set the first one to inclusive and the second one to exclusive, if not already set otherwise
		if (r.rangeType & RANGE_EXCLUDE_START) == 0 {
			r.rangeType |= RANGE_INCLUDE_START
		}
		if (r.rangeType & RANGE_INCLUDE_STOP) == 0 {
			r.rangeType |= RANGE_EXCLUDE_STOP
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

	return r, nil
}

// MustRange ris the same as NewRange, but panics if given an invalid input string
func MustRange(inputString string) *Range {
	r, err := NewRange(inputString)
	if err != nil {
		panic(err)
	}
	return r
}

// Slice returns a slice of numbers, generated from the range
func (r *Range) Slice() []float64 {
	var xs []float64
	r.ForEach(func(x float64) {
		xs = append(xs, x)
	})
	return xs
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

func (r *Range) Sum() float64 {
	var sum float64
	r.ForEach(func(x float64) {
		sum += x
	})
	return sum
}
