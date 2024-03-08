package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var str string

func peek() string {
	return (string)(str[0])
}

func next() string {
	ans := (string)(str[0])
	str = str[1:]
	return ans
}

func removespaces(str string) string {
	var ans string
	for _, val := range str {
		if string(val) != " " {
			ans += string(val)
		}
	}
	return ans
}

// <expr> ::= "(" <inner> ")" | DIGIT
// <inner> ::= "*" <expr> <expr> | <inner2>
// <inner2> ::= "+" <expr> <expr> | <inner3>
// <inner3> ::= "-" <expr> <expr> | <expr>

func expr() int {
	if peek() == "(" {
		next()
		k := inner()
		next()
		return k
	} else {
		ans, _ := strconv.Atoi(next())
		return ans
	}
}

func inner() int {
	if peek() == "*" {
		next()
		return expr() * expr()
	} else {
		return inner2()
	}
}

func inner2() int {
	if peek() == "+" {
		next()
		return expr() + expr()
	} else {
		return inner3()
	}
}

func inner3() int {
	if peek() == "-" {
		next()
		return expr() - expr()
	} else {
		return expr()
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str = scanner.Text()
	str = removespaces(str)
	ans := expr()
	fmt.Printf("%d\n", ans)
}
