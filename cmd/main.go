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
	llst := s.Lpush("first", []string{"4", "5"})
	s.Rpush("second", []string{"1"})
	rlst := s.Rpush("second", []string{"2", "3"})
	fmt.Println(llst, rlst)
	s.Raddtoset("second", []string{"3", "5", "8", "4", "8", "6"})

	s.Rpush("third", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})
	fmt.Println(s.Check_arr("third"))
	fisrt_step := s.Lpop("third", 2)
	fmt.Println(fisrt_step)
	second_step := s.Lpop("third", 2, -2)
	fmt.Println(second_step)
	fmt.Println(s.Check_arr("third"))

	res_getkindstr := s.GetKind("string_val")
	res_getkindint := s.GetKind("int_val")
	res_getkind_unkonown := s.GetKind("unknown")
	fmt.Println(res_str, res_int, res_unknown_val)
	fmt.Println(res_getkindstr, res_getkindint, res_getkind_unkonown)
}

//list = [1, 2, 3, 4]
//del [0, s]
//s = 1
//s = 4

//1. Создать словарь с ключами - исходными значениями нашего настоящего словваря
//2. Идем циклом по новым значениями и если такого ключа не существует то добавляем к исходному списку новое значение, которое заведомо уникально

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
