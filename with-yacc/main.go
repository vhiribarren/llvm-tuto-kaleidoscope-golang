//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go kaleido_grammar.y

package main

func main() {
	data := "1 + 1"
	Parse(data)
}
