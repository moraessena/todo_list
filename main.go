package main

import (
	"log"
	"net/http"

	"github.com/danielAang/todo_list/internal"
	"github.com/danielAang/todo_list/todo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Routes(config *internal.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Compress(5),
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/todo", todo.New(config)())
	})

	return router
}

func main() {
	configuration, err := internal.New()
	if err != nil {
		log.Panic("Configuration error", err)
	}
	router := Routes(configuration)
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Printf("Serving application at port %s\n", configuration.Constants.Port)
	log.Fatal(http.ListenAndServe(":8080", router))
}
