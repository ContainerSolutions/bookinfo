package redisdb

import (
	"github.com/go-redis/redis/v9"
	"github.com/nicholasjackson/env"
)

var databaseName = env.Int("DatabaseName", false, 0, "The database name for redis")
var connectionString = env.String("ConnectionString", false, "localhost:6379", "Database connection string")
var password = env.String("DbPassword", false, "secret", "Database password")

// DataContext represents a struct that holds concrete repositories
type DataContext struct {
	BookStockRepository BookStockRepository
	HealthRepository    HealthRepository
}

// NewDataContext returns a new mongoDB backed DataContext
func NewDataContext() (DataContext, error) {

	env.Parse()

	// We try to get connectionstring value from the environment variables, if not found it falls back to local database

	client := redis.NewClient(&redis.Options{
		Addr:     *connectionString,
		Password: *password,     // no password set
		DB:       *databaseName, // use default DB
	})

	dataContext := DataContext{}
	dataContext.BookStockRepository = newBookStockRepository(client)
	dataContext.HealthRepository = newHealthRepository(client)
	return dataContext, nil
}
