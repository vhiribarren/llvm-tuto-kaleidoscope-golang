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
	processFile(*filePtr)

}

func processFile(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	kaleidoVisitor := visitor.NewVisitorKaleido()
	if err := consumeAndProcess(string(data), &kaleidoVisitor); err != nil {
		fmt.Println(err)
	}
}

func startREPL() {
	reader := bufio.NewReader(os.Stdin)
	kaleidoVisitor := visitor.NewVisitorKaleido()
	for {
		fmt.Print("kaleido> ")
		input, _ := reader.ReadString('\n')
		if err := consumeAndProcess(input, &kaleidoVisitor); err != nil {
			fmt.Println(err)
		}
	}
}

func consumeAndProcess(input string, kaleidoVisitor *visitor.VisitorKaleido) error {
	kaleidoAST, err := yacc.BuildKaleidoAST(input)
	if err != nil {
		return err
	}
	err = kaleidoVisitor.FeedAST(kaleidoAST)
	if err != nil {
		return err
	}
	println(kaleidoVisitor.GenerateLastModuleIR())
	res, err := kaleidoVisitor.EvalutateMain()
	if err != nil {
		return err
	}
	fmt.Printf("Main evaluated to: %v\n", res)
	return nil
}
