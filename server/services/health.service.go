package services

type HealthService struct {
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) CheckHealth() string {
	return "ok"
}
