package testdata

import "fmt"

func _() {
	a := make([]int, 0)
	b := make([]int, 10) // want ".+"
	c := make([]int, 5, 10)
	var d []int
	f := make([]int, 10)

	for i := 0; i < 10; i++ {
		a = append(a, i)
		b = append(b, i)
		c = append(c, i)
		d = append(d, i)
		f[i] = i
	}
}

func _() {
	a := make([]int, 10) // want ".+"
	for i := 0; i < 10; i++ {
		a = append(a, i)
		a := make([]int, 10)
		for j := 0; j < 10; j++ {
			a[j] = j
		}
	}
	fmt.Println(a)
}
