package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"go.uber.org/zap"
)

var (
	Root_dict = "/Users/vadim/Desktop/golang/fifth lesson/BolshoiGolangProject"
)

type Storage struct {
	InnerString map[string]string
	InnerInt    map[string]int
	InnerArray  map[string][]string
	InnerKeys   map[string]struct{}
	logger      *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	defer logger.Sync()
	logger.Info("created new storage")
	return Storage{
		InnerString: make(map[string]string),
		InnerInt:    make(map[string]int),
		InnerArray:  make(map[string][]string),
		InnerKeys:   make(map[string]struct{}),
		logger:      logger,
	}, nil
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

	r.logger.Info("json file read")
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

	r.logger.Info("json file saved")
	return nil
}

func (r *Storage) Lpush(key string, list []string) []string {
	defer r.logger.Sync()
	slices.Reverse(list)

	if _, ok := r.InnerArray[key]; !ok {
		r.InnerArray[key] = list
		r.InnerKeys[key] = struct{}{}
		r.logger.Info("List set")
		return r.InnerArray[key]
	} else {
		r.InnerArray[key] = append(list, r.InnerArray[key]...)
		r.logger.Info("values append in list in left")
		return r.InnerArray[key]
	}
}

func (r Storage) Rpush(key string, list []string) []string {
	defer r.logger.Sync()
	if _, ok := r.InnerArray[key]; !ok {
		r.InnerArray[key] = list
		r.InnerKeys[key] = struct{}{}
		r.logger.Info("List set")
		return r.InnerArray[key]
	} else {
		r.InnerArray[key] = append(r.InnerArray[key], list...)
		r.logger.Info("values append in list in right")
		return r.InnerArray[key]
	}
}

func (r Storage) Raddtoset(key string, list []string) {
	NewSet := make(map[string]struct{})
	for _, value_set := range r.InnerArray[key] {
		NewSet[value_set] = struct{}{}
	}
	for _, value := range list {
		if _, check := NewSet[value]; !check {
			r.InnerArray[key] = append(r.InnerArray[key], value)
			NewSet[value] = struct{}{}
		}
	}
}

func (r Storage) Check_arr(key string) ([]string, error) {
	if _, err := r.InnerArray[key]; err {
		return r.InnerArray[key], nil
	}
	return nil, errors.New("key does not exist")
}

func (r Storage) Lpop(key string, values []int) ([]string, error) {
	defer r.logger.Info("LPop done")
	defer r.logger.Sync()
	if _, err := r.InnerArray[key]; err {
		if len(values) == 1 {
			if int(math.Abs(float64(values[0]))) > len(r.InnerArray[key]) {
				deleted := r.InnerArray[key]
				r.InnerArray[key] = nil
				return deleted, nil
			}
			end := values[0]
			if end < 0 {
				end = len(r.InnerArray[key]) + end
			}
			deleted := r.InnerArray[key][:end]
			r.InnerArray[key] = r.InnerArray[key][end:]
			return deleted, nil
		} else if len(values) == 2 {
			if int(math.Abs(float64(values[0])))+int(math.Abs(float64(values[1]))) >
				len(r.InnerArray[key]) {
				deleted := r.InnerArray[key]
				r.InnerArray[key] = nil
				return deleted, nil
			}
			start := values[0]
			end := values[1]
			if start < 0 {
				start = len(r.InnerArray[key]) + start
			}
			if end < 0 {
				end = len(r.InnerArray[key]) + end
			}
			end += 1
			if start < 0 || start >= len(r.InnerArray[key]) || end <= start || end > len(r.InnerArray[key]) {
				return nil, errors.New("index does not exit")
			}
			deleted := make([]string, end-start)
			copy(deleted, r.InnerArray[key][start:end])
			r.InnerArray[key] = append(r.InnerArray[key][:start], r.InnerArray[key][end:]...)
			return deleted, nil
		}
	}
	return nil, errors.New("key does not exit")
}

func (r Storage) Rpop(key string, values []int) ([]string, error) {
	defer r.logger.Info("Rpop done")
	defer r.logger.Sync()
	if _, err := r.InnerArray[key]; err {
		if len(values) == 1 {
			deleted := r.InnerArray[key]
			start := values[0]
			end := len(r.InnerArray[key])
			if start < 0 {
				start = -start
				deleted = r.InnerArray[key][0:start]
				r.InnerArray[key] = r.InnerArray[key][start:]
			} else {
				start = len(r.InnerArray[key]) - start
				deleted = r.InnerArray[key][start:end]
				r.InnerArray[key] = r.InnerArray[key][:start]
			}
			return deleted, nil
		} else if len(values) == 2 {
			start := values[0]
			end := values[1]
			if start < 0 {
				start = -start
			} else {
				start = len(r.InnerArray[key]) - values[0]
			}
			if end < 0 {
				end = -end - 1
			} else {
				end = len(r.InnerArray[key]) - values[1]
			}
			start_index, end_index := min(start, end), max(start, end)
			deleted := make([]string, end_index-start_index)
			copy(deleted, r.InnerArray[key][start_index:end_index])
			r.InnerArray[key] = append(r.InnerArray[key][:start_index], r.InnerArray[key][end_index:]...)
			return deleted, nil
		}
		return nil, errors.New("key does not exit")
	}
	return nil, errors.New("key does not exist")
}

func (r Storage) LSet(key string, index uint64, element string) (string, error) {
	if int(index) > len(r.InnerArray[key]) {
		return "", errors.New("index out of range")
	}
	if _, err := r.InnerArray[key]; err {
		r.InnerArray[key][index] = element
		return "OK", nil
	}
	return "", errors.New("key does not exist")
}

func (r Storage) LGet(key string, index int) (string, error) {
	if index < 0 || index > len(r.InnerArray[key]) {
		return "", errors.New("key does not exist")
	}
	return r.InnerArray[key][index], nil
}

func (r *Storage) Set(key string, value interface{}) error {
	defer r.logger.Sync()
	if !r.CheckKeys(key) {
		switch state := value.(type) {
		case string:
			r.InnerString[key] = state
			r.InnerKeys[key] = struct{}{}
			r.logger.Info("key with int value set")
		case int:
			r.InnerInt[key] = state
			r.InnerKeys[key] = struct{}{}
			r.logger.Info("key with string value set")
		default:
			return errors.New("value must be equal a string or a integer")

		}
	}
	return errors.New("keys existed")
}

func (r Storage) Get(key string) (string, error) {
	if _, okint := r.InnerInt[key]; okint {
		return strconv.Itoa(r.InnerInt[key]), nil
	} else if resstring, okstring := r.InnerString[key]; okstring {
		return resstring, nil
	}
	defer r.logger.Sync()
	return "", errors.New("key does not exist")
}

func (r Storage) GetKind(key string) (interface{}, error) {
	defer r.logger.Sync()
	if _, okint := r.InnerInt[key]; okint {
		r.logger.Info("key D sent")
		return "D", nil
	} else if _, okstring := r.InnerString[key]; okstring {
		r.logger.Info("key S sent")
		return "S", nil
	}
	return "", errors.New("key does not exist")
}

func (r Storage) CheckKeys(key string) bool {
	_, exist := r.InnerKeys[key]
	return exist
}
