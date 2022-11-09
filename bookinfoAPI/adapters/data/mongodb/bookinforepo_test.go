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

func TestBookInfoRepository_Delete_Error(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 0, errors.New("Whatever error")
	}
	err := pr.Delete("id")
	assert.EqualError(t, err, "Error deleting the bookInfo")
}

func TestBookInfoRepository_Delete_ResultNotOne(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 0, nil
	}
	err := pr.Delete("this_id")
	assert.EqualError(t, err, "Cannot find the bookInfo with the ID this_id")
}

func TestBookInfoRepository_Delete_ResultSuccess(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetDeleteFunc = func(ctx context.Context, id string) (int, error) {
		return 1, nil
	}
	err := pr.Delete("this_id")
	assert.Nil(t, err)
}

func TestBookInfoRepository_Update_Error(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 0, errors.New("Whatever error")
	}
	err := pr.Update("id", domain.BookInfo{})
	assert.EqualError(t, err, "Error updating the bookInfo")
}

func TestBookInfoRepository_Update_ResultNotOne(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 0, nil
	}
	err := pr.Update("this_id", domain.BookInfo{})
	assert.EqualError(t, err, "Cannot find the bookInfo with the ID this_id")
}

func TestBookInfoRepository_Update_ResultSuccess(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetUpdateFunc = func(ctx context.Context, id string, update interface{}) (int, error) {
		return 1, nil
	}
	err := pr.Update("id", domain.BookInfo{})
	assert.Nil(t, err)
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
			ID:                      "id",
			Name:                    "name",
			BookInfoClass:           1,
			Sex:                     "male",
			Age:                     34,
			Survived:                false,
			SiblingsOrSpousesAboard: 3,
			ParentsOrChildrenAboard: 2,
			Fare:                    4.97,
		}, nil
	}
	pDAO, err := pr.Get("this_id")
	assert.Equal(t, pDAO, domain.BookInfo{
		ID:                      "id",
		Name:                    "name",
		BookInfoClass:           1,
		Sex:                     "male",
		Age:                     34,
		Survived:                false,
		SiblingsOrSpousesAboard: 3,
		ParentsOrChildrenAboard: 2,
		Fare:                    4.97,
	})
	assert.Nil(t, err)
}

func TestBookInfoRepository_InsertOne_Error(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetInsertOneFunc = func(ctx context.Context, document interface{}) (string, error) {
		return "", errors.New("Whatever error")
	}
	bookInfo, err := pr.Add(domain.BookInfo{
		ID: "this_id",
	})
	assert.Equal(t, bookInfo, domain.BookInfo{})
	assert.EqualError(t, err, "Cannot insert the bookInfo")
}

func TestBookInfoRepository_InsertOne_Success(t *testing.T) {
	pr := BookInfoRepository{MockMongoHelper{}}
	GetInsertOneFunc = func(ctx context.Context, document interface{}) (string, error) {
		return "new_id", nil
	}
	bookInfo, err := pr.Add(domain.BookInfo{
		ID: "this_id",
	})
	assert.Equal(t, bookInfo, domain.BookInfo{
		ID: "this_id",
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
			ID:                      "id1",
			Name:                    "name1",
			BookInfoClass:           1,
			Sex:                     "male",
			Age:                     14,
			Survived:                false,
			SiblingsOrSpousesAboard: 1,
			ParentsOrChildrenAboard: 2,
			Fare:                    1.97,
		},
		dao.BookInfoDAO{
			ID:                      "id2",
			Name:                    "name2",
			BookInfoClass:           2,
			Sex:                     "male",
			Age:                     24,
			Survived:                false,
			SiblingsOrSpousesAboard: 2,
			ParentsOrChildrenAboard: 3,
			Fare:                    2.97,
		},
	}

	GetListFunc = func(ctx context.Context) ([]dao.BookInfoDAO, error) {
		return pDAOs, nil
	}
	result, err := pr.List()
	assert.Nil(t, err)
	assert.Equal(t, result, []domain.BookInfo{
		domain.BookInfo{
			ID:                      "id1",
			Name:                    "name1",
			BookInfoClass:           1,
			Sex:                     "male",
			Age:                     14,
			Survived:                false,
			SiblingsOrSpousesAboard: 1,
			ParentsOrChildrenAboard: 2,
			Fare:                    1.97,
		},
		domain.BookInfo{
			ID:                      "id2",
			Name:                    "name2",
			BookInfoClass:           2,
			Sex:                     "male",
			Age:                     24,
			Survived:                false,
			SiblingsOrSpousesAboard: 2,
			ParentsOrChildrenAboard: 3,
			Fare:                    2.97,
		},
	})
}
