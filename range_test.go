package rangetype

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bmizerany/assert"
)

func TestSimpleRange(t *testing.T) {
	r := New("[1,3]")
	assert.Equal(t, r.All(), []float64{1.0, 2.0, 3.0})
}

func TestPythonStyleRanges(t *testing.T) {
	r := New("0:3")
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	// The "from" value is inclusive, and the -3 step leads nowhere,
	// so the only number in the range is 4.
	r = New("[4:2:-3")
	assert.Equal(t, r.All(), []float64{4})

	r = New("0:3") // Python syntax, exclusive
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	r = New("[0:3)") // Extended syntax, with ")" for exclusive. Same as above.
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	r = New("1:6:2") // Python syntax, with step size 2. Output: 1 3 5
	assert.Equal(t, r.All(), []float64{1.0, 3.0, 5.0})
	sum := 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 9.0)

	r = New("[4:2:-0.1]") // Python style, extended syntax, step size -0.1
	s := ""
	r.ForEach(func(x float64) {
		s += fmt.Sprintf("%.1f ", x)
	})
	s = s[:len(s)-1]
	assert.Equal(t, s, "4.0 3.9 3.8 3.7 3.6 3.5 3.4 3.3 3.2 3.1 3.0 2.9 2.8 2.7 2.6 2.5 2.4 2.3 2.2 2.1 2.0")

	// Exclusive end
	assert.Equal(t, New("0:3").All(), []float64{0.0, 1.0, 2.0})

	// Inclusive end
	assert.Equal(t, New("0:3]").All(), []float64{0.0, 1.0, 2.0, 3.0})
}

func TestSum(t *testing.T) {
	r, _ := New2("[4,5]")
	sum := 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 9.0)
	assert.Equal(t, sum, r.Sum())
}

func TestDescription(t *testing.T) {
	// "Math style", with "["/"]" for inclusive and "("/")" for exclusive
	r := New("[0,1]")
	assert.Equal(t, r.String(), "[0, 1], integer range")
	r = New("(0,1]")
	assert.Equal(t, r.String(), "(0, 1], integer range")
	r = New("[0,1)")
	assert.Equal(t, r.String(), "[0, 1), integer range")
	r = New("(0,1)")
	assert.Equal(t, r.String(), "(0, 1), integer range")
	r = New("0,1")
	assert.Equal(t, r.String(), "[0, 1], integer range")
}

func TestRubyStyleRanges(t *testing.T) {
	// Ruby style
	r := New("0..1")
	assert.Equal(t, r.String(), "[0, 1], integer range")

	r = New("0..3)") // Extended syntax, with ")" for exclusive. Output: 0, 1, 2
	assert.Equal(t, r.String(), "[0, 3), integer range")

	s := ""
	r.ForEach(func(x float64) {
		s += strconv.Itoa(int(x)) + ", "
	})
	s += "EOL"
	assert.Equal(t, s, "0, 1, 2, EOL")

	assert.Equal(t, New("1..3 step 0.5").Join(";", 2), "1.00;1.50;2.00;2.50;3.00")
	assert.Equal(t, New("1..3~ step 0.5").Join(", ", 1), "1.0, 1.5, 2.0")
	assert.Equal(t, New("1..3**2~) step 0.7").Join("->", 0), "1->2->2->3->4->4->5->6->7->7")
}

func TestSlices(t *testing.T) {
	assert.Equal(t, New("[0:99999999999999999999999:0.3]").Take(3), []float64{0.0, 0.3, 0.6})
	assert.Equal(t, New("(2,15) step 4").All(), []float64{6, 10, 14})
	assert.Equal(t, Slice([]float64{1.0, 2.0, 3.0, 4.0}, "1..3 step 2"), []float64{2.0, 4.0})
	assert.Equal(t, Slice([]float64{1.0, 2.0, 3.0, 4.0}, "(0,3)"), []float64{2.0, 3.0})
	assert.Equal(t, New("[3:1:-1]").All(), []float64{3.0, 2.0, 1.0})
	assert.Equal(t, New("[3:1:-1)").All(), []float64{3.0, 2.0})
	assert.Equal(t, New("3:1:-1").All(), []float64{3.0, 2.0})
}

func TestByteTypes(t *testing.T) {
	// Various ways to define a byte type
	byteTypes := []*Range{
		Byte,
		Char,
		U8,
		New("0..255"),
		New("[0,255]"),
		New("[0,2**8)"),
		New("..2**8)"),
		New("..2**8~"),
		New("[0,256)"),
	}

	x := 2.000000001

	for _, byteType := range byteTypes {
		assert.Equal(t, byteType.Valid(255), true)
		assert.Equal(t, byteType.Valid(42), true)
		assert.Equal(t, byteType.Valid(254), true)
		assert.Equal(t, byteType.Valid(256), false)
		assert.Equal(t, byteType.Valid(0), true)
		assert.Equal(t, byteType.Valid(-1), false)
		assert.Equal(t, byteType.Valid(0.1), false)
		assert.Equal(t, byteType.Valid(-0.1), false)
		assert.Equal(t, byteType.Valid(x), false)
		assert.Equal(t, byteType.ValidInt(int(x)), true)
		assert.Equal(t, byteType.Bits(), 8)
	}
}

func TestMinusOneOperator(t *testing.T) {
	result, err := eval("10**2~", false)
	assert.Equal(t, err, nil)
	assert.Equal(t, result, 99.0)

	r, err := New2("[1~..4~~**2)") // From 0 up to, but not including, 4
	assert.Equal(t, err, nil)
	assert.Equal(t, r.Len(), uint(4))

	IntType := New("..2**16~") // from 0 up to and including 65536-1
	assert.Equal(t, IntType.Valid(42), true)
}

func TestFloatBits(t *testing.T) {
	SmallFloat := New("[0:1.0:0.1)")

	// 10 possible values, should fit in a 4 bit number
	assert.Equal(t, SmallFloat.Bits(), 4)

	assert.Equal(t, SmallFloat.Valid(0.0), true)
	assert.Equal(t, SmallFloat.Valid(0.1), true)
	assert.Equal(t, SmallFloat.Valid(0.05), false)
	assert.Equal(t, SmallFloat.Valid(-0.1), false)
	assert.Equal(t, SmallFloat.Valid(1.0), false)
	assert.Equal(t, SmallFloat.Valid(1.05), false)
	assert.Equal(t, SmallFloat.Valid(1.1), false)
	assert.Equal(t, SmallFloat.Valid(0.9), true)

	// 16 possible values, should still fit in a 4 bit number
	SmallFloat2 := New("[0:1.6:0.1)") // same as [0..1.6) step 0.1
	assert.Equal(t, SmallFloat2.Bits(), 4)

	// 17 possible values, should still fit in a 5 bit number
	SmallFloat3 := New("[0..1.7) step 0.1") // same as [0:1.7:0.1]
	assert.Equal(t, SmallFloat3.Bits(), 5)
}

func TestSmallInt(t *testing.T) {
	Integer8 := New("-2**7 .. 2**7~")
	assert.Equal(t, Integer8.Valid(100), true)
}

func TestRsplit(t *testing.T) {
	s := "abc (123 (abc) cheese) asdf"
	_, right := rsplit(s, ")")
	assert.Equal(t, right, " asdf")
}

func TestAda(t *testing.T) {
	Integer8 := NewAda("-(2**7) .. (2**7)-1")
	assert.Equal(t, Integer8.Valid(100), true)

	TinyType := NewAda("-5 .. 10")
	l := TinyType.All()
	assert.Equal(t, l[0], -5.0)
	assert.Equal(t, l[1], -4.0)
	assert.Equal(t, l[5], 0.0)
	assert.Equal(t, l[14], 9.0)
	assert.Equal(t, l[len(l)-1], 10.0)

	Short := NewAda("-128 .. +127")
	l = Short.All()
	assert.Equal(t, l[0], -128.0)
	assert.Equal(t, l[len(l)-1], 127.0)

	Integer := NewAda("0 .. Integer'Last")
	assert.Equal(t, Integer.Len64(), float64(MaxInt))
}
