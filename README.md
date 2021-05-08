# LLVM Kaleidoscope Tutorial in Go

_**Warning: WORK IN PROGRESS**_


Up to step 3
Missing Go bindings to pursue.

Base tutorial: https://llvm.org/docs/tutorial/index.html

Launch tests:

    go test ./...

If changed, some files may need to be regenerated:

    go generate lexer/kaleido_lexer.go
    go generate with-yacc/main.go

Launch:

    cd with-yacc
    go run .

To use with LLVM on MacOS 11.2.2

    brew install llvm@12
    brew install libffi

    export CGO_CPPFLAGS="`/usr/local/Cellar/llvm/12.0.0/bin/llvm-config --cppflags`"
    export CGO_CXXFLAGS=-std=c++14
    export CGO_LDFLAGS="`/usr/local/Cellar/llvm/12.0.0/bin/llvm-config --ldflags --libs --system-libs all` -L/usr/local/Cellar/libffi/3.3_3/lib -lffi"
    export CGO_LDFLAGS_ALLOW='-Wl,(-search_paths_first|-headerpad_max_install_names)'
    go build -tags byollvm

