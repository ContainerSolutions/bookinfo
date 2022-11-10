package redisdb

import (
	"context"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/data/redisdb/dao"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

type dbHelper interface {
	FindOne(ctx context.Context, id string) (dao.BookStockDAO, error)
}
