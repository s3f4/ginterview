package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/s3f4/ginterview/api/mocks"
	"github.com/s3f4/ginterview/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoData = []struct {
	name         string
	method       string
	url          string
	body         string
	expectedCode int
	expectedBody string
}{
	{"mongo_handler_parse_error", http.MethodPost, "/mongo", `{"startDate: "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`, 422, `{"code":422,"msg":"json parse error"}`},
	{"mongo_handler_startDate", http.MethodPost, "/mongo", `{"endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`, 422, `{"code":422,"msg":"you must provide a valid startDate"}`},
	{"mongo_handler_endDate", http.MethodPost, "/mongo", `{"startDate": "2016-01-26","minCount": 2700,"maxCount": 3000}`, 422, `{"code":422,"msg":"you must provide a valid endDate"}`},
	{"mongo_handler_minCount", http.MethodPost, "/mongo", `{"startDate": "2016-01-26","endDate": "2018-02-02","maxCount": 3000}`, 422, `{"code":422,"msg":"you must provide a valid minCount"}`},
	{"mongo_handler_minCount", http.MethodPost, "/mongo", `{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700}`, 422, `{"code":422,"msg":"you must provide a valid maxCount"}`},
	{"mongo_handler_not_found", http.MethodGet, "/mongo", `{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`, 404, `{"code":404,"msg":"Not Found"}`},
}

func Test_MongoHandler(t *testing.T) {
	time.Local = time.UTC
	createdAt := time.Now()
	records := []*models.Record{
		{
			ID:         primitive.NewObjectID(),
			Key:        "Key",
			Value:      "Value",
			Counts:     []int{100, 100, 100},
			TotalCount: 300,
			CreatedAt:  createdAt,
		},
	}

	mongoRepository := new(mocks.MongoRepository)

	mongoRepository.On("List", mock.Anything, mock.Anything).Return(records, nil)
	mongoHandler := NewMongoHandler(mongoRepository)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/mongo", strings.NewReader(`{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`))

	mongoHandler.ServeHTTP(w, req)
	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fmt.Sprintf(`{"code":0,"msg":"Success","records":[{"key":"Key","createdAt":"%s","totalCount":300}]}`, createdAt.Format(time.RFC3339Nano)), string(body))
	assert.Equal(t, 200, w.Code)
}

func Test_MongoHandler_Errors(t *testing.T) {
	records := []*models.Record{
		{
			ID:         primitive.NewObjectID(),
			Key:        "Key",
			Value:      "Value",
			Counts:     []int{100, 100, 100},
			TotalCount: 300,
		},
	}
	mongoRepository := new(mocks.MongoRepository)
	mongoRepository.On("List", mock.Anything, mock.Anything).Return(records, nil)

	for _, data := range mongoData {
		mongoHandler := NewMongoHandler(mongoRepository)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
		mongoHandler.ServeHTTP(w, req)

		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data.expectedBody, string(body))
		assert.Equal(t, data.expectedCode, w.Code)
	}
}

func Test_MongoHandler_Repository_Error(t *testing.T) {
	mongoRepository := new(mocks.MongoRepository)

	mongoRepository.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("mongo eror"))

	mongoHandler := NewMongoHandler(mongoRepository)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/mongo", strings.NewReader(`{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`))
	mongoHandler.ServeHTTP(w, req)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{"code":500,"msg":"Internal Server Error"}`, string(body))
	assert.Equal(t, 500, w.Code)
}

func Test_MongoHandler_Not_Found(t *testing.T) {
	mongoRepository := new(mocks.MongoRepository)

	mongoRepository.On("List", mock.Anything, mock.Anything).Return([]*models.Record{}, nil)

	mongoHandler := NewMongoHandler(mongoRepository)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/mongo", strings.NewReader(`{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}`))
	mongoHandler.ServeHTTP(w, req)

	res := w.Result()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{"code":404,"msg":"Not Found"}`, string(body))
	assert.Equal(t, 404, w.Code)
}
