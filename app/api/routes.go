package api

import (
	"net/http"

	handler "echo-postgres-sample/app/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

// start api with the pgdb and return a chi router
func StartAPI(pgdb *pg.DB) *chi.Mux {
	//get the router
	r := chi.NewRouter()
	//add middleware
	//in this case we will store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	//routes for our service
	r.Route("/comments", func(r chi.Router) {
		r.Post("/", handler.CreateComment)
		r.Get("/", handler.GetComments)
		r.Get("/{commentID}", handler.GetCommentByID)
		r.Put("/{commentID}", handler.UpdateCommentByID)
		r.Delete("/{commentID}", handler.DeleteCommentByID)
	})

	//test route to make sure everything works
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}
