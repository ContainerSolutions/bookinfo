package redisdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

// HealthRepository represent a structure that will communicate to MongoDB to accomplish health related transactions
type HealthRepository struct {
	dbClient *redis.Client
}

func newHealthRepository(client *redis.Client) HealthRepository {
	return HealthRepository{
		dbClient: client,
	}
}

// Ready checks the mongodb connection
func (hr HealthRepository) Ready() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check the connection
	status := hr.dbClient.Ping(ctx)
	if status.Err() != nil {
		log.Error().Err(status.Err()).Msg("An error occured while connecting to tha database")
		return false
	}
	log.Info().Msg("Connection to Redis checked successfuly!")
	return true

}
