package main

import (
	"github.com/robatussum/personal_web_be/articles"
	"context"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gorilla/mux"
)


func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully waits for existing connections to finish - e.g. 15s or 1m ")
	flag.Parse()

	r := mux.NewRouter()
    
	r.HandleFunc("/articles", articles.ListHandler).Methods("GET")
	r.HandleFunc("/articles/{id}", articles.ContentHandler).Methods("GET", "PUT")
	r.HandleFunc("/articles/tags/{tag}", articles.CategoryHandler).Methods("GET")
	http.Handle("/", r)
	
	srv := &http.Server{
		Addr: "0.0.0.0:5000",
		// Good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
    // SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// block until we receive our signal
	<-c

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	srv.Shutdown(ctx)

	// optionally, you could run srv.Shutdown in a goroutine and block on
    // <-ctx.Done() if your application should wait for other services
    // to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}