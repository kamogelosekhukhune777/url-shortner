package main

import (
	"fmt"
	"net/http"

	"github.com/kamogelosekhukhune777/url-shortner/internals/handlers"
	"github.com/kamogelosekhukhune777/url-shortner/internals/storage"
)

func main() {
	db := storage.NewInMemoryDB()

	urlHandler := handlers.NewURLHandler(db)

	http.HandleFunc("/shorten", urlHandler.ShortenURL)
	http.HandleFunc("/redirect/", urlHandler.Redirect)
	http.HandleFunc("/update", urlHandler.UpdateURL)
	http.HandleFunc("/delete", urlHandler.Delete)

	fmt.Println("starting server at port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
