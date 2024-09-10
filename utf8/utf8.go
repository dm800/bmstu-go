package main

import (
	"fmt"
	"math"
)

func encode(utf32 []rune) []byte {
	var ans = make([]byte, 0)
	var c = 0
	for _, val := range utf32 {
		num := (int64)(val)
		if num <= 128 {
			ans = append(ans, (byte)(num))
			c += 1
		} else if num <= 2048 {
			k1 := (byte)(num / 64)
			k1 += 192
			ans = append(ans, k1)
			k2 := (byte)(num % 64)
			k2 += 128
			ans = append(ans, k2)
		} else if num <= 65536 {
			k1 := (byte)(num / 4096)
			k1 += 224
			ans = append(ans, k1)
			k2 := (byte)((num % 4096) / 64)
			k2 += 128
			ans = append(ans, k2)
			k3 := (byte)(num % 64)
			k3 += 128
			ans = append(ans, k3)
		} else {
			k1 := (byte)(num / 262144)
			k1 += 240
			ans = append(ans, k1)
			k2 := (byte)((num % 262144) / 4096)
			k2 += 128
			ans = append(ans, k2)
			k3 := (byte)((num % 4096) / 64)
			k3 += 128
			ans = append(ans, k3)
			k4 := (byte)(num % 64)
			k4 += 128
			ans = append(ans, k4)
		}
	}
	return ans
}

func decode(utf8 []byte) []rune {
	var ans = make([]rune, 0)
	var su rune
	su = 0
	count := 0
	for _, val := range utf8 {
		if count == 0 {
			if val < 128 {
				ans = append(ans, (rune)(val))
			} else if val < 224 {
				count += 1
				su += (64 * (rune)(val-192))
			} else if val < 240 {
				count += 2
				su += (4096 * (rune)(val-224))
			} else if val < 248 {
				count += 3
				su += (262144 * (rune)(val-240))
			}
		} else {
			count -= 1
			su += ((rune)(math.Pow(64, (float64)(count))) * (rune)(val-128))
			if count == 0 {
				ans = append(ans, su)
				su = 0
			}
		}
	}
	return ans
}

func main() {
	var A string
	fmt.Scanf("%s", &A)
	B := encode(([]rune)(A))
	C := decode(B)
	for _, val := range C {
		fmt.Printf("%c ", val)
	}
}
