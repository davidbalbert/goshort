package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var urls map[string]string = make(map[string]string)

func slugForUrl(url string) string {
	var slug string

	for {
		slug = fmt.Sprintf("%x", rand.Int31n(1000000))

		if urls[slug] == "" {
			return slug
		}
	}
}

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(422) // UnprocessableEntity
		fmt.Fprintf(w, err.Error())
		return
	}

	if len(r.PostForm["url"]) == 0 {
		w.WriteHeader(422)
		fmt.Fprintf(w, "Missing required param 'url'")
		return
	}

	url := r.PostForm["url"][0]
	slug := slugForUrl(url)

	urls[slug] = url

	fmt.Fprintf(w, "http://%s/%s", r.Host, slug)
}

func lengthen(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
		return
	}

	slug := r.URL.Path[1:]

	if urls[slug] == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, urls[slug], http.StatusMovedPermanently)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", lengthen)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
