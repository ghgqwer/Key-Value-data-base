package filework

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"project_1/internal/storage/storage"
)

var (
	Root_dict = "/data/"
)

func WriteAtomic(r storage.Storage, path string) error {
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

func ReadFromJSON(r storage.Storage, path string) error {
	file_path := filepath.Join(Root_dict, path)
	fromFile, err := os.ReadFile(file_path)
	if err != nil {
		return SaveToJSON(r, path)
	}

	err = json.Unmarshal(fromFile, &r)
	if err != nil {
		return err
	}

	return nil
}

func SaveToJSON(r storage.Storage, path string) error {
	file_path := filepath.Join(Root_dict, path)
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

	return nil
}
