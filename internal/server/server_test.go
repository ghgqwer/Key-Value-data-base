package server

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"project_1/internal/storage/storage"
	"strconv"
	"sync"
	"testing"
	"time"

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
		InnerScalar: make(map[string]storage.Scalar),
		InnerArray:  make(map[string]storage.Array),
		InnerKeys:   make(map[string]struct{}),
		Logger:      logger,
	}
}

func TestHandlerGet(t *testing.T) {
	s := NewStorageTest()
	s.Set("testKey", "testValue", 0)

	recorder := httptest.NewRecorder()
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Could not create listener: %v", err)
	}
	defer ln.Close()
	serv := New(ln.Addr().String(), &s)
	router := serv.newApi()

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
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Could not create listener: %v", err)
	}
	defer ln.Close()
	serv := New(ln.Addr().String(), &s)
	router := serv.newApi()

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

	if _, _, err := s.Get(data.Key); err != nil {
		t.Errorf("value doesnt exist: %v", err)
	}
}

func TestRpushArr(t *testing.T) {
	s := NewStorageTest()

	recorder := httptest.NewRecorder()
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Could not create listener: %v", err)
	}
	defer ln.Close()
	serv := New(ln.Addr().String(), &s)
	router := serv.newApi()

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

	if _, _, err := s.CheckArr("testLpush"); err != nil {
		t.Errorf("%v", err)
	}
}

func TestRaddtoset(t *testing.T) {
	s := NewStorageTest()

	recorder := httptest.NewRecorder()
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Could not create listener: %v", err)
	}
	defer ln.Close()
	serv := New(ln.Addr().String(), &s)
	router := serv.newApi()

	s.Rpush("testRaddtoset", []string{"1", "2"}, 0)

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

	arr, _, err := s.CheckArr(testData.Key)
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
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Could not create listener: %v", err)
	}
	defer ln.Close()
	serv := New(ln.Addr().String(), &s)
	router := serv.newApi()

	testData := dataList{
		Key:     "testPop",
		List:    []string{"1", "2", "3"},
		ListInt: []int{1},
	}

	s.Rpush(testData.Key, testData.List, 0)
	if _, ok := s.InnerArray[testData.Key]; !ok {
		t.Errorf("key and value does not add")
	}
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal err: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost,
		"/array/Lpop/"+testData.Key,
		bytes.NewBuffer(jsonData))
	router.ServeHTTP(recorder, req)

	if _, _, err := s.CheckArr(testData.Key); err != nil {
		t.Errorf("key doesnt exist: %v", err)
	}
	assert.Equal(t, http.StatusOK, recorder.Code)

	value, err := s.LGet(testData.Key, 0)
	if value != "2" {
		t.Errorf("value doesnt delete: %v", err)
	}
}

//create recorder
//create server
//create router

//create request
//write down response
//asserts

type ValueBench struct {
	ExpireAt int64
}

type StorageBench struct {
	InnerScalar map[string]ValueBench
	m           sync.RWMutex
}

func (r *StorageBench) SingleClean() {
	realTime := time.Now().UnixMilli()
	r.m.Lock()
	defer r.m.Unlock()

	for key, value := range r.InnerScalar {
		if realTime >= value.ExpireAt {
			delete(r.InnerScalar, key)
		}
	}
}

func (r *StorageBench) MultiClean() {
	r.m.Lock()
	defer r.m.Unlock()

	var wg sync.WaitGroup
	for key, value := range r.InnerScalar {
		wg.Add(1)
		go func(k string, v ValueBench) {
			defer wg.Done()
			if time.Now().Unix() > v.ExpireAt {
				r.m.Lock()
				delete(r.InnerScalar, k)
				r.m.Unlock()
			}
		}(key, value)
	}
	wg.Wait()
}

const numIter = 1000000

func BenchmarkSingleClean(b *testing.B) {
	s := &StorageBench{
		InnerScalar: make(map[string]ValueBench),
	}

	for i := 0; i < numIter; i++ {
		s.InnerScalar[strconv.Itoa(i)] = ValueBench{
			ExpireAt: time.Now().Add(time.Duration(i)).UnixMilli(),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.SingleClean()
	}
}

func BenchmarkMultiClean(b *testing.B) {
	s := &StorageBench{
		InnerScalar: make(map[string]ValueBench),
	}

	for i := 0; i < numIter; i++ {
		s.InnerScalar[strconv.Itoa(i)] = ValueBench{
			ExpireAt: time.Now().Add(time.Duration(i)).UnixMilli(),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.MultiClean()
	}
}

// 10000 iter: BenchmarkSingleClean-8   	22256586	        54.12 ns/op	       0 B/op	       0 allocs/op
// 10000 iter: BenchmarkMultiClean-8   	     478	   2279669 ns/op	  720185 B/op	   20001 allocs/op
// 1000000 iter: BenchmarkSingleClean-8   	20520859	        54.46 ns/op	       0 B/op	       0 allocs/op
// 1000000 iter: BenchmarkMultiClean-8   	       5	 225649667 ns/op	72000412 B/op	 2000002 allocs/op
