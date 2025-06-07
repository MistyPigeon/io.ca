package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	urlStore = make(map[string]string)
	mutex    = &sync.Mutex{}
	baseURL  = "http://localhost:8080/"
	chars    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func generateShortCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]rune, n)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	shortCode := generateShortCode(6)

	mutex.Lock()
	urlStore[shortCode] = longURL
	mutex.Unlock()

	shortURL := baseURL + shortCode
	fmt.Fprintf(w, "Shortened URL: %s\n", shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // Remove leading "/"

	mutex.Lock()
	longURL, ok := urlStore[code]
	mutex.Unlock()

	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
