package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		os.Exit(1)
	}

	s := os.Args[1]

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	var n int
	n, s = readNum(s)
	fmt.Printf("	mov	rax, %d\n", n)

	for len(s) > 0 {
		switch s[0] {
		case byte('+'):
			n, s = readNum(s[1:])
			fmt.Printf("	add	rax, %d\n", n)

		case byte('-'):
			n, s = readNum(s[1:])
			fmt.Printf("	sub	rax, %d\n", n)

		default:
			fmt.Fprintf(os.Stderr, "予期しない文字です: %s", s)
			os.Exit(1)
		}
	}

	fmt.Println("	ret")
}

func readNum(s string) (int, string) {
	var n int
	for len(s) > 0 {
		tmp := s[0] - '0'
		if tmp < 0 || 10 < tmp {
			break
		}
		n *= 10
		n += int(tmp)
		s = s[1:]
	}
	return n, s
}
