package server

import (
	// "bytes"
	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type HandlerGet struct {
	Name string `json:"name"`
}

func testGet(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

// func testSet(ctx *gin.Context) {

// }
func testRequests() *gin.Engine {
	r := gin.Default()
	r.GET("/scalar/get/fisrt", testGet)
	return r
}

func TestHandlerGet(t *testing.T) {
	router := testRequests()

	w := httptest.NewRecorder()

	//requestBody := HandlerGet{Name: "first"}
	// jsonData, err := json.Marshal(requestBody)
	// if err != nil {
	// 	t.Errorf("new request: %v", err)
	// }

	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/fisrt", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
