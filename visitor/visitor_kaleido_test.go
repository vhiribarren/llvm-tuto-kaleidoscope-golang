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

package visitor

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/vhiribarren/tuto-llvm-kaleidoscope-golang/parser/yacc"
)

const validProgramDirectory = "../samples/valid"
const invalidProgramDirectory = "../samples/invalid"

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestParseValidProgram(t *testing.T) {
	entries, err := os.ReadDir(validProgramDirectory)
	if err != nil {
		t.Error(err)
	}
	for _, entry := range entries {
		file := path.Join(validProgramDirectory, entry.Name())
		fileContent, err := os.ReadFile(file)
		if err != nil {
			t.Error("File", file, err)
		}
		ast, err := yacc.BuildKaleidoAST(string(fileContent))
		if err != nil {
			t.Error("File", file, err)
		}
		visitor := NewVisitorKaleido()
		if err = visitor.FeedAST(ast); err != nil {
			t.Error("File", file, err)
		}
	}
}

func TestParseInvalidProgram(t *testing.T) {
	entries, err := os.ReadDir(invalidProgramDirectory)
	if err != nil {
		t.Error(err)
	}
	for _, entry := range entries {
		file := path.Join(invalidProgramDirectory, entry.Name())
		fileContent, err := os.ReadFile(file)
		if err != nil {
			t.Error("File", file, err)
		}
		ast, err := yacc.BuildKaleidoAST(string(fileContent))
		if err != nil {
			t.Error("File", file, err)
		}
		visitor := NewVisitorKaleido()
		if err = visitor.FeedAST(ast); err == nil {
			t.Error("Should have triggered an error, file", file, err)
		}
	}
}
