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

func _() {
	xs := []int{1, 2, 3}
	ys := make([]int, len(xs))
	copy(ys, xs)
	for i := 0; i < 10; i++ {
		ys = append(ys, i)
	}
	fmt.Println(ys)
}

func _() {
	xs := []string{"a", "b", "c"}
	ys := make([]int, len(xs))
	zs := make([]int, len(xs)) // want ".+"
	for i, x := range xs {
		ys[i] = len(x)
		zs = append(zs, len(x))
	}
	for i := 0; i < 10; i++ {
		ys = append(ys, i)
		zs = append(zs, i)
	}
	fmt.Println(ys)
}

func _() {
	xs := []string{"a", "b", "c"}
	ys := make([]int, len(xs))

	ys = ys[:0]
	for _, x := range xs {
		ys = append(ys, len(x))
	}
	fmt.Println(ys)
}
