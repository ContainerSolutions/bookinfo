package dao

// BookInfoDAO represents the struct of BookInfo type to be stored in mongoDB
type BookInfoDAO struct {
	ID     string `bson:"uuid"`
	Name   string `bson:"Name"`
	Author string `bson:"Author"`
}
