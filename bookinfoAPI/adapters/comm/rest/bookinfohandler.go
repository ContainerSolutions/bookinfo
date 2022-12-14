package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/dto"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/adapters/comm/rest/mappers"
	"github.com/ContainerSolutions/bookinfo/bookInfoAPI/application"
	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
)

type validatedBookInfo struct{}

// swagger:route GET /book BookInfo GetBookInfos
// Return all the bookInfos
// responses:
//	200: OK
//	500: errorResponse

// GetBookInfos gets all the bookInfos of the BookInfo
func (ctx *APIContext) GetBookInfos(rw http.ResponseWriter, r *http.Request) {
	span, _ := createSpan("BookInfo.ListAll", r)
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

// swagger:route GET /book/{id} BookInfo GetBookInfo
// Return the the bookInfo with the given id
// responses:
//	200: OK
//  400: Bad Request
//	500: errorResponse

// GetBookInfo gets the bookInfos of the BookInfo with the given id
func (ctx *APIContext) GetBookInfo(rw http.ResponseWriter, r *http.Request) {
	span, tracer := createSpan("BookInfo.GetOne", r)
	defer span.Finish()

	// parse the BookInfo id from the url
	vars := mux.Vars(r)
	id := vars["id"]
	bookInfoService := application.NewBookInfoService(ctx.bookInfoRepo)
	bookInfo, err := bookInfoService.Get(id)
	if err != nil {
		respondWithError(rw, r, 500, "Cannot get bookInfo from database")
		return
	}
	// Call stocks service
	url := os.Getenv("STOCK_URL")
	if url == "" {
		url = "http://localhost:5555"
	}

	url = url + "/book/" + id
	// First prepare the tracing info
	netClient := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("GET", url, nil)
	// Inject the client span context into the headers
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	stockresponse, err := netClient.Do(req)

	stockInfo := &dto.StockInfo{
		CurrentStock: 0,
	}
	if err == nil {
		buf, _ := ioutil.ReadAll(stockresponse.Body)
		json.Unmarshal(buf, &stockInfo)
	}
	bookInfo.CurrentStock = stockInfo.CurrentStock
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
