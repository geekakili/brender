package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"brender/api/requestlog"
	"brender/api/resource/render"
	"brender/api/router/middleware"
	"brender/util/logger"
)

func New(l *logger.Logger, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		renderer := render.New(l, v)
		r.Method("POST", "/render", requestlog.NewHandler(renderer.Render, l))
	})

	return r
}
