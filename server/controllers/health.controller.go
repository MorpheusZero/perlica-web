package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moprheuszero/perlica-web/server/services"
)

type HealthController struct {
	Router        *chi.Mux
	healthService *services.HealthService
}

func NewHealthController(healthService *services.HealthService) *HealthController {
	controller := &HealthController{
		healthService: healthService,
	}
	controller.Router = chi.NewRouter()
	controller.mapController()
	return controller
}

func (c *HealthController) mapController() {
	c.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(c.healthService.CheckHealth()))
	})
}
