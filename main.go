package main

import (
	"log"
	"net/http"
	"strings"
)

func Http2HttpsRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/_ping" {
			w.Write([]byte("pong!"))
			return
		}

		if strings.HasPrefix(r.Proto, "HTTP/") {
			url := strings.Join([]string{
				"https://",
				r.Host,
				r.RequestURI,
			}, "")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	http.Handle("/", Http2HttpsRedirector(nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
