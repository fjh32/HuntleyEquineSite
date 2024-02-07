package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var PORT = 1954

func setupLogger(logFileName string) *os.File {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// log to stdout and log file
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("Starting new instance of huntleyequine.com")
	return f
}

func loggingMiddleware(next http.Handler) http.Handler {
	hndlr := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr + " requesting " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hndlr)
}

func main() {
	f := setupLogger("./huntleyequine.com.log")
	defer f.Close() // needs to be closed in this scope

	http.Handle("/", loggingMiddleware(http.FileServer(http.Dir("./public"))))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
