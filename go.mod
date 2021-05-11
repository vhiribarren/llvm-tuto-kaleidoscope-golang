module github.com/vhiribarren/tuto-llvm-kaleidoscope-golang

go 1.16

require (
	github.com/llvm/llvm-project v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.1.0
)

replace github.com/llvm/llvm-project => ./llvm-bindings/
