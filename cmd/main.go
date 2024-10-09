package main

import (
	"fmt"
	"homework1/internal/storage"
	"log"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	s.Set("string_val", "value1")
	s.Set("int_val", "123")
	s.Set("", "Val")

	res_str := s.Get("string_val")
	res_int := s.Get("int_val")
	res_unknown_val := s.Get("unknown")

	res_getkindstr := s.GetKind("string_val")
	res_getkindint := s.GetKind("int_val")
	res_getkind_unkonown := s.GetKind("unknown")
	fmt.Println(res_str, res_int, res_unknown_val)
	fmt.Println(res_getkindstr, res_getkindint, res_getkind_unkonown)
}
