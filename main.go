package main

import (
	"net/http"
	"strings"
	"flag"
	"fmt"
	"log"
)



var args struct {
	port int
	http2https bool
	www2wild bool
}

func init(){
	flag.BoolVar(&args.http2https,"http2https", true, "redirect from http to https if set to true")
	flag.BoolVar(&args.www2wild,"www2wild", true, "redirect from www.example.org to example.org if set to true")
	flag.IntVar(&args.port, "port", 8080, "the port to start the http server on")
}

func main() {
	flag.Parse()
	if !args.http2https && !args.www2wild {
		log.Println("No redirection has ben set")
	}
	initHandlers()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",args.port), nil))
}

func initHandlers(){
	http.HandleFunc("/_ping",Pong)

	var handler http.Handler
	if args.http2https {
		log.Println("enabling http2https")
		handler = HttpsRedirector(handler);
	}

	if args.www2wild {
		log.Println("enabling www2wild")
		handler = WildcardRedirector(handler);
	}

	http.Handle("/", handler)
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
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func WildcardRedirector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var scheme string
		if strings.HasPrefix(r.Proto, "HTTP/") {
			scheme = "http"
		}else if strings.HasPrefix(r.Proto, "HTTPS/"){
			scheme = "https"
		}
		if strings.HasPrefix(r.Host, "www.") {
			url := strings.Join([]string{
				scheme+"://",
				strings.Replace(r.Host,"www.","",1),
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
