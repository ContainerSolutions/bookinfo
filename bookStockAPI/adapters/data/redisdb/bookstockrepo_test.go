package redisdb

import (
	"context"
	"errors"
	"testing"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/data/redisdb/dao"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/domain"
	"github.com/stretchr/testify/assert"
)

// // MockHTTPClient is the client that mocks original http.Client
type MockMongoHelper struct {
}

var (
	// GetFindOneFunc will be used to get different FindOne functions for testing purposes
	GetFindOneFunc func(ctx context.Context, id string) (dao.BookStockDAO, error)
)

func (mh MockMongoHelper) FindOne(ctx context.Context, id string) (dao.BookStockDAO, error) {
	return GetFindOneFunc(ctx, id)
}

func TestBookStockRepository_FindOne_Error(t *testing.T) {
	pr := BookStockRepository{MockMongoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.BookStockDAO, error) {
		return dao.BookStockDAO{}, errors.New("Cannot find the bookInfo with the ID this_id")
	}
	pDAO, err := pr.Get("this_id")
	assert.Equal(t, pDAO, domain.BookStock{})
	assert.EqualError(t, err, "Cannot find the bookInfo with the ID this_id")
}

func TestBookStockRepository_FindOne_Success(t *testing.T) {
	pr := BookStockRepository{MockMongoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.BookStockDAO, error) {
		return dao.BookStockDAO{
			ID:           "id",
			CurrentStock: 10,
		}, nil
	}
	pDAO, err := pr.Get("this_id")
	assert.Equal(t, pDAO, domain.BookStock{
		ID:           "id",
		CurrentStock: 10,
	})
	assert.Nil(t, err)
}
