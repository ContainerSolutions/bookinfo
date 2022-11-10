package redisdb

import (
	"context"
	"strconv"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/data/redisdb/dao"
	"github.com/go-redis/redis/v9"
)

type redisHelper struct {
	client *redis.Client
}

func (rh redisHelper) FindOne(ctx context.Context, id string) (dao.BookStockDAO, error) {
	val, err := rh.client.Get(ctx, id).Result()
	if err != nil {
		return dao.BookStockDAO{}, err
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return dao.BookStockDAO{}, err
	}
	return dao.BookStockDAO{
		ID:           id,
		CurrentStock: intVal,
	}, nil
}
