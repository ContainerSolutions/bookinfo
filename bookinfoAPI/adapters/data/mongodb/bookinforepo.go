package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/mappers"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/application"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookInfoRepository holds the mongodb client and database name for methods to use
type BookInfoRepository struct {
	helper dbHelper
}

func newBookInfoRepository(client *mongo.Client, databaseName string) BookInfoRepository {
	return BookInfoRepository{
		helper: mongoHelper{coll: client.Database(databaseName).Collection("bookInfos")},
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

// Add adds a new bookInfo to the underlying database.
// It returns the bookInfo inserted on success or error
func (pr BookInfoRepository) Add(p domain.BookInfo) (domain.BookInfo, error) {
	pass := mappers.MapBookInfo2BookInfoDAO(p)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := pr.helper.InsertOne(ctx, pass)
	if err != nil {
		log.Error().Err(err).Msg("Error while writing user")
		return domain.BookInfo{}, errors.New("Cannot insert the bookInfo")
	}
	log.Info().Msgf("User written: %s", result)
	p.ID = pass.ID
	return p, nil
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

// Update updates fields of a single bookInfo from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr BookInfoRepository) Update(id string, p domain.BookInfo) error {
	p.ID = id
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pDAO := mappers.MapBookInfo2BookInfoDAO(p)
	upDoc := bson.D{{Key: "$set", Value: pDAO}}
	result, err := pr.helper.UpdateOne(ctx, id, upDoc)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating the bookInfo with ID: %s", id)
		return errors.New("Error updating the bookInfo")
	}
	if result != 1 {
		log.Error().Err(err).Msgf("Could not found the bookInfo with ID: %s", id)
		return &application.ErrorCannotFindBookInfo{ID: id}
	}
	return nil
}

// Delete selects a single bookInfo from the database with the given unique identifier
// Returns an error if database fails to provide service
func (pr BookInfoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := pr.helper.DeleteOne(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting BookInfo with ID: %s", id)
		return errors.New("Error deleting the bookInfo")
	}
	if result != 1 {
		log.Error().Err(err).Msgf("Could not found the bookInfo with ID: %s", id)
		return &application.ErrorCannotFindBookInfo{ID: id}
	}
	return nil
}
