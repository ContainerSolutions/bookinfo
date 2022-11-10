package redisdb

import (
	"context"
	"time"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/data/redisdb/mappers"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/application"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/domain"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

// BookStockRepository holds the mongodb client and database name for methods to use
type BookStockRepository struct {
	helper dbHelper
}

func newBookStockRepository(client *redis.Client) BookStockRepository {
	return BookStockRepository{
		helper: redisHelper{client: client},
	}
}

// Get selects a single bookInfo from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr BookStockRepository) Get(id string) (domain.BookStock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	bookInfoDAO, err := pr.helper.FindOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting BookStock")
		return domain.BookStock{}, &application.ErrorCannotFindBookStock{ID: id}
	}
	return mappers.MapBookStockDAO2BookStock(bookInfoDAO), nil
}
