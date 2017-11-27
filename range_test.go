package rangetype

import (
	"fmt"
	"testing"
	"strconv"

	"github.com/bmizerany/assert"
)

func TestRange(t *testing.T) {
	// Evaluate ranges

	// "Math style", with "["/"]" for inclusive and "("/")" for exclusive
	r, _ := NewRange("[0,1]")
	assert.Equal(t, r.String(), "[0, 1], integer range")
	r, _ = NewRange("(0,1]")
	assert.Equal(t, r.String(), "(0, 1], integer range")
	r, _ = NewRange("[0,1)")
	assert.Equal(t, r.String(), "[0, 1), integer range")
	r, _ = NewRange("(0,1)")
	assert.Equal(t, r.String(), "(0, 1), integer range")
	r, _ = NewRange("0,1")
	assert.Equal(t, r.String(), "[0, 1], integer range")

	// Ruby style
	r, _ = NewRange("0..1")
	assert.Equal(t, r.String(), "[0, 1], integer range")
	r, _ = NewRange("0..3)") // Extended syntax, with ")" for exclusive. Output: 0, 1, 2
	assert.Equal(t, r.String(), "[0, 3), integer range")
	s := ""
	r.ForEach(func(x float64) {
		s += strconv.Itoa(int(x)) + ", "
	})
	s += "EOL"
	assert.Equal(t, s, "0, 1, 2, EOL")

	// Python style
	r, _ = NewRange("0:3")
	assert.Equal(t, r.All(), []float64{0.0, 1.0, 2.0})
	r, _ = NewRange("[0:3]") // Python syntax, exclusive
	fmt.Println(r)
	r, _ = NewRange("[0:3)") // Extended syntax, with ")" for exclusive. Same as above.
	fmt.Println(r)

	r, _ = NewRange("[1:6:2]") // Python syntax, with step size 2. Output: 1 3 5
	assert.Equal(t, r.All(), []float64{1.0, 3.0, 5.0})

	sum := 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 9.0)

	r, _ = NewRange("[4:2:-3]") // Python style, step size -3. The "from" value is inclusive, so this results in just 4
	fmt.Println(r.All())
	assert.Equal(t, r.All(), []float64{4.0})

	sum = 0.0
	r.ForEach(func(x float64) {
		sum += x
	})
	assert.Equal(t, sum, 4.0)

	r, _ = NewRange("[4:2:-0.1]") // Python style, extended syntax, step size -0.1
	s = ""
	r.ForEach(func(x float64) {
		s += fmt.Sprintf("%.1f ", x)
	})
	s = s[:len(s)-1]
	assert.Equal(t, s, "4.0 3.9 3.8 3.7 3.6 3.5 3.4 3.3 3.2 3.1 3.0 2.9 2.8 2.7 2.6 2.5 2.4 2.3 2.2 2.1 2.0")

	assert.Equal(t, MustRange("[0:99999999999999999999999:0.3]").Take(3), []float64{0.0, 0.3, 0.6})
	assert.Equal(t, MustRange("(2,15) step 4").All(), []float64{6, 10, 14})
	assert.Equal(t, MustSlice([]float64{1.0, 2.0, 3.0, 4.0}, "1..3 step 2"), []float64{2.0, 4.0})
	assert.Equal(t, MustSlice([]float64{1.0, 2.0, 3.0, 4.0}, "(0,3)"), []float64{2.0, 3.0})
}
