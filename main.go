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

		if s[0] == '+' || s[0] == '-' || s[0] == '*' || s[0] == '/' || s[0] == '(' || s[0] == ')' {
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

// 抽象構文木のノードの種類
const (
	NdAdd = iota // +
	NdSub        // -
	NdMul        // *
	NdDiv        // /
	NdNum        // 整数
)

type NodeKind = int

type Node struct {
	kind     NodeKind // ノードの型
	lhs, rhs *Node    // 左辺，右辺
	val      int      // kind が NdNum の場合のみ使う
}

func newNode(kind NodeKind, lhs, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNodeNum(val int) *Node {
	return &Node{
		kind: NdNum,
		val:  val,
	}
}

func expr() *Node {
	node := mul()

	for {
		if consume('+') {
			node = newNode(NdAdd, node, mul())
		} else if consume('-') {
			node = newNode(NdSub, node, mul())
		} else {
			return node
		}
	}
}

func mul() *Node {
	node := primary()

	for {
		if consume('*') {
			node = newNode(NdMul, node, primary())
		} else if consume('/') {
			node = newNode(NdDiv, node, primary())
		} else {
			return node
		}
	}
}

func primary() *Node {
	if consume('(') {
		node := expr()
		expect(')')
		return node
	}

	return newNodeNum(expectNumber())
}

func gen(node *Node) {
	if node.kind == NdNum {
		fmt.Printf("	push	%d\n", node.val)
		return
	}

	gen(node.lhs)
	gen(node.rhs)

	fmt.Println("	pop	rdi")
	fmt.Println("	pop	rax")

	switch node.kind {
	case NdAdd:
		fmt.Println("	add	rax, rdi")
	case NdSub:
		fmt.Println("	sub	rax, rdi")
	case NdMul:
		fmt.Println("	imul	rax, rdi")
	case NdDiv:
		fmt.Println("	cqo")
		fmt.Println("	idiv	rdi")
	}

	fmt.Println("	push	rax")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "引数の個数が正しくありません")
		os.Exit(1)
	}

	userInput = os.Args[1]
	token = tokenize(userInput)
	node := expr()

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	gen(node)

	fmt.Println("	pop	rax")
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
