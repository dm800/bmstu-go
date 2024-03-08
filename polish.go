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

// <expr> ::= "(" <inner> ")" | DIGIT
// <inner> ::= "*" <expr> <expr> | <inner2>
// <inner2> ::= "+" <expr> <expr> | <inner3>
// <inner3> ::= "-" <expr> <expr> | DIGIT

func expr() int {
	if peek() == "(" {
		return inner()
	} else {
		ans, _ := strconv.Atoi(next())
		return ans
	}
}

func inner() int {
	if peek() == "*" {
		return expr() * expr()
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str = scanner.Text()
	fmt.Printf("%c\n", peek())
	fmt.Printf("%c\n", next())
	fmt.Printf("%c\n", peek())
	ans := expr()
}
