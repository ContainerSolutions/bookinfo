package mongodb

import (
	"context"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/dao"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

type dbHelper interface {
	Find(ctx context.Context) ([]dao.BookInfoDAO, error)
	FindOne(ctx context.Context, id string) (dao.BookInfoDAO, error)
}
