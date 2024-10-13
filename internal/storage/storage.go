package storage

import (
	"slices"
	"strconv"

	"go.uber.org/zap"
)

type Storage struct {
	innerString map[string]string
	innerInt    map[string]string
	innerArray  map[string][]string
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
		innerString: make(map[string]string),
		innerInt:    make(map[string]string),
		innerArray:  make(map[string][]string),
		logger:      logger,
	}, nil
}

func (r Storage) Lpush(key string, list []string) []string { //, string
	defer r.logger.Sync()
	slices.Reverse(list)
	if _, ok := r.innerArray[key]; !ok {
		r.innerArray[key] = list
		r.logger.Info("List set")
		return r.innerArray[key] //, key
	} else {
		r.innerArray[key] = append(list, r.innerArray[key]...)
		r.logger.Info("values append in list in left")
		return r.innerArray[key] //, key
	}

}

func (r Storage) Rpush(key string, list []string) []string {
	defer r.logger.Sync()
	if _, ok := r.innerArray[key]; !ok {
		r.innerArray[key] = list
		r.logger.Info("List set")
		return r.innerArray[key]
	} else {
		r.innerArray[key] = append(r.innerArray[key], list...)
		r.logger.Info("values append in list in right")
		return r.innerArray[key]
	}

}

func (r Storage) Raddtoset(key string, list []string) {
	new_set := make(map[string]struct{})
	for _, value_set := range r.innerArray[key] {
		new_set[value_set] = struct{}{}
	}
	for _, value := range list {
		if _, check := new_set[value]; !check {
			r.innerArray[key] = append(r.innerArray[key], value)
			new_set[value] = struct{}{}

		}
	}
	defer r.logger.Info("New unique value set")
	defer r.logger.Sync()
}

func (r Storage) Check_arr(key string) []string {
	defer r.logger.Info("arr sent")
	defer r.logger.Sync()
	return r.innerArray[key]
}

func (r Storage) Lpop(key string, values ...int) []string {
	defer r.logger.Info("Pop done")
	defer r.logger.Sync()
	if len(values) == 1 {
		end := values[0]
		if end < 0 {
			end = len(r.innerArray[key]) + end
		}
		deleted := r.innerArray[key][:end]
		r.innerArray[key] = r.innerArray[key][end:]
		return deleted
	} else if len(values) == 2 { //переделать!
		start := values[0]
		end := values[1]
		if start < 0 {
			start = len(r.innerArray[key]) + start
		}
		if end < 0 {
			end = len(r.innerArray[key]) + end
		}
		end += 1
		if start < 0 || start >= len(r.innerArray[key]) || end <= start || end > len(r.innerArray[key]) {
			return nil
		}
		deleted := make([]string, end-start)
		copy(deleted, r.innerArray[key][start:end])
		r.innerArray[key] = append(r.innerArray[key][:start], r.innerArray[key][end:]...)
		return deleted
	}
	return nil
}

func (r Storage) Set(key string, value string) {
	_, errint := strconv.Atoi(value)

	if errint == nil {
		r.innerInt[key] = value
		r.logger.Info("key with int value set")
	} else {
		r.innerString[key] = value
		r.logger.Info("key with string value set")
	}
	defer r.logger.Sync()
}

func (r Storage) Get(key string) string {
	if resint, okint := r.innerInt[key]; okint {
		return resint
	} else if resstring, okstring := r.innerString[key]; okstring {
		return resstring
	}
	defer r.logger.Sync()
	return ""
}

func (r Storage) GetKind(key string) interface{} {
	defer r.logger.Sync()
	if _, okint := r.innerInt[key]; okint {
		r.logger.Info("key D sent")
		return "D"
	} else if _, okstring := r.innerString[key]; okstring {
		r.logger.Info("key S sent")
		return "S"
	}
	return nil
}
