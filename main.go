package main

import (
	"github.com/HarunBuyuktepe/go-minimal-cache"
	"net/http"
)

var c memory.Cache

func main() {
	c = memory.New()
	http.ListenAndServe(":8080", Router().InitRouter())
}