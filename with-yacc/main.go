//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go kaleido_grammar.y

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
	Parse(string(data))
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
