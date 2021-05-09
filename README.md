# LLVM Kaleidoscope Tutorial in Go

_**Warning: WORK IN PROGRESS**_

This is a toy project to test LLVM and manipulate other tools (Go, YACC, ...). It
works, but I do not necessarily follow all the best practices and the code is
not necessarily robust against some errors.

It follows the tutorial: https://llvm.org/docs/tutorial/index.html

Here what is currently done, with some differences with the original tutorial:

- Step1: Lexer
    - https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/LangImpl01.html
    - Usage of Go instead of C++ for the whole tutorial
    - Lexing is manually done (no usage of Go `Scanner`, ...)

- Step 2: Parser and AST
    - https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/LangImpl02.html
    - Parsing done using YACC / GoYACC

- Step 3: Intermediate Representation (IR) code generation
    - https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/LangImpl03.html
    - Usage of a Visitor pattern to analyze the abstract syntax tree
    - Usage of official LLVM v12 Go binding (based on LLVM C binding)


## How to run

You must have working/compiled LLVM v12 libraries on your system.

Launch tests:

    go test ./...

If changed, some files may need to be regenerated:

    go generate lexer/kaleido_lexer.go
    go generate with-yacc/main.go

Launch:

    cd with-yacc
    go run .

## Note on LLVM

I had issue in adding LLVM bindings as a Go module. For me, adding the
`github.com/llvm/llvm-project/llvm/bindings/go/llvm` or 
`github.com/llvm/llvm-project/tree/main/llvm/bindings/go/llvm` either created
some `no matching versions for query "latest"` issues or some zip file creation issues.

I tried to use an alternative - https://github.com/llir/llvm - however it does
not really provide bindings with LLVM, it is a lite wrapper to generate IR code
without optimization, etc.

So... I did manually copy files from
`https://github.com/llvm/llvm-project/tree/llvmorg-12.0.0/llvm/bindings/go`,
and added a `go.mod` file to mark it as a Go module.

A working LLVM installation is anyway needed.

To use with LLVM 12 on MacOS 11.2.2, with HomeBrew:


```bash
brew install llvm@12
brew install libffi
```

Before compiling the files, path to LLVM and FFI must be declared.
In my case, if using HomeBrew:

```bash
export CGO_CPPFLAGS="`/usr/local/Cellar/llvm/12.0.0/bin/llvm-config --cppflags`"
export CGO_CXXFLAGS=-std=c++14
export CGO_LDFLAGS="`/usr/local/Cellar/llvm/12.0.0/bin/llvm-config --ldflags --libs --system-libs all` -L/usr/local/Cellar/libffi/3.3_3/lib -lffi"
export CGO_LDFLAGS_ALLOW='-Wl,(-search_paths_first|-headerpad_max_install_names)'
```
Then, the Go program can be compiled.

