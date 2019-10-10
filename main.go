package main

import (
	"fmt"
	"os"
)

const (
	TkReserved = iota // 記号
	TkNum             // 整数トークン
	TkEOF             // 入力の終わりを表すトークン
)

type TokenKind = int

// Token はトークン型です
type Token struct {
	kind TokenKind // トークンの型
	next *Token    // 次の入力トークン
	val  int       // kind が TkNum の場合，その数値
	str  string    // トークン文字列
}

var token *Token // 現在着目しているトークン

var userInput string

func errorAt(loc string, message string) {
	pos := len(userInput) - len(loc)
	fmt.Fprintf(os.Stderr, "%s\n", userInput)
	for i := 0; i < pos; i++ {
		fmt.Fprint(os.Stderr, " ")
	}
	fmt.Fprintln(os.Stderr, "^", message)
	os.Exit(1)
}

// 次のトークンが期待している記号のときには，トークンを1つ読み進めて真を返す。
// それ以外の場合には偽を返す。
func consume(op byte) bool {
	if token.kind != TkReserved || token.str[0] != op {
		return false
	}
	token = token.next
	return true
}

// 次のトークンが期待している記号のときにはトークンを1つよみすすめる
// それ以外の場合にはエラーを報告する
func expect(op byte) {
	if token.kind != TkReserved || token.str[0] != op {
		errorAt(token.str, fmt.Sprintf("'%s' ではありません", string(op)))
	}
	token = token.next
}

// 次のトークンが数値の場合，トークンを1つ読み進めてその数値を返す。
// それ以外の場合にはエラーを報告する。
func expectNumber() int {
	if token.kind != TkNum {
		errorAt(token.str, "数ではありません")
	}
	val := token.val
	token = token.next
	return val
}

func atEOF() bool {
	return token.kind == TkEOF
}

func newToken(kind TokenKind, cur *Token, str string) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
	}
	cur.next = tok
	return tok
}

func tokenize(s string) *Token {
	var head Token
	head.next = nil
	cur := &head

	for len(s) > 0 {
		if s[0] == ' ' || s[0] == '\n' {
			s = s[1:]
			continue
		}

		if s[0] == '+' || s[0] == '-' {
			cur = newToken(TkReserved, cur, s)
			s = s[1:]
			continue
		}

		if '0' <= s[0] && s[0] <= '9' {
			cur = newToken(TkNum, cur, s)
			var n int
			n, s = readNum(s)
			cur.val = n
			continue
		}

		errorAt(s, "トークナイズできません")
	}

	newToken(TkEOF, cur, s)
	return head.next
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		os.Exit(1)
	}

	userInput = os.Args[1]
	token = tokenize(userInput)

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	// 式の最初は数でなければいけないので，それをチェックして最初の mov 命令を出力
	fmt.Printf("	mov	rax, %d\n", expectNumber())

	// `+ <数>` あるいは `- <数>` というトークンの並びを消費しつつ
	// アセンブリを出力
	for !atEOF() {
		if consume('+') {
			fmt.Printf("	add	rax, %d\n", expectNumber())
			continue
		}

		expect('-')
		fmt.Printf("	sub	rax, %d\n", expectNumber())
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
