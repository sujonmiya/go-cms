package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"service/app/utils"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	router   *mux.Router
	endpoint = fmt.Sprintf("%s%s/", ApiBaseEndpoint, UserApiEndpoint)
)

func init() {
	ua := InitUserApi(mux.NewRouter())
	router = ua.GetRouter()
}

func TestUserApiListHandlerWithoutSendingHeader(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest(Get, endpoint, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusNotFound)
}

func TestUserApiListHandlerUnauthorized(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest(Get, endpoint, nil)
	req.Header.Add(HeaderContentType, MediaTypeJson)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusUnauthorized)
}

func TestUserApiCreateHandlerWithoutSendingHeaders(t *testing.T) {
	assert := assert.New(t)

	data, _ := json.Marshal(utils.NewFakeAdministrator())
	req, _ := http.NewRequest(Post, endpoint, bytes.NewReader(data))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusNotFound)
}

func TestUserApiCreateHandlerCsrf(t *testing.T) {
	assert := assert.New(t)

	data, _ := json.Marshal(utils.NewFakeAdministrator())
	req, _ := http.NewRequest(Post, endpoint, bytes.NewReader(data))
	req.Header.Add(HeaderAccept, MediaTypeJson)
	req.Header.Add(HeaderContentType, MediaTypeJson)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusForbidden)
}
