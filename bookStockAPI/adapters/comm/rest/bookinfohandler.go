package rest

import (
	"net/http"

	"github.com/ContainerSolutions/bookinfo/bookStockAPI/adapters/comm/rest/mappers"
	"github.com/ContainerSolutions/bookinfo/bookStockAPI/application"
	"github.com/gorilla/mux"
)

type validatedBookStock struct{}

// swagger:route GET /book/{id} BookStock GetBookStock
// Return the the bookInfo with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// GetBookStock gets the bookInfos of the BookStock with the given id
func (ctx *APIContext) GetBookStock(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("BookStock.GetOne", r)
	defer span.Finish()

	// parse the BookStock id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	bookInfoService := application.NewBookStockService(ctx.bookInfoRepo)
	bookInfo, err := bookInfoService.Get(id)
	if err != nil {
		switch err.(type) {
		case *application.ErrorIDFormat:
			respondWithError(rw, r, 400, "Cannot process with the given id")
		case *application.ErrorCannotFindBookStock:
			respondWithError(rw, r, 404, "Cannot get bookInfo from database")
		default:
			respondWithError(rw, r, 500, "Internal server error")
		}
	} else {
		pDTO := mappers.MapBookStock2BookStockResponseDTO(bookInfo)
		respondWithJSON(rw, r, 200, pDTO)
	}
}
