package rest

import (
	"net/http"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/dto"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/mappers"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/application"
	"github.com/gorilla/mux"
)

type validatedBookInfo struct{}

// swagger:route GET /people BookInfo GetBookInfos
// Return all the bookInfos
// responses:
//	200: OK
//	500: errorResponse

// GetBookInfos gets all the bookInfos of the BookInfo
func (ctx *APIContext) GetBookInfos(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("BookInfo.ListAll", r)
	defer span.Finish()

	bookInfoService := application.NewBookInfoService(ctx.bookInfoRepo)
	bookInfos, err := bookInfoService.List()
	if err != nil {
		respondWithError(rw, r, 500, "Cannot get bookInfos from database")
	} else {
		bookInfoDTOs := make([]dto.BookInfoResponseDTO, 0)
		for _, p := range bookInfos {
			pDTO := mappers.MapBookInfo2BookInfoResponseDTO(p)
			bookInfoDTOs = append(bookInfoDTOs, pDTO)
		}
		respondWithJSON(rw, r, 200, bookInfoDTOs)
	}

}

// swagger:route GET /people/{id} BookInfo GetBookInfo
// Return the the bookInfo with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// GetBookInfo gets the bookInfos of the BookInfo with the given id
func (ctx *APIContext) GetBookInfo(rw http.ResponseWriter, r *http.Request) {
	span := createSpan("BookInfo.GetOne", r)
	defer span.Finish()

	// parse the BookInfo id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	bookInfoService := application.NewBookInfoService(ctx.bookInfoRepo)
	bookInfo, err := bookInfoService.Get(id)
	if err != nil {
		switch err.(type) {
		case *application.ErrorIDFormat:
			respondWithError(rw, r, 400, "Cannot process with the given id")
		case *application.ErrorCannotFindBookInfo:
			respondWithError(rw, r, 404, "Cannot get bookInfo from database")
		default:
			respondWithError(rw, r, 500, "Internal server error")
		}
	} else {
		pDTO := mappers.MapBookInfo2BookInfoResponseDTO(bookInfo)
		respondWithJSON(rw, r, 200, pDTO)
	}
}
