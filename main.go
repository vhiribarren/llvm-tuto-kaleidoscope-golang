//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser/yacc/parser.go parser/yacc/kaleido_grammar.y

package main

import (
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
		usage()
		return
	}
	data, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(err)
	}
	kaleidoAST := yacc.Parse(string(data)).Result()
	log.Printf("AST: %#v", kaleidoAST)
	visitor := parser.NewVisitorKaleido()
	kaleidoAST.Accept(&visitor)
	print(visitor.Module.String())
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
