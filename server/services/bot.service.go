package services

import (
	"fmt"

	"github.com/moprheuszero/perlica-web/server/database/repositories"
	"github.com/moprheuszero/perlica-web/server/util"
)

type BotService struct {
	botRepo *repositories.BotRepository
}

func NewBotService(botRepo *repositories.BotRepository) *BotService {
	return &BotService{
		botRepo: botRepo,
	}
}

// CreateBot creates a new bot and returns it.
func (s *BotService) CreateBot(name string, description *string, dockerImage string) (*repositories.BotEntity, error) {
	newBot, err := s.botRepo.CreateBot(name, description, dockerImage)
	if err != nil {
		return nil, err
	}
	return newBot, nil
}

// GetBotByID retrieves a bot by its ID.
func (s *BotService) GetBotByID(id int) (*repositories.BotEntity, error) {
	bot, err := s.botRepo.GetBotByID(id)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

// GetBotByObjectID retrieves a bot by its object ID.
func (s *BotService) GetBotByObjectID(objectID string) (*repositories.BotEntity, error) {
	bot, err := s.botRepo.GetBotByObjectID(objectID)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

// DeleteBotByObjectID deletes a bot by its object ID.
func (s *BotService) DeleteBotByObjectID(objectID string) error {
	return s.botRepo.DeleteBotByObjectID(objectID)
}

// StartBotInstance starts a new instance of the bot with the given object ID.
func (s *BotService) StartBotInstance(botObjectID string) {
	bot, err := s.botRepo.GetBotByObjectID(botObjectID)
	if err != nil {
		// Log the error but don't return it to the client since the instance is already starting
		fmt.Printf("Failed to retrieve bot with object ID %s: %v\n", botObjectID, err)
		return
	}
	if bot == nil {
		// Log the error but don't return it to the client since the instance is already starting
		fmt.Printf("Bot with object ID %s not found\n", botObjectID)
		return
	}
	err = util.StartDockerContainer(bot.DockerImage, botObjectID)
	if err != nil {
		fmt.Printf("Failed to start Docker container for bot %s: %v\n", botObjectID, err)
	}
}
