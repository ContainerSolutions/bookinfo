package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/mappers"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/application"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookInfoRepository holds the mongodb client and database name for methods to use
type BookInfoRepository struct {
	helper dbHelper
}

func newBookInfoRepository(client *mongo.Client, databaseName string) BookInfoRepository {
	return BookInfoRepository{
		helper: mongoHelper{coll: client.Database(databaseName).Collection("details")},
	}
}

// List loads all the bookInfo records from tha database and returns it
// Returns an error if database fails to provide service
func (pr BookInfoRepository) List() ([]domain.BookInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	//var bookInfoDAO dao.BookInfoDAO
	bookInfoDAOs, err := pr.helper.Find(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting BookInfos")
		return nil, errors.New("Error getting BookInfos")
	}
	bookInfos := make([]domain.BookInfo, 0)
	for _, bookInfoDAO := range bookInfoDAOs {
		bookInfo := mappers.MapBookInfoDAO2BookInfo(bookInfoDAO)
		bookInfos = append(bookInfos, bookInfo)
	}
	return bookInfos, nil
}

// Get selects a single bookInfo from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr BookInfoRepository) Get(id string) (domain.BookInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	bookInfoDAO, err := pr.helper.FindOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting BookInfo")
		return domain.BookInfo{}, &application.ErrorCannotFindBookInfo{ID: id}
	}
	return mappers.MapBookInfoDAO2BookInfo(bookInfoDAO), nil
}
