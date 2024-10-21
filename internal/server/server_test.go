package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project_1/internal/storage/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type HandlerGet struct {
	Value string `json:"value"`
}

func NewStorageTest() storage.Storage {
	logger, err := zap.NewProduction()
	if err != nil {
		return storage.Storage{}
	}
	return storage.Storage{
		InnerString: make(map[string]string),
		InnerInt:    make(map[string]int),
		InnerArray:  make(map[string][]string),
		InnerKeys:   make(map[string]struct{}),
		Logger:      logger,
	}
}

func TestHandlerGet(t *testing.T) {
	s := NewStorageTest()
	s.Set("testKey", "testValue")

	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/testKey", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response HandlerGet

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("json: %v", err)
	}

	assert.Equal(t, "testValue", response.Value)
}

func TestHandlerSet(t *testing.T) {
	s := NewStorageTest()

	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	data := map[string]string{
		"key":   "testKey",
		"value": "testValue",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("marshal json: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"/scalar/set/"+data["key"],
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	if _, err := s.Get(data["key"]); err != nil {
		t.Errorf("value doesnt exist: %v", err)
	}
}

//create recorder
//create server
//create router

//create request
//write down response
//assert

func TestRpushArr(t *testing.T) {
	s := NewStorageTest()
	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	type data struct {
		Key  string
		List []string
	}

	testData := data{
		Key:  "testLpush",
		List: []string{"1", "2", "3"},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal json: %v", err)
	}
	req, _ := http.NewRequest(http.MethodPost,
		"/array/Lpush/testLpush",
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	if _, err := s.Check_arr("testLpush"); err != nil {
		t.Errorf("%v", err)
	}
}
