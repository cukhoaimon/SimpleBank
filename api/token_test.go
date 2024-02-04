package api

import (
	mockdb "github.com/cukhoaimon/SimpleBank/db/mock"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestServer_renewAccessTokenUser(t *testing.T) {
	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{ /* TODO: add test cases*/ },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
