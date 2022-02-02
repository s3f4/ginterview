package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/s3f4/ginterview/api/library"
	"github.com/s3f4/ginterview/api/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Parse parses json body and makes validation on struct
func Parse(r *http.Request, model interface{}) *library.ApiError {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		return library.NewApiError(422, "json parse error")
	}

	if err := validate.Struct(model); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			switch err.Field() {
			case "Key":
				return library.NewApiError(422, "you must provide a valid Key")
			case "Value":
				return library.NewApiError(422, "you must provide a valid Value")
			case "MinCount":
				return library.NewApiError(422, "you must provide a valid minCount")
			case "MaxCount":
				return library.NewApiError(422, "you must provide a valid maxCount")
			case "StartDate":
				return library.NewApiError(422, "you must provide a valid startDate")
			case "EndDate":
				return library.NewApiError(422, "you must provide a valid endDate")
			}

		}
		return library.NewApiError(422, "you must provide a valid JSON body")
	}

	return nil
}

//SendResponse returns json response
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp, err := json.Marshal(data)
	if err != nil {
		log.Println("Body JSON Marshal error")
		return
	}

	w.Write(resp)
}

// Send404 sends not found error
func Send404(w http.ResponseWriter) {
	SendResponse(
		w,
		library.Err404.Code, models.Response{
			Code: library.Err404.Code,
			Msg:  library.Err404.Error(),
		})
}

// Send500 sends internal server error
func Send500(w http.ResponseWriter) {
	SendResponse(
		w,
		library.Err500.Code,
		models.Response{
			Code: library.Err500.Code,
			Msg:  library.Err500.Error(),
		},
	)
}
