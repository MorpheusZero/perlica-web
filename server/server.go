package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/moprheuszero/perlica-web/config"
	"github.com/moprheuszero/perlica-web/constants"
	"github.com/moprheuszero/perlica-web/server/controllers"
	"github.com/moprheuszero/perlica-web/server/database"
	"github.com/moprheuszero/perlica-web/server/services"
	"github.com/moprheuszero/perlica-web/server/valkey"
)

type AppServer struct {
	serverStartTime int64
	envProvider     *config.EnvProvider
}

func NewAppServer() *AppServer {
	return &AppServer{
		serverStartTime: int64(time.Now().Unix()),
		envProvider:     config.NewEnvProvider(),
	}
}

func (s *AppServer) Start() error {
	fmt.Printf("Starting Perlica Web Server - v%s (Build: %s)\n", constants.AppReleaseVersion, constants.APIBuildVersion)

	// Initialize Valkey Client
	valkeyClient := valkey.NewValkeyClient()
	err := valkeyClient.Initialize(s.envProvider.GetValkeyConnectionString())
	if err != nil {
		return err
	}

	// Setup DB Connection
	database := database.NewDatabase()
	err = database.Initialize(s.envProvider.GetDBConnectionString())
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Setup Services
	healthService := services.NewHealthService()
	templateService := services.NewTemplateService()
	staticService := services.NewStaticService()

	// Setup Controllers & Routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/api/health", controllers.NewHealthController(healthService).Router)

	// Configure UI Controller (at root level)
	router.Mount("/", controllers.NewUIController(templateService, staticService).MapController())

	fmt.Printf("Server started in %d seconds\n", int(time.Now().Unix()-s.serverStartTime))
	fmt.Println("Server started successfully on http://0.0.0.0:3000")
	return http.ListenAndServe("0.0.0.0:3000", router)
}
