package main

import (
	"fmt"
)

func add(a, b []int32, p int) []int32 {
	ans := make([]int32, 0)
	var mem = 0
	var l1 = len(a)
	var l2 = len(b)
	var req int
	if l1 < l2 {
		req = l1
	} else {
		req = l2
	}
	for i := 0; i < req; i++ {
		elem1 := (int)(a[i])
		elem2 := (int)(b[i])
		if elem1+elem2+mem >= p {
			ans = append(ans, (int32)((elem1+elem2+mem)%p))
			mem = 1
		} else {
			ans = append(ans, (int32)((elem1+elem2+mem)%p))
			mem = 0
		}
	}
	if mem == 1 && len(a) == len(b) {
		ans = append(ans, 1)
	} else if len(a) != len(b) {
		var bigger []int32
		var smaller []int32
		if len(a) > len(b) {
			bigger = a
			smaller = b
		} else {
			bigger = b
			smaller = a
		}
		for i := 0; i < len(bigger)-len(smaller); i++ {
			if (int32)(mem)+bigger[i+len(smaller)] == (int32)(p) {
				mem = 1
				ans = append(ans, 0)
			} else {
				ans = append(ans, (int32)(mem)+bigger[i+len(smaller)])
				mem = 0
			}
		}
		if mem == 1 {
			ans = append(ans, 1)
		}
	}
	return ans
}

func main() {
	A := []int32{1, 2}
	B := []int32{2}
	var C = add(A, B, 3)
	for _, val := range C {
		fmt.Printf("%d", val)
	}
}
