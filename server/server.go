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
	"github.com/moprheuszero/perlica-web/server/database/repositories"
	"github.com/moprheuszero/perlica-web/server/guards"
	"github.com/moprheuszero/perlica-web/server/services"
	"github.com/moprheuszero/perlica-web/server/util"
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

	// Setup Repositories
	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)
	botRepo := repositories.NewBotRepository(database)

	// Setup Guards
	authGuard := guards.NewAuthGuard(sessionRepo, userRepo)

	// Setup Services
	healthService := services.NewHealthService()
	templateService := services.NewTemplateService()
	staticService := services.NewStaticService()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(sessionRepo, userService)
	botService := services.NewBotService(botRepo)

	// First Run Check - Create Default Admin User if not exists
	err = s.FirstRunCheck(userService)
	if err != nil {
		return fmt.Errorf("first run check failed: %w", err)
	}

	// Setup Controllers & Routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/api/health", controllers.NewHealthController(healthService).Router)
	router.Mount("/api/auth", controllers.NewAuthController(authGuard, s.envProvider, authService, userService).MapController())
	router.Mount("/api/bots", controllers.NewBotController(authGuard, s.envProvider, botService).MapController())

	// Configure UI Controller (at root level)
	router.Mount("/", controllers.NewUIController(authGuard, templateService, staticService).MapController())

	fmt.Printf("Server started in %d seconds\n", int(time.Now().Unix()-s.serverStartTime))
	fmt.Println("Server started successfully on http://0.0.0.0:3000")
	return http.ListenAndServe("0.0.0.0:3000", router)
}

func (s *AppServer) FirstRunCheck(userService *services.UserService) error {
	_, err := userService.GetUserByUsername("admin")
	if err != nil {
		fmt.Printf("%s", err.Error())
		randomPassword, err := util.GenerateRandomPassword()
		if err != nil {
			return fmt.Errorf("failed to generate random password: %w", err)
		}
		fmt.Printf("This appears to be the first run. Creating default admin user with username 'admin' and password '%s'.\n", randomPassword)
		_, err = userService.CreateUser("admin", "admin", randomPassword)
		if err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}
		fmt.Println("Default admin user created successfully.")
	}
	return nil
}
