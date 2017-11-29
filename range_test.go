package rangetype

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRange(t *testing.T) {
	r, _ := New2("[4:2:-3") // Python style, step size -3. The "from" value is inclusive, so this results in just 4

	assert.Equal(t, r.All(), []float64{4.0})

	sum := 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 4.0)

	// Evaluate ranges

	// "Math style", with "["/"]" for inclusive and "("/")" for exclusive
	r = New("[0,1]")
	assert.Equal(t, r.String(), "[0, 1], integer range")
	r = New("(0,1]")
	assert.Equal(t, r.String(), "(0, 1], integer range")
	r = New("[0,1)")
	assert.Equal(t, r.String(), "[0, 1), integer range")
	r = New("(0,1)")
	assert.Equal(t, r.String(), "(0, 1), integer range")
	r = New("0,1")
	assert.Equal(t, r.String(), "[0, 1], integer range")

	// Ruby style
	r = New("0..1")
	assert.Equal(t, r.String(), "[0, 1], integer range")
	r = New("0..3)") // Extended syntax, with ")" for exclusive. Output: 0, 1, 2
	assert.Equal(t, r.String(), "[0, 3), integer range")
	s := ""
	r.ForEach(func(x float64) {
		s += strconv.Itoa(int(x)) + ", "
	})
	s += "EOL"
	assert.Equal(t, s, "0, 1, 2, EOL")

	// Python style
	r = New("0:3")
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	r = New("0:3") // Python syntax, exclusive
	//fmt.Println(r.String())
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	r = New("[0:3)") // Extended syntax, with ")" for exclusive. Same as above.
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})

	r = New("1:6:2") // Python syntax, with step size 2. Output: 1 3 5
	assert.Equal(t, r.All(), []float64{1.0, 3.0, 5.0})

	sum = 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 9.0)

	r = New("[4:2:-0.1]") // Python style, extended syntax, step size -0.1
	s = ""
	r.ForEach(func(x float64) {
		s += fmt.Sprintf("%.1f ", x)
	})
	s = s[:len(s)-1]
	assert.Equal(t, s, "4.0 3.9 3.8 3.7 3.6 3.5 3.4 3.3 3.2 3.1 3.0 2.9 2.8 2.7 2.6 2.5 2.4 2.3 2.2 2.1 2.0")

	assert.Equal(t, New("[0:99999999999999999999999:0.3]").Take(3), []float64{0.0, 0.3, 0.6})
	assert.Equal(t, New("(2,15) step 4").All(), []float64{6, 10, 14})
	assert.Equal(t, Slice([]float64{1.0, 2.0, 3.0, 4.0}, "1..3 step 2"), []float64{2.0, 4.0})
	assert.Equal(t, Slice([]float64{1.0, 2.0, 3.0, 4.0}, "(0,3)"), []float64{2.0, 3.0})

	assert.Equal(t, New("[3:1:-1]").All(), []float64{3.0, 2.0, 1.0})
	assert.Equal(t, New("[3:1:-1)").All(), []float64{3.0, 2.0})
	assert.Equal(t, New("3:1:-1").All(), []float64{3.0, 2.0})

	assert.Equal(t, New("1..3 step 0.5").Join(";", 2), "1.00;1.50;2.00;2.50;3.00")
	assert.Equal(t, New("1..3~ step 0.5").Join(", ", 1), "1.0, 1.5, 2.0")
	assert.Equal(t, New("1..3**2~) step 0.7").Join("->", 0), "1->2->2->3->4->4->5->6->7->7")
}

func TestByteType(t *testing.T) {
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

func TestEval(t *testing.T) {
	result, err := eval("10**2~")
	assert.Equal(t, err, nil)
	assert.Equal(t, result, 99.0)
}

//func TestIntType(t *testing.T) {
//	intType := RangeType("[0,2**16)")
//
//	assert.Equal(t, intType.Valid(42), true)
//}
