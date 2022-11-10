package dao

// BookStockDAO represents the struct of BookStock type to be stored in mongoDB
type BookStockDAO struct {
	ID           string
	CurrentStock int
}
