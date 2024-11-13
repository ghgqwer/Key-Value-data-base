package storage

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"go.uber.org/zap"
)

//const deafultExpireTime = int64(0)

type Scalar struct {
	Value    string `json:"value"`
	Kind     string
	ExpireAt int64 `json:"expireAt"`
}

type Array struct {
	Values   []string `json:"values"`
	ExpireAt int64    `json:"expireAt"`
}

type Hash struct {
	Value    map[string]string
	ExpireAt int64
}

type Storage struct {
	InnerScalar  map[string]Scalar
	InnerArray   map[string]Array
	InnerHashmap map[string]Hash
	InnerKeys    map[string]struct{}
	Logger       *zap.Logger
}

func NewStorage() (Storage, error) {
	Logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	Logger.Info("created new storage")
	return Storage{
		InnerScalar:  make(map[string]Scalar),
		InnerArray:   make(map[string]Array),
		InnerHashmap: make(map[string]Hash),
		InnerKeys:    make(map[string]struct{}),
		Logger:       Logger,
	}, nil
}

func (r *Storage) Hset(key string, keyMap string, value string, expireTime int64) {
	hash, exist := r.InnerHashmap[key]
	if !exist {
		hash = Hash{
			Value:    make(map[string]string),
			ExpireAt: expireTime,
		}
	}
	hash.Value[keyMap] = value
	r.InnerHashmap[key] = hash
}

func (r *Storage) Hget(key, keyMap string) (any, error) {
	hashMap, ok := r.InnerHashmap[key]
	if !ok {
		return "", errors.New("undefind key")
	}
	value, ok := hashMap.Value[keyMap]
	if !ok {
		return "", errors.New("undefind hash key")
	}
	return value, nil
}

func (r Storage) GarbageCollection(closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			r.Cleaner(&sync.Mutex{})
		}
	}
}

func (r Storage) Cleaner(m *sync.Mutex) {
	realTime := time.Now().UnixMilli()
	m.Lock()
	defer m.Unlock()

	for key := range r.InnerKeys {
		if value, ok := r.InnerScalar[key]; ok {
			if realTime >= value.ExpireAt && value.ExpireAt != 0 {
				delete(r.InnerScalar, key)
				delete(r.InnerKeys, key)
			}
		}

		if value, ok := r.InnerArray[key]; ok {
			if realTime >= value.ExpireAt && value.ExpireAt != 0 {
				delete(r.InnerArray, key)
				delete(r.InnerKeys, key)
			}
		}

	}
}

func (r *Storage) LoggerSync(closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			r.Logger.Sync()
		}
	}
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
	slices.Reverse(list)
	if _, ok := r.InnerScalar[key]; !ok {
		if _, ok := r.InnerArray[key]; !ok {
			if expireTime == 0 {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: 0,
				}
			} else {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(expireTime * int64(time.Second))).UnixMilli(),
				}
			}
			r.Logger.Info("List set")
		} else {
			currentArray := r.InnerArray[key]
			currentArray.Values = append(list, currentArray.Values...)
			r.InnerArray[key] = currentArray
			//r.InnerArray[key].Values = append(list, r.InnerArray[key].Values...)
			r.Logger.Info("Values append in list in left")
		}
		r.InnerKeys[key] = struct{}{}
		return r.InnerArray[key].Values, nil
	}
	return nil, errors.New("key existed")
}

func (r Storage) Rpush(key string, list []string, expireTime int64) ([]string, error) {
	if _, ok := r.InnerKeys[key]; !ok {
		if _, ok := r.InnerArray[key]; !ok {
			if expireTime == 0 {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: 0,
				}
			} else {
				r.InnerArray[key] = Array{
					Values:   list,
					ExpireAt: time.Now().Add(time.Duration(expireTime * int64(time.Second))).UnixMilli(),
				}
			}
			r.Logger.Info("List set")
		} else {
			currentArray := r.InnerArray[key]
			currentArray.Values = append(currentArray.Values, list...)
			r.InnerArray[key] = currentArray
			r.Logger.Info("Values append in list in left")
		}
		r.InnerKeys[key] = struct{}{}
		return r.InnerArray[key].Values, nil
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
		if time.Now().UnixMilli() >= array.ExpireAt && array.ExpireAt != 0 {
			delete(r.InnerArray, key)
			delete(r.InnerKeys, key)
			return nil, 0, errors.New("key does not exist")
		}
		return r.InnerArray[key].Values, r.InnerArray[key].ExpireAt, nil
	}
	return nil, 0, errors.New("key does not exist")
}

func (r Storage) Lpop(key string, start, end int) ([]string, error) {
	defer r.Logger.Info("LPop done")
	if array, err := r.InnerArray[key]; err {
		if int(math.Abs(float64(start)))+int(math.Abs(float64(end))) >
			len(array.Values) {
			deleted := array.Values
			array.Values = nil
			r.InnerArray[key] = array
			return deleted, nil
		}
		if start < 0 {
			start = len(array.Values) + start
		}
		if end < 0 {
			end = len(array.Values) + end
		}
		end += 1
		if start < 0 || start >= len(array.Values) || end <= start || end > len(array.Values) {
			return nil, errors.New("index out of range")
		}
		deleted := make([]string, end-start)
		copy(deleted, array.Values[start:end])
		array.Values = append(array.Values[:start], array.Values[end:]...)
		r.InnerArray[key] = array
		return deleted, nil
	}
	return nil, errors.New("key does not exit")
}

func (r Storage) Rpop(key string, start, end int) ([]string, error) {
	defer r.Logger.Info("Rpop done")
	if array, err := r.InnerArray[key]; err {
		if int(math.Abs(float64(start)))+int(math.Abs(float64(end))) >
			len(array.Values) {
			deleted := array.Values
			array.Values = nil
			r.InnerArray[key] = array
			return deleted, nil
		}
		if start < 0 {
			start = -start
		} else {
			start = len(array.Values) - start
		}
		if end < 0 {
			end = -end - 1
		} else {
			end = len(array.Values) - end
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
	if _, ok := r.InnerKeys[key]; !ok {
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
				ExpireAt: 0, //deafultExpireTime
			}
		}
		r.InnerKeys[key] = struct{}{}
		return nil
	}
	return errors.New("keys existed")
}

func (r Storage) Get(key string) (string, int64, error) {
	if val, ok := r.InnerScalar[key]; ok {
		if time.Now().UnixMilli() >= val.ExpireAt && val.ExpireAt != 0 {
			delete(r.InnerScalar, key)
			delete(r.InnerKeys, key)
			return "", 0, errors.New("expired")
		}
		return val.Value, val.ExpireAt, nil
	}
	return "", 0, errors.New("key does not exist")
}

func (r Storage) GetKind(key string) (string, error) {
	if _, ok := r.InnerKeys[key]; ok {
		return r.InnerScalar[key].Kind, nil
	}
	return "", errors.New("key does not exist")
}

// func (r Storage) CheckKeys(key string) bool {
// 	_, exist := r.InnerKeys[key]
// 	return exist
// }
