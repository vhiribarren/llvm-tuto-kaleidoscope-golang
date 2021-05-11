/*
MIT License

Copyright (c) 2021 Vincent Hiribarren

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser/yacc/parser.go parser/yacc/kaleido_grammar.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/parser/yacc"
	"github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/visitor"
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
	kaleidoVisitor := visitor.NewVisitorKaleido()
	IRResult, err := kaleidoVisitor.GenerateIR(kaleidoAST)
	if err != nil {
		fmt.Println(err)
		return
	}
	print(IRResult)
}

func startREPL() {
	reader := bufio.NewReader(os.Stdin)
	kaleidoVisitor := visitor.NewVisitorKaleido()
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
