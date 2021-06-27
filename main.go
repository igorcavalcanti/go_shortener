package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/igorcavalcanti/go_shortener/api"
	"github.com/igorcavalcanti/go_shortener/repository/mongo"
	"github.com/igorcavalcanti/go_shortener/shortener"
)

func main() {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewRedirectHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on Port: 8000")
		errs <- http.ListenAndServe(httpPort(), router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"

	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {
	mongoURL := "mongodb://localhost/shortner"
	mongoDB := "shortner"
	mongoTimeout := 30

	repo, err := mongo.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
