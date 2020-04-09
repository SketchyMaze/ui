// +build disabled

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const wasm = "/app.wasm"

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc(wasm, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		http.ServeFile(w, r, "."+wasm)
	})

	fmt.Println("Listening at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
