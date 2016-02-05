package main

import (
	"net/http"
	"strings"
	"flag"
	"fmt"
	"log"
)

const (
	HTTP2HTTPS = 1 << iota
	WWW2WILD
)

var redirectFlag int
var port int

func init(){
	if *flag.Bool("http2https", true, "redirect from http to https if set to true") {
		addFlag(HTTP2HTTPS)
	}

	if *flag.Bool("www2wild", true, "redirect from www.example.org to example.org if set to true") {
		addFlag(WWW2WILD)
	}

	flag.IntVar(&port, "port", 8080, "the port to start the http server on")
}

func main() {
	flag.Parse()
	if redirectFlag == 0 {
		log.Fatal("No redirection has ben set")
	}
	initHandlers()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port), nil))
}

func initHandlers(){
	http.HandleFunc("/_ping",Pong)

	var handler http.Handler
	if isFlagSet(HTTP2HTTPS){
		handler = HttpsRedirector(handler);
	}

	if isFlagSet(WWW2WILD){
		handler = WildcardRedirector(handler);
	}

	http.Handle("/", handler)
}

func isFlagSet(flag int) bool {
	return redirectFlag&flag != 0
}

func addFlag(flag int){
	redirectFlag |= flag
}

func Pong(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("pong!"))
}

func HttpsRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Proto, "HTTP/") {
			url := strings.Join([]string{
				"https://",
				r.Host,
				r.RequestURI,
			}, "")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func WildcardRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Host, "www.") {
			url := strings.Join([]string{
				r.URL.Scheme+"://",
				strings.Replace(r.URL.Host,"www.","",1),
				r.RequestURI,
			}, "")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
