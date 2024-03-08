package main

import (
	"fmt"
)

type operation struct {
	op string
	v1 string
	v2 string
}

var str string

func peek() string {
	return (string)(str[0])
}

func next() string {
	ans := (string)(str[0])
	str = str[1:]
	return ans
}

var list = make([]operation, 0)

// <expr> = "(" + <inner> + ")" | VAR
// <inner> = <operation> + <expr> + <expr>

func expr() string {
	if peek() == "(" {
		next()
		k := inner()
		next()
		return "(" + k + ")"
	} else {
		return next()
	}
}

func inner() string {
	oper := next()
	var1 := expr()
	var2 := expr()
	flag := 0
	for _, val := range list {
		if (val.op == oper) && (val.v1 == var1) && (val.v2 == var2) {
			flag = 1
		}
	}
	if flag == 0 {
		var elem operation
		elem.op = oper
		elem.v1 = var1
		elem.v2 = var2
		list = append(list, elem)
	}
	return oper + var1 + var2
}

func main() {
	fmt.Scanf("%s", &str)
	expr()
	fmt.Printf("%d", len(list))
}
