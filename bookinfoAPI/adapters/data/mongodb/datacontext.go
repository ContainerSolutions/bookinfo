package mongodb

import (
	"context"
	"strings"
	"time"

	"github.com/nicholasjackson/env"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseName = env.String("DatabaseName", false, "bookInfo", "The database name for monfodb")
var connectionString = env.String("ConnectionString", false, "mongodb://{username}:{password}@localhost:27017", "Database connection string")
var username = env.String("DbUserName", false, "mongoadmin", "Database username")
var password = env.String("DbPassword", false, "secret", "Database password")

// DataContext represents a struct that holds concrete repositories
type DataContext struct {
	BookInfoRepository BookInfoRepository
	HealthRepository   HealthRepository
}

// NewDataContext returns a new mongoDB backed DataContext
func NewDataContext() (DataContext, error) {

	env.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// We try to get connectionstring value from the environment variables, if not found it falls back to local database

	connstr := strings.Replace(*connectionString, "{username}", *username, -1)
	connstr = strings.Replace(connstr, "{password}", *password, -1)

	client, err := mongo.NewClient(options.Client().ApplyURI(connstr))
	err = client.Connect(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("An error occured while connecting to tha database")
	} else {
		// Check the connection
		err = client.Ping(ctx, nil)

		if err != nil {
			log.Error().Err(err).Msg("An error occured while connecting to tha database")
		} else {
			log.Info().Msg("Connected to MongoDB!")
		}
	}
	dataContext := DataContext{}
	dataContext.BookInfoRepository = newBookInfoRepository(client, *databaseName)
	dataContext.HealthRepository = newHealthRepository(client, *databaseName)
	return dataContext, nil
}
