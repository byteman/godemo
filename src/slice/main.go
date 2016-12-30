// slice project main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("insider hello ")
	fmt.Fprintf(w, "<h1>%s</h1>", "hello")
}
func main() {

	http.HandleFunc("/", HelloServer)
	be
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Listen failed", err.Error())
	}
}
