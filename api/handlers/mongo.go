package handlers

import (
	"log"
	"net/http"

	"github.com/s3f4/ginterview/api/models"
	"github.com/s3f4/ginterview/api/repository"
)

type mongoHandler struct {
	repository repository.MongoRepository
}

// NewMongoHandler creates mongoHandler pointer and returns
func NewMongoHandler(
	repository repository.MongoRepository,
) *mongoHandler {
	return &mongoHandler{
		repository: repository,
	}
}

// ServeHTTP the http method of mongoHandler
// works on url/mongo, accepts only post requests
func (h *mongoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == "POST" {
		var mongoReq models.MongoRequest
		if err := Parse(r, &mongoReq); err != nil {
			SendResponse(w, err.Code, models.Response{
				Code: err.Code,
				Msg:  err.Msg,
			})
			return
		}

		// get records
		records, err := h.repository.List(ctx, &mongoReq)
		if err != nil {
			log.Println(err)
			Send500(w)
			return
		}

		// fill response records from records
		var responseRecords []*models.ResponseRecord
		for _, record := range records {
			responseRecord := &models.ResponseRecord{
				Key:        record.Key,
				CreatedAt:  record.CreatedAt,
				TotalCount: record.TotalCount,
			}

			responseRecords = append(responseRecords, responseRecord)
		}

		if len(responseRecords) == 0 {
			Send404(w)
			return
		}

		SendResponse(w, http.StatusOK, models.Response{
			Code:    0,
			Msg:     "Success",
			Records: responseRecords,
		})
		return
	}

	Send404(w)
}
