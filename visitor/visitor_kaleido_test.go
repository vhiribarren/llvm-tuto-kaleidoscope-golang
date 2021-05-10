package visitor

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"alea.net/xp/llvm/kaleidoscope/parser/yacc"
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
		_, err = visitor.GenerateIR(ast)
		if err != nil {
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
		_, err = visitor.GenerateIR(ast)
		if err == nil {
			t.Error("Should have triggered an error, file", file, err)
		}
	}
}
