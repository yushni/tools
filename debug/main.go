package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	if _, err := fmt.Fprintf(w, "Hello %s!", r.UserAgent()); err != nil {
		panic(err)
	}
}

func main() {
	router := &http.ServeMux{}
	router.HandleFunc("/", indexHandler)

	port := os.Getenv("HTTP_PORT")

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      enableHttpLogging()(router),
	}

	fmt.Println("server started on port:", port)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func enableHttpLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r)
		})
	}
}
