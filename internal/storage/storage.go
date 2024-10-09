package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Storage struct {
	innerString map[string]string
	innerInt    map[string]string
	//
	logger *zap.Logger
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
		logger:      logger,
	}, nil
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
