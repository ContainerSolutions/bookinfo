package mongodb

import (
	"context"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/data/mongodb/dao"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/application"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoHelper struct {
	coll *mongo.Collection
}

func (mh mongoHelper) Find(ctx context.Context) ([]dao.BookInfoDAO, error) {
	var bookInfoDAOs = make([]dao.BookInfoDAO, 0)
	cur, err := mh.coll.Find(ctx, bson.M{})
	if err != nil {
		log.Error().Err(err).Msgf("Error getting BookInfos")
		return nil, err
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &bookInfoDAOs)
	return bookInfoDAOs, err
}

func (mh mongoHelper) FindOne(ctx context.Context, id string) (dao.BookInfoDAO, error) {
	var bookInfoDAO dao.BookInfoDAO
	err := mh.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&bookInfoDAO)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting BookInfo")
		return dao.BookInfoDAO{}, &application.ErrorCannotFindBookInfo{ID: id}
	}
	return bookInfoDAO, nil
}
