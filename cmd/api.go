package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sdumax/ecom/internal/products"
)

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID) // for rate limitng
	r.Use(middleware.RealIP) // for rate limiting / analytics / tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// timer on request context
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All Good!"))
	})

	productService := products.NewService()
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	

	return r
}

// run
func (app *application) run(h http.Handler) error {
	srv  := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
			ReadTimeout: time.Second * 10,
			IdleTimeout: time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	// logger
	// db driver
}

type config struct {
	addr     string
	dbConfig dbConfig
}

type dbConfig struct {
	dsn string
}
