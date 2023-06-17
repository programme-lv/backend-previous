package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/programme-lv/backend/internal/environment"
)

func main() {
	conf := environment.ReadEnvConfig()
	conf.Print()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sveika, pasaule!"))
	})
	http.ListenAndServe(":3000", r)
}
