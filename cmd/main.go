package main

import (
	"fmt"
	"log"
	"project_1/internal/storage"
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

	s.Lpush("first", []string{"1", "2", "3"})
	s.Lpush("first", []string{"4", "5"})

	s.Rpush("second", []string{"1", "2", "3"})
	s.Rpush("second", []string{"4", "5"})

	res_getkindstr := s.GetKind("string_val")
	res_getkindint := s.GetKind("int_val")
	res_getkind_unkonown := s.GetKind("unknown")
	fmt.Println(res_str, res_int, res_unknown_val)
	fmt.Println(res_getkindstr, res_getkindint, res_getkind_unkonown)

	//fmt.Println(q, k)
}
