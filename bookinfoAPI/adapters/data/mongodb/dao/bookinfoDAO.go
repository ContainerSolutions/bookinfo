package dao

// BookInfoDAO represents the struct of BookInfo type to be stored in mongoDB
type BookInfoDAO struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Author string `bson:"author"`
}
