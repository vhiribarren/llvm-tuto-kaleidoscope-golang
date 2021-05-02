# LLVM Kaleidoscope Tutorial in Go

_**Warning: WORK IN PROGRESS**_

Base tutorial: https://llvm.org/docs/tutorial/index.html

Launch tests:

    go test ./...

If changed, some files may need to be regenerated:

    go generate lexer/kaleido_lexer.go
    go generate with-yacc/main.go

Launch:

    cd with-yacc
    go run .
