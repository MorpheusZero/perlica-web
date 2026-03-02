package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/moprheuszero/perlica-web/server/services"
)

type UIController struct {
	templateService services.ITemplateService
	staticService   services.IStaticService
}

func NewUIController(templateService services.ITemplateService, staticService services.IStaticService) *UIController {
	return &UIController{
		templateService: templateService,
		staticService:   staticService,
	}
}

func (c *UIController) MapController() *chi.Mux {
	router := chi.NewRouter()

	// Root redirect to welcome
	router.Get("/", c.redirectToWelcome)

	router.Get("/welcome", c.welcome)

	// Static file serving
	router.Get("/static/*", c.serveStatic)

	return router
}

func (c *UIController) serveStatic(w http.ResponseWriter, r *http.Request) {
	// Extract the file path from the URL
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")

	err := c.staticService.ServeStaticFile(w, r, filePath)
	if err != nil {
		// Error is already handled and logged in the static service
		return
	}
}

func (c *UIController) redirectToWelcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirecting root to welcome page")
	http.Redirect(w, r, "/welcome", http.StatusMovedPermanently)
}

func (c *UIController) welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving welcome page")

	pageData := services.PageData{
		Title:       "Welcome",
		Description: "Welcome to Parallax - A journey through the rifts of space-time",
		Data: map[string]interface{}{
			"PrelaunchMode": false,
		},
	}

	err := c.templateService.RenderTemplate(w, "welcome", pageData)
	if err != nil {
		fmt.Println("Failed to render welcome page: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
