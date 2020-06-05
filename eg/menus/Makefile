.PHONY: run
run:
	go run main.go

.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -v -o windows.wasm main_wasm.go

.PHONY: wasm-serve
wasm-serve: wasm
	../wasm-common/serve.sh
