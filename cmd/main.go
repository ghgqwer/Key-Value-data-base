package main

import (
	"fmt"
	"log"
	"project_1/internal/storage"
)

func main() {
	s, err := storage.NewStorage() //s
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
	llst := s.Lpush("first", []string{"4", "5"})

	s.Rpush("second", []string{"1"})
	rlst := s.Rpush("second", []string{"2", "3"})

	s.Raddtoset("second", []string{"3", "5", "8", "4", "8"})

	chahge_lst := s.Check_arr("second")

	//1. Создать словарь с ключами - исходными значениями нашего настоящего словваря
	//2. Идем циклом по новым значениями и если такого ключа не существует то добавляем к исходному списку новое значение, которое заведомо уникально

	res_getkindstr := s.GetKind("string_val")
	res_getkindint := s.GetKind("int_val")
	res_getkind_unkonown := s.GetKind("unknown")
	fmt.Println(res_str, res_int, res_unknown_val)
	fmt.Println(res_getkindstr, res_getkindint, res_getkind_unkonown)
	fmt.Println(llst, rlst)
	fmt.Println(chahge_lst)
}

// func check() {
// 	true_lst := []string{
// 		"first", "second",
// 	}
// 	// fmt.Println(true_lst)
// 	lst := []string{
// 		"first", "fourth",
// 	}

// 	fmt.Println(true_lst)

// 	set := make(map[string]struct{})

// 	for _, value := range true_lst {
// 		set[value] = struct{}{}
// 	}

// 	for _, val := range lst {
// 		if _, err := set[val]; !err {
// 			true_lst = append(true_lst, val)
// 		}
// 	}
// 	fmt.Println(true_lst)
// }

// true_lst := map[string][]string{
// 	"first": {"1", "2"},
// }
// // fmt.Println(true_lst)

// key := "first"
// list := []string{"1", "113"}

// new_set := make(map[string]struct{})
// for _, value := range true_lst[key] {
// 	new_set[value] = struct{}{}
// }

// for _, val := range list {
// 	if _, err := new_set[val]; !err {
// 		true_lst[key] = append(true_lst[key], val)
// 	}
// }

//fmt.Println(true_lst)
