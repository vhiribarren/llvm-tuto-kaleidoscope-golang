//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go kaleido_grammar.y

package main

func main() {
	data := "def test(arg) 1+2+go()"
	Parse(data)
}
