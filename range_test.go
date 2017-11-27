package range2

import (
	"fmt"
	"testing"
)

func TestRange(t *testing.T) {
	// Evaluate ranges

	// "Math style", with "["/"]" for inclusive and "("/")" for exclusive
	r, _ := NewRange("[0,1]")
	fmt.Println(r)
	r, _ = NewRange("(0,1]")
	fmt.Println(r)
	r, _ = NewRange("[0,1)")
	fmt.Println(r)
	r, _ = NewRange("(0,1)")
	fmt.Println(r)
	r, _ = NewRange("0,1")
	fmt.Println(r)

	// Ruby style
	r, _ = NewRange("0..1")
	fmt.Println(r)
	r, _ = NewRange("0..3)") // Extended syntax, with ")" for exclusive. Output: 0, 1, 2
	fmt.Println(r)
	r.ForEach(func(x float64) {
		fmt.Println(x)
	})

	// Python style
	r, _ = NewRange("0:3")
	fmt.Println(r)
	r, _ = NewRange("[0:3]") // Python syntax, exclusive
	fmt.Println(r)
	r, _ = NewRange("[0:3)") // Extended syntax, with ")" for exclusive. Same as above.
	fmt.Println(r)
	r, _ = NewRange("[1:4:2]") // Python syntax, with step size 2. Output: 1 3
	fmt.Println(r)
	r.ForEach(func(x float64) {
		fmt.Println(x)
	})
	r, _ = NewRange("[4:2:-3]") // Python style, step size -3
	fmt.Println(r)
	r.ForEach(func(x float64) {
		fmt.Println(x)
	})
	r, _ = NewRange("[4:2:-0.1]") // Python style, extended syntax, step size -0.1
	fmt.Println(r)
	r.ForEach(func(x float64) {
		fmt.Printf("%.1f\n", x)
	})
	fmt.Println(r.Slice())

	fmt.Println(MustRange("[0:99999999999999999999999:0.3]").Take(3)) // 0, 0.3, 0.6
}
