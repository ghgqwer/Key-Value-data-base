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

// type dataString {
// 	Key string
// 	Value string
// }

type dataList struct {
	Key     string
	Value   string
	List    []string
	ListInt []int
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

	data := dataList{
		Key:   "testKey",
		Value: "testValue",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("marshal json: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"/scalar/set/"+data.Key,
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	if _, err := s.Get(data.Key); err != nil {
		t.Errorf("value doesnt exist: %v", err)
	}
}

func TestRpushArr(t *testing.T) {
	s := NewStorageTest()
	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	testData := dataList{
		Key:  "testLpush",
		List: []string{"1", "2", "3"},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal json: %v", err)
	}
	req, _ := http.NewRequest(http.MethodPost,
		"/array/Lpush/"+testData.Key,
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	if _, err := s.Check_arr("testLpush"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestRaddtoset(t *testing.T) {
	s := NewStorageTest()
	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	s.Rpush("testRaddtoset", []string{"1", "2"})

	testData := dataList{
		Key:  "testRaddtoset",
		List: []string{"1", "2", "3"},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal err: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"/array/Raddtoset/"+testData.Key,
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	arr, err := s.Check_arr(testData.Key)
	if err != nil {
		t.Errorf("key doesnt exist: %v", err)
	}
	if len(arr) != 3 { //magic numb(
		t.Errorf("req doesnt work")
	}
}

func TestLpop(t *testing.T) {
	s := NewStorageTest()

	recorder := httptest.NewRecorder()
	server := New("localhost:8080", &s)
	router := server.newApi()

	testData := dataList{
		Key:     "testPop",
		List:    []string{"1", "2", "3"},
		ListInt: []int{1},
	}

	s.Lpush(testData.Key, testData.List)

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal err: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"/array/Lpop/"+testData.Key,
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	if _, err := s.Check_arr(testData.Key); err != nil {
		t.Errorf("key doesnt exist: %v", err)
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
} //check Pop

//create recorder
//create server
//create router

//create request
//write down response
//assert
