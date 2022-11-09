package application

// HealthRepository is the interface to interact with database
type HealthRepository interface {
	Ready() bool
}

// HealthService is the struct to let outer layers to interact to the health applicatopn
type HealthService struct {
	HealthRepository HealthRepository
}

// NewHealthService creates a new HealthService instance and sets its repository
func NewHealthService(hr HealthRepository) HealthService {
	if hr == nil {
		panic("missing HealthRepository")
	}
	return HealthService{
		HealthRepository: hr,
	}
}

// Ready returns true if underlying reposiroty and its connection is up and running, false otherwise
func (hs HealthService) Ready() bool {
	return hs.HealthRepository.Ready()
}
