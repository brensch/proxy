package main

import (
	"log"
	"net/http"

	"github.com/brensch/proxy"
)

func main() {
	http.HandleFunc("/", proxy.HandleProxy)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
