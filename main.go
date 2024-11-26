package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"net/http"
	"strings"
	"math/rand" // not safety critical --- no use for crypto/rand
)

var randomChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// @via: https://stackoverflow.com/a/22892986
func makeRand(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = randomChars[rand.Intn(len(randomChars))]
	}
	return string(b)
}

func main() {
	outsideDomain := os.Getenv("OUTSIDE_DOMAIN")
	if outsideDomain == "" {
		outsideDomain = "http://localhost:8080"
	}
	outsideDomain = strings.TrimSuffix(outsideDomain, "/")

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		random := makeRand(5 + rand.Intn(30))
		h := w.Header()
		h.Set("X-Hell", "Empty")
		h.Set("X-Repo", "github.com/cvanloo/purgatory")
		h.Set("Location", outsideDomain + "/" + random)
		w.WriteHeader(http.StatusFound)
	})

	srv := http.Server{
		Addr: listenAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("listening on: %s", listenAddr)
		log.Printf("redirection base path: %s", outsideDomain)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("interrupt received, shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("server shutdown with error: %v", err)
	}
}
