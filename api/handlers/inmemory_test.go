package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/s3f4/ginterview/api/mocks"
	"github.com/stretchr/testify/assert"
)

var inMemoryData = []struct {
	name         string
	method       string
	url          string
	body         string
	expectedCode int
	expectedBody string
}{
	{"inmemory_handler", http.MethodPost, "/in-memory", `{"key":"key","value":"value"}`, 200, `{"key":"key","value":"value"}`},
	{"inmemory_handler_parse_error", http.MethodPost, "/in-memory", `{"key":"key,"value":"value"}`, 422, `{"code":422,"msg":"json parse error"}`},
	{"inmemory_handler_no_value", http.MethodPost, "/in-memory", `{"key":"key"}`, 422, `{"code":422,"msg":"you must provide a valid Value"}`},
	{"inmemory_handler_no_key", http.MethodPost, "/in-memory", `{"value":""}`, 422, `{"code":422,"msg":"you must provide a valid Key"}`},
	{"inmemory_handler_no_url_param", http.MethodGet, "/in-memory", ``, 422, `{"code":422,"msg":"Url Param 'key' is missing"}`},
	{"inmemory_handler_sucess", http.MethodGet, "/in-memory?key=key", ``, 200, `{"key":"key","value":"value"}`},
	{"inmemory_handler_not_found", http.MethodGet, "/in-memory?key=key2", ``, 404, `{"code":404,"msg":"Not Found"}`},
	{"inmemory_handler_url_not_found", http.MethodPut, "/", ``, 404, `{"code":404,"msg":"Not Found"}`},
}

func Test_InmemoryHandler(t *testing.T) {
	InMemoryRepository := new(mocks.InMemoryRepository)
	InMemoryRepository.On("Create", "key", "value").Return(nil)
	InMemoryRepository.On("Get", "key").Return("value")
	InMemoryRepository.On("Exist", "key").Return(true)
	InMemoryRepository.On("Exist", "key2").Return(false)

	for _, data := range inMemoryData {
		inMemoryHandler := NewInMemoryHandler(InMemoryRepository)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(data.method, data.url, strings.NewReader(data.body))
		inMemoryHandler.ServeHTTP(w, req)

		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, data.expectedBody, string(body))
		assert.Equal(t, data.expectedCode, w.Code)
	}
}
