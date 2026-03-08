package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/moprheuszero/perlica-web/config"
	"github.com/moprheuszero/perlica-web/server/guards"
	"github.com/moprheuszero/perlica-web/server/services"
)

type BotController struct {
	authGuard  *guards.AuthGuard
	config     *config.EnvProvider
	botService *services.BotService
}

func NewBotController(authGuard *guards.AuthGuard, config *config.EnvProvider, botService *services.BotService) *BotController {
	return &BotController{
		authGuard:  authGuard,
		config:     config,
		botService: botService,
	}
}

func (c *BotController) MapController() *chi.Mux {
	router := chi.NewRouter()

	router.With(c.authGuard.ValidateSession).Post("/", c.createBot)
	router.With(c.authGuard.ValidateSession).Post("/{id}/start", c.startBotInstance)

	return router
}

func (c *BotController) createBot(w http.ResponseWriter, r *http.Request) {

	// session := guards.GetSessionFromContext(r)

	data := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		DockerImage string `json:"docker_image"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newBot, err := c.botService.CreateBot(data.Name, &data.Description, data.DockerImage)
	if err != nil {
		http.Error(w, "Failed to create bot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBot)
}

func (c *BotController) startBotInstance(w http.ResponseWriter, r *http.Request) {
	// session := guards.GetSessionFromContext(r)

	botObjectID := chi.URLParam(r, "id")

	// Start the bot instance using the bot service but don't wait for it to complete since it may take a while and we want to return a response immediately
	go func() {
		c.botService.StartBotInstance(botObjectID)
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bot instance is starting"})
}
