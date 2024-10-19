package main

import (
	"fmt"
	"log"
	"project_1/internal/storage/storage" 
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	err = s.ReadFromJSON("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := s.SaveToJSON("data.json")
		if err != nil {
			log.Fatal(err)
		}
	}()

	s.Set("name", "Anton")
	s.Set("name", "Vadim")
	fmt.Println(s.Get("name"))

	fmt.Println(s.Get("pue"))

	s.Set("key1", 23)
	fmt.Println(s.Get("key1"))

	s.Rpush("fifth", []string{"1", "2", "3"})
	fmt.Println(s.Check_arr("fifth"))
	s.LSet("fifth", 0, "67")
	fmt.Println(s.Check_arr("fifth"))

	s.Lpush("first", []string{"1", "2", "3"})
	llst := s.Lpush("first", []string{"4", "5"})
	s.Rpush("second", []string{"1"})
	rlst := s.Rpush("second", []string{"2", "3"})
	fmt.Println(llst, rlst)
	s.Raddtoset("second", []string{"3", "5", "8", "4", "8", "6"})

	s.Rpush("third", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})

	fmt.Println(s.Check_arr("third"))
	fisrt_step, _ := s.Lpop("third", 2)
	fmt.Println(fisrt_step)
	second_step, _ := s.Lpop("third", 2, -2)
	fmt.Println(second_step)
	fmt.Println(s.Check_arr("third"))

	s.Rpush("fourth", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})
	fmt.Println(s.Check_arr("fourth"))
	first_deleted, _ := s.Rpop("fourth", 4, -2)
	fmt.Println(first_deleted)
	fmt.Println(s.Check_arr("fourth"))

	s.Set("string_val", "value1")
	s.Set("int_val", 123)
	s.Set("", "Val")

	res_str, _ := s.Get("string_val")
	res_int, _ := s.Get("int_val")
	res_unknown_val, _ := s.Get("unknown")
	res_getkindstr, _ := s.GetKind("string_val")
	res_getkindint, _ := s.GetKind("int_val")
	res_getkind_unkonown, _ := s.GetKind("unknown")
	fmt.Println(res_str, res_int, res_unknown_val)
	fmt.Println(res_getkindstr, res_getkindint, res_getkind_unkonown)

	err = s.WriteAtomic("data.json")
	if err != nil {
		log.Fatal(err)
	}

}
