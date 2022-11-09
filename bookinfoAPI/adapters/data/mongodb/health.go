package mongodb

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthRepository represent a structure that will communicate to MongoDB to accomplish health related transactions
type HealthRepository struct {
	dbClient *mongo.Client
	dbName   string
}

func newHealthRepository(client *mongo.Client, databaseName string) HealthRepository {
	return HealthRepository{
		dbClient: client,
		dbName:   databaseName,
	}
}

// Ready checks the mongodb connection
func (hr HealthRepository) Ready() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check the connection
	err := hr.dbClient.Ping(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("An error occured while connecting to tha database")
		return false
	}
	log.Info().Msg("Connection to MongoDB checked successfuly!")
	return true

}
