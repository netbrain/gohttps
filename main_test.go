package main
import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"strings"
)

func TestPong(t *testing.T) {
	initHandlers()
	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	resp,err := http.Get(server.URL+"/_ping")
	if err != nil {
		t.Fatal(err)
	}

	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "pong!" {
		t.Fatalf("Invalid response of %s",b)
	}
}

func TestHttpsRedirector(t *testing.T) {
	redirectFlag = HTTP2HTTPS
	initHandlers()
	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	req,err := http.NewRequest("GET",server.URL,nil)
	if err != nil {
		t.Fatal(err)
	}

	resp,err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("Expected temporary redirect status, instead got: %d",resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	locationExpectation := strings.Replace(server.URL,"http","https",1)+"/";
	if location !=  locationExpectation{
		t.Fatalf("Expected location to be: %s, instead got: %s",locationExpectation,location)
	}
}

type responseMock struct {
	header http.Header
}

func (r *responseMock) Header() http.Header {
	return r.header
}

func (r *responseMock) Write(b []byte) (int, error){
	return len(b),nil
}

func (r *responseMock) WriteHeader(i int){

}

func TestWildcardRedirector(t *testing.T) {
	//Cant test this End2End due to dns is impossible to override through golang.

	wildcardRedirector := WildcardRedirector(nil)
	req,err := http.NewRequest("GET","http://www.example.org",nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := &responseMock{
		header: make(http.Header),
	}
	wildcardRedirector.ServeHTTP(resp,req)

	if resp.header.Get("Location") != "http://example.org" {
		t.Fatal("Failed to redirect to http://example.org")
	}
}