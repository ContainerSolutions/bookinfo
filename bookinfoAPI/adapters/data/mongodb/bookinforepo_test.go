package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/dao"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
	"github.com/stretchr/testify/assert"
)

// // MockHTTPClient is the client that mocks original http.Client
type MockMongoHelper struct {
}

var (
	// GetFindOneFunc will be used to get different FindOne functions for testing purposes
	GetFindOneFunc func(ctx context.Context, id string) (dao.BookInfoDAO, error)
	// GetListFunc will be used to get different List functions for testing purposes
	GetListFunc func(ctx context.Context) ([]dao.BookInfoDAO, error)
)

func (mh MockMongoHelper) Find(ctx context.Context) ([]dao.BookInfoDAO, error) {
	return GetListFunc(ctx)
}
func (mh MockMongoHelper) FindOne(ctx context.Context, id string) (dao.BookInfoDAO, error) {
	return GetFindOneFunc(ctx, id)
}

func TestBookInfoRepository_FindOne_Error(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.BookInfoDAO, error) {
		return dao.BookInfoDAO{}, errors.New("Cannot find the bookInfo with the ID this_id")
	}
	pDAO, err := pr.Get("this_id")
	assert.Equal(t, pDAO, domain.BookInfo{})
	assert.EqualError(t, err, "Cannot find the bookInfo with the ID this_id")
}

func TestBookInfoRepository_FindOne_Success(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetFindOneFunc = func(ctx context.Context, id string) (dao.BookInfoDAO, error) {
		return dao.BookInfoDAO{
			ID:     "id",
			Name:   "name",
			Author: "author",
		}, nil
	}
	pDAO, err := pr.Get("this_id")
	assert.Equal(t, pDAO, domain.BookInfo{
		ID:     "id",
		Name:   "name",
		Author: "author",
	})
	assert.Nil(t, err)
}

func TestBookInfoRepository_List_Error(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetListFunc = func(ctx context.Context) ([]dao.BookInfoDAO, error) {
		return nil, errors.New("Whatever error")
	}
	result, err := pr.List()
	assert.Nil(t, result)
	assert.EqualError(t, err, "Error getting BookInfos")
}

func TestBookInfoRepository_List_Success(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	pDAOs := []dao.BookInfoDAO{
		dao.BookInfoDAO{
			ID:     "id1",
			Name:   "name1",
			Author: "author1",
		},
		dao.BookInfoDAO{
			ID:     "id2",
			Name:   "name2",
			Author: "author2",
		},
	}

	GetListFunc = func(ctx context.Context) ([]dao.BookInfoDAO, error) {
		return pDAOs, nil
	}
	result, err := pr.List()
	assert.Nil(t, err)
	assert.Equal(t, result, []domain.BookInfo{
		domain.BookInfo{
			ID:     "id1",
			Name:   "name1",
			Author: "author1",
		},
		domain.BookInfo{
			ID:     "id2",
			Name:   "name2",
			Author: "author2",
		},
	})
}
