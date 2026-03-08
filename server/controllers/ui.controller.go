package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/moprheuszero/perlica-web/server/guards"
	"github.com/moprheuszero/perlica-web/server/services"
)

type UIController struct {
	authGuard       *guards.AuthGuard
	templateService services.ITemplateService
	staticService   services.IStaticService
}

func NewUIController(authGuard *guards.AuthGuard, templateService services.ITemplateService, staticService services.IStaticService) *UIController {
	return &UIController{
		authGuard:       authGuard,
		templateService: templateService,
		staticService:   staticService,
	}
}

func (c *UIController) MapController() *chi.Mux {
	router := chi.NewRouter()

	// Root redirect to welcome
	router.Get("/", c.redirectToLogin)

	router.Get("/login", c.login)
	router.With(c.authGuard.ValidateSession).Get("/dashboard", c.dashboard)
	router.With(c.authGuard.ValidateSession).Get("/bots", c.bots)

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

func (c *UIController) dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving dashboard page")

	user := c.authGuard.GetUserFromSessionContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pageData := services.PageData{
		Title:       "Dashboard",
		Description: "Welcome to your dashboard",
		Data: map[string]interface{}{
			"User": user,
		},
	}

	err := c.templateService.RenderTemplate(w, "dashboard", pageData)
	if err != nil {
		fmt.Println("Failed to render dashboard page: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (c *UIController) bots(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving bots page")

	user := c.authGuard.GetUserFromSessionContext(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pageData := services.PageData{
		Title:       "Bots Management",
		Description: "Manage your bots",

		Data: map[string]interface{}{
			"User": user,
		},
	}

	err := c.templateService.RenderTemplate(w, "bots", pageData)
	if err != nil {
		fmt.Println("Failed to render bots page: " + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
