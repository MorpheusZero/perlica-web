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
	router.Get("/", c.redirectToLogin)

	router.Get("/login", c.login)

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

func (c *UIController) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirecting root to login page")
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (c *UIController) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving login page")

	pageData := services.PageData{
		Title:       "Login",
		Description: "Login to Perlica",
		Data: map[string]interface{}{
			"PrelaunchMode": false,
		},
	}

	err := c.templateService.RenderTemplate(w, "login", pageData)
	if err != nil {
		fmt.Println("Failed to render login page: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
