package handlers

import (
	"log"
	"net/http"

	"github.com/s3f4/ginterview/api/models"
	"github.com/s3f4/ginterview/api/repository"
)

type inMemoryHandler struct {
	repository repository.InMemoryRepository
}

// NewInMemoryHandler creates inMemoryHandler pointer and returns
func NewInMemoryHandler(
	repository repository.InMemoryRepository,
) *inMemoryHandler {
	return &inMemoryHandler{
		repository: repository,
	}
}

// ServeHTTP the http method of inMemoryHandler
// works on url/in-memory?key=key
func (h *inMemoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keys, ok := r.URL.Query()["key"]

		if !ok || len(keys[0]) < 1 {
			SendResponse(w, 422, models.Response{
				Code: 422,
				Msg:  "Url Param 'key' is missing",
			})
			return
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]
		if h.repository.Exist(key) {
			value := h.repository.Get(key)
			SendResponse(w, http.StatusOK, map[string]interface{}{
				"key":   key,
				"value": value,
			})

			return
		}

		Send404(w)
	case "POST":
		inMemoryReq := models.InMemoryRequest{}
		if err := Parse(r, &inMemoryReq); err != nil {
			log.Println(err)
			SendResponse(w, err.Code, models.Response{
				Code: err.Code,
				Msg:  err.Msg,
			})
			return
		}

		// Create inMemory data
		h.repository.Create(inMemoryReq.Key, inMemoryReq.Value)
		SendResponse(w, http.StatusOK, inMemoryReq)
	default:
		Send404(w)
	}

}
