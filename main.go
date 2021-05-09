//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser/yacc/parser.go parser/yacc/kaleido_grammar.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"alea.net/xp/llvm/kaleidoscope/parser"
	"alea.net/xp/llvm/kaleidoscope/parser/yacc"
)

const EMPTY_STRING = ""

func main() {
	filePtr := flag.String("file", EMPTY_STRING, "File container Kaleidoscope program")
	flag.Parse()
	if *filePtr == EMPTY_STRING {
		startREPL()
		return
	}
	data, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(err)
	}
	kaleidoAST, err := yacc.BuildKaleidoAST(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("AST: %#v", kaleidoAST)
	kaleidoVisitor := parser.NewVisitorKaleido()
	IRResult, err := kaleidoVisitor.GenerateIR(kaleidoAST)
	if err != nil {
		fmt.Println(err)
		return
	}
	print(IRResult)
}

func startREPL() {
	reader := bufio.NewReader(os.Stdin)
	kaleidoVisitor := parser.NewVisitorKaleido()
	for {
		fmt.Print("kaleido> ")
		input, _ := reader.ReadString('\n')
		kaleidoAST, err := yacc.BuildKaleidoAST(string(input))
		if err != nil {
			fmt.Println(err)
			continue
		}
		IRResult, err := kaleidoVisitor.GenerateIR(kaleidoAST)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(IRResult)
	}
}
