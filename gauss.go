package main

import (
	"fmt"
	"math"
)

type Rational struct {
	ch int
	zn int
}

func simplify(a Rational) Rational {
	var ans Rational
	if (a.zn < 0) && (a.ch < 0) {
		a.ch = (int)(math.Abs((float64)(a.ch)))
		a.zn = (int)(math.Abs((float64)(a.zn)))
	} else if (a.zn < 0) && (a.ch >= 0) {
		a.ch = -a.ch
		a.zn = -a.zn
	}
	for i := 2; i < a.zn+1; i++ {
		for (a.zn%i == 0) && (a.ch%i == 0) {
			a.zn /= i
			a.ch /= i
		}
	}
	ans.ch = a.ch
	ans.zn = a.zn
	return ans
}

func subtract(a Rational, b Rational) Rational {
	var ans Rational
	ans.ch = a.ch*b.zn - b.ch*a.zn
	ans.zn = a.zn * b.zn
	ans = simplify(ans)
	return ans
}

func multiply(a Rational, b Rational) Rational {
	var ans Rational
	ans.ch = a.ch * b.ch
	ans.zn = a.zn * b.zn
	ans = simplify(ans)
	return ans
}

func swap(num []Rational, ind int, ind2 int) []Rational {
	var ans = num
	var dl = (int)(math.Pow((float64)(len(num)), 0.5))
	for i := 0; i < dl; i++ {
		var temp = ans[ind*dl+i]
		ans[ind*dl+i] = ans[ind2*dl+i]
		ans[ind2*dl+i] = temp
	}
	return ans
}

func printer(numbers []Rational, answers []Rational) {
	N := (int)(math.Pow((float64)(len(numbers)), 0.5))
	for t := 0; t < N; t++ {
		for k := 0; k < N; k++ {
			fmt.Printf("%d/%d ", numbers[t*N+k].ch, numbers[t*N+k].zn)
		}
		fmt.Printf("| %d/%d\n", answers[t].ch, answers[t].zn)
	}
	fmt.Printf("\n")
}

func main() {
	var N int
	var in int
	fmt.Scan(&N)
	numbers := make([]Rational, 0)
	answers := make([]Rational, 0)
	for i := 0; i < N; i++ {
		for k := 0; k < N; k++ {
			var add Rational
			add.zn = 1
			fmt.Scan(&in)
			add.ch = in
			numbers = append(numbers, add)
		}
		var add Rational
		add.zn = 1
		fmt.Scan(&in)
		add.ch = in
		answers = append(answers, add)
	}
	count := 1
	flag := 0
	for i := 0; i < N; i++ {
		cur := numbers[(count-1)*N : count*N]
		if cur[count-1].ch == 0 {
			flag += 1
			if count+flag > N {
				flag = -1
				break
			}
			numbers = swap(numbers, count-1, count-1+flag)
			temp := answers[count-1]
			answers[count-1] = answers[count-1+flag]
			answers[count-1+flag] = temp
			i--
			continue
		}
		flag = 0
		for ind := count; ind < N; ind++ {
			cur[ind].zn *= cur[count-1].ch
			cur[ind].ch *= cur[count-1].zn
			cur[ind] = simplify(cur[ind])
		}
		answers[count-1].zn *= cur[count-1].ch
		answers[count-1].ch *= cur[count-1].zn
		answers[count-1] = simplify(answers[count-1])
		cur[count-1].ch = 1
		cur[count-1].zn = 1
		for k := 0; k < N; k++ {
			if k != count-1 {
				change := numbers[k*N : (k+1)*N]
				koef := change[count-1]
				for t := count - 1; t < N; t++ {
					change[t] = subtract(change[t], multiply(cur[t], koef))
				}
				answers[k] = subtract(answers[k], multiply(answers[count-1], koef))
			}
		}
		count += 1
	}
	if (flag == -1) || ((numbers[N*N-1].ch == 0) && (answers[N-1].ch != 0)) {
		fmt.Printf("No solution\n")
	} else {
		for i := 0; i < N; i++ {
			fmt.Printf("%d/%d\n", answers[i].ch, answers[i].zn)
		}
	}
}
