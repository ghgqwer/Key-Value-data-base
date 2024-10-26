package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	Root_dict = "/Users/vadim/Desktop/golang/sixth lessson/BolshoiGolangProject"
)

const deafultExpireTime = int64(10)

type Scalar struct {
	Value    string `json:"value"`
	Kind     string
	ExpireAt int64 `json:"expireAt"`
}

type Array struct {
	Values   []string `json:"values"`
	ExpireAt int64    `json:"expireAt"`
}

type Storage struct {
	InnerScalar map[string]Scalar
	InnerArray  map[string]Array
	InnerKeys   map[string]struct{}
	Logger      *zap.Logger
}

func NewStorage() (Storage, error) {
	Logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	defer Logger.Sync()
	Logger.Info("created new storage")
	return Storage{
		InnerScalar: make(map[string]Scalar),
		InnerArray:  make(map[string]Array),
		InnerKeys:   make(map[string]struct{}),
		Logger:      Logger,
	}, nil
}

func (r Storage) GarbageCollection(closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			r.Clean()
		}
	}
}

//First realisation

// func (r Storage) Clean() {
// 	var m sync.Mutex
// 	m.Lock()
// 	defer m.Unlock()

// 	for key, value := range r.InnerScalar {
// 		if time.Now().UnixMilli() >= value.ExpireAt {
// 			delete(r.InnerArray, key)
// 		}
// 	}

// 	for key, value := range r.InnerArray {
// 		if time.Now().UnixMilli() >= value.ExpireAt {
// 			delete(r.InnerArray, key)
// 		}
// 	}

// }

func (r Storage) Clean() {
	var (
		wg sync.WaitGroup
		m  sync.Mutex
	)
	numWorkers := 4
	tasks := make(chan string)
	for i := 0; i < numWorkers; i++ {
		m.Lock()
		defer m.Unlock()
		wg.Add(1)
		go r.cleaner(tasks, &wg)
	}

	for key := range r.InnerKeys {
		tasks <- key
	}
}

func (r Storage) cleaner(task <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for key := range task {
		if value, err := r.InnerScalar[key]; err {
			if time.Now().UnixMilli() > value.ExpireAt {
				delete(r.InnerScalar, key)
			}
		}

		if value, err := r.InnerArray[key]; err {
			if time.Now().UnixMilli() >= value.ExpireAt {
				delete(r.InnerArray, key)
			}
		}
	}
}

func (r Storage) WriteAtomic(path string) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	filename := filepath.Base(path)
	tmpPathName := filepath.Join(Root_dict, filename+".tmp")

	err = os.WriteFile(tmpPathName, b, 0o777)
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(tmpPathName)
	}()

	return os.Rename(tmpPathName, Root_dict+path)
}

func (r *Storage) ReadFromJSON(path string) error {
	file_path := filepath.Join(Root_dict, path)
	fromFile, err := os.ReadFile(file_path)
	if err != nil {
		return r.SaveToJSON(path)
	}

	err = json.Unmarshal(fromFile, &r)
	if err != nil {
		return err
	}

	r.Logger.Info("json file read")
	return nil
}

func (r *Storage) SaveToJSON(path string) error {
	file_path := filepath.Join(Root_dict, path)
	file, err := os.Create(file_path)
	if err != nil {
		fmt.Println("Error creating file", err)
		return err
	}
	defer file.Close()

	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error write file", err)
		return err
	}

	err = os.WriteFile(file_path, b, 0o777)
	if err != nil {
		fmt.Println("Error write file", err)
		return err
	}

	r.Logger.Info("json file saved")
	return nil
}

func (r *Storage) Expire(key string, expire int64) error {
	if array, err := r.InnerScalar[key]; err {
		array = r.InnerScalar[key]
		array.ExpireAt = time.Now().Add(time.Duration(expire * int64(time.Second))).UnixMilli()
		r.InnerScalar[key] = array
		return nil
	} else if array, err := r.InnerArray[key]; err {
		array = r.InnerArray[key]
		array.ExpireAt = time.Now().Add(time.Duration(expire * int64(time.Second))).UnixMilli()
		r.InnerArray[key] = array
		return nil
	}
	return errors.New("key doesnt exist")
}

func (r *Storage) Lpush(key string, list []string, expireTime int64) ([]string, error) {
	defer r.Logger.Sync()
	slices.Reverse(list)
	if _, err := r.InnerScalar[key]; !err {
		if _, ok := r.InnerArray[key]; !ok {
			if expireTime == 0 {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(deafultExpireTime * int64(time.Second))).UnixMilli(),
				}
			} else {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(expireTime * int64(time.Second))).UnixMilli(),
				}
			}
			r.InnerKeys[key] = struct{}{}
			r.Logger.Info("List set")
			return r.InnerArray[key].Values, nil
		} else {
			currentArray := r.InnerArray[key]
			currentArray.Values = append(list, currentArray.Values...)
			r.InnerArray[key] = currentArray
			//r.InnerArray[key].Values = append(list, r.InnerArray[key].Values...)
			r.Logger.Info("Values append in list in left")
			return r.InnerArray[key].Values, nil
		}
	}
	return nil, errors.New("key existed")
}

func (r Storage) Rpush(key string, list []string, expireTime int64) ([]string, error) {
	defer r.Logger.Sync()
	if !r.CheckKeys(key) {
		if _, ok := r.InnerArray[key]; !ok {
			if expireTime == 0 {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(deafultExpireTime * int64(time.Second))).UnixMilli(),
				}
			} else {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(expireTime * int64(time.Second))).UnixMilli(),
				}
			}
			r.InnerKeys[key] = struct{}{}
			r.Logger.Info("List set")
			return r.InnerArray[key].Values, nil
		} else {
			currentArray := r.InnerArray[key]
			currentArray.Values = append(currentArray.Values, list...)
			r.InnerArray[key] = currentArray
			//r.InnerArray[key].Values = append(list, r.InnerArray[key].Values...)
			r.Logger.Info("Values append in list in left")
			return r.InnerArray[key].Values, nil
		}
	}
	return nil, errors.New("key existed")
}

func (r Storage) Raddtoset(key string, list []string) {
	NewSet := make(map[string]struct{})
	for _, value_set := range r.InnerArray[key].Values {
		NewSet[value_set] = struct{}{}
	}
	for _, Value := range list {
		if _, check := NewSet[Value]; !check {
			currentArray := r.InnerArray[key]
			currentArray.Values = append(currentArray.Values, Value)
			r.InnerArray[key] = currentArray
		}
	}
}

func (r Storage) CheckArr(key string) ([]string, int64, error) {
	if array, err := r.InnerArray[key]; err {
		if time.Now().UnixMilli() >= array.ExpireAt {
			delete(r.InnerArray, key)
			delete(r.InnerKeys, key)
			return nil, 0, errors.New("key does not exist")
		}
		array.ExpireAt = time.Now().Add(time.Duration(deafultExpireTime * int64(time.Second))).UnixMilli()
		return r.InnerArray[key].Values, r.InnerArray[key].ExpireAt, nil
	}
	return nil, 0, errors.New("key does not exist")
}

func (r Storage) Lpop(key string, Values []int) ([]string, error) {
	defer r.Logger.Info("LPop done")
	defer r.Logger.Sync()
	if array, err := r.InnerArray[key]; err {
		if len(Values) == 1 {
			if int(math.Abs(float64(Values[0]))) > len(r.InnerArray[key].Values) {
				deleted := array.Values
				array.Values = nil
				r.InnerArray[key] = array
				return deleted, nil
			}
			end := Values[0]
			if end < 0 {
				end = len(r.InnerArray[key].Values) + end
			}
			deleted := array.Values[:end]
			array.Values = array.Values[end:]
			r.InnerArray[key] = array
			return deleted, nil
		} else if len(Values) == 2 {
			if int(math.Abs(float64(Values[0])))+int(math.Abs(float64(Values[1]))) >
				len(array.Values) {
				deleted := array.Values
				array.Values = nil
				r.InnerArray[key] = array
				return deleted, nil
			}
			start := Values[0]
			end := Values[1]
			if start < 0 {
				start = len(array.Values) + start
			}
			if end < 0 {
				end = len(array.Values) + end
			}
			end += 1
			if start < 0 || start >= len(array.Values) || end <= start || end > len(array.Values) {
				return nil, errors.New("index does not exit")
			}
			deleted := make([]string, end-start)
			copy(deleted, array.Values[start:end])
			array.Values = append(array.Values[:start], array.Values[end:]...)
			r.InnerArray[key] = array
			return deleted, nil
		}
	}
	return nil, errors.New("key does not exit")
}

func (r Storage) Rpop(key string, Values []int) ([]string, error) {
	defer r.Logger.Info("Rpop done")
	defer r.Logger.Sync()
	if array, err := r.InnerArray[key]; err {
		if len(Values) == 1 {
			deleted := array.Values
			start := Values[0]
			end := len(array.Values)
			if start < 0 {
				start = -start
				deleted = array.Values[0:start]
				array.Values = array.Values[start:]
				r.InnerArray[key] = array
			} else {
				start = len(array.Values) - start
				deleted = array.Values[start:end]
				array.Values = array.Values[:start]
				r.InnerArray[key] = array
			}
			return deleted, nil
		} else if len(Values) == 2 {
			start := Values[0]
			end := Values[1]
			if start < 0 {
				start = -start
			} else {
				start = len(array.Values) - Values[0]
			}
			if end < 0 {
				end = -end - 1
			} else {
				end = len(array.Values) - Values[1]
			}
			start_index, end_index := min(start, end), max(start, end)
			deleted := make([]string, end_index-start_index)
			copy(deleted, array.Values[start_index:end_index])
			array.Values = append(array.Values[:start_index], array.Values[end_index:]...)
			r.InnerArray[key] = array
			return deleted, nil
		}
		return nil, errors.New("key does not exit")
	}
	return nil, errors.New("key does not exist")
}

func (r Storage) LSet(key string, index uint64, element string) (string, error) {
	if int(index) > len(r.InnerArray[key].Values) {
		return "", errors.New("index out of range")
	}
	if _, err := r.InnerArray[key]; err {
		r.InnerArray[key].Values[index] = element
		return "OK", nil
	}
	return "", errors.New("key does not exist")
}

func (r Storage) LGet(key string, index int) (string, error) {
	if index < 0 || index > len(r.InnerArray[key].Values) {
		return "", errors.New("key does not exist")
	}
	return r.InnerArray[key].Values[index], nil
}

func (r *Storage) Set(key string, Value any, expireTime int64) error {
	defer r.Logger.Sync()
	stringVal := fmt.Sprintf("%v", Value)
	Kind := ""
	switch Value.(type) {
	case string:
		Kind = "S"
	case int:
		Kind = "D"
	default:
		Kind = "NonType"
	}
	if !r.CheckKeys(key) {
		if expireTime != 0 {
			r.InnerScalar[key] = Scalar{
				Value:    stringVal,
				Kind:     Kind,
				ExpireAt: time.Now().Add(time.Duration(expireTime * int64(time.Second))).UnixMilli(), //expireTime
			}
		} else if expireTime == 0 {
			r.InnerScalar[key] = Scalar{
				Value:    stringVal,
				Kind:     Kind,
				ExpireAt: time.Now().Add(time.Duration(deafultExpireTime * int64(time.Second))).UnixMilli(), //deafultExpireTime
			}
		}
		r.InnerKeys[key] = struct{}{}
	}
	return errors.New("keys existed")
}

func (r Storage) Get(key string) (string, int64, error) {
	if val, ok := r.InnerScalar[key]; ok {
		if time.Now().UnixMilli() >= val.ExpireAt {
			delete(r.InnerScalar, key)
			delete(r.InnerKeys, key)
			return "", 0, errors.New("expired")
		}
		val.ExpireAt = time.Now().Add(time.Duration(deafultExpireTime * int64(time.Second))).UnixMilli()
		return val.Value, val.ExpireAt, nil
	}
	defer r.Logger.Sync()
	return "", 0, errors.New("key does not exist")
}

// func (r Storage) GetKind(key string) (interface{}, error) {
// 	defer r.Logger.Sync()
// 	if r.CheckKeys(key) {
// 		if _, okint := strconv.Atoi(r.InnerScalar[key].Value); okint == nil {
// 			r.Logger.Info("key D sent")
// 			return "D", nil
// 		} else {
// 			r.Logger.Info("key S sent")
// 			return "S", nil
// 		}
// 	}
// 	return "", errors.New("key does not exist")
// }

func (r Storage) GetKind(key string) (string, error) {
	defer r.Logger.Sync()
	if r.CheckKeys(key) {
		return r.InnerScalar[key].Kind, nil
	}
	return "", errors.New("key does not exist")
}

func (r Storage) CheckKeys(key string) bool {
	_, exist := r.InnerKeys[key]
	return exist
}
