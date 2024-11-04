package storage

import (
	"strconv"
	"testing"
)

type testCase struct {
	key   string
	value string
}

func TestSetGet(t *testing.T) {
	cases := []testCase{
		{"key1", "string_value"},
		{"key2", "123"},
		{"key3", ""},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for numb, c := range cases {
		t.Run(strconv.Itoa(numb), func(t *testing.T) {
			s.Set(c.key, c.value, 0)
			sValue, _, _ := s.Get(c.key)

			if sValue != c.value {
				t.Errorf("values not equal")
			}
		})
	}
}

type TestCaseGetKind struct {
	key   string
	value string
	kind  string
}

func TestGetKind(t *testing.T) {
	cases := []TestCaseGetKind{
		{"key1", "string_value", "S"},
		{"key2", "23", "D"},
		{"key3", "", "S"},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for numb, c := range cases {
		t.Run(strconv.Itoa(numb), func(t *testing.T) {
			s.Set(c.key, c.value, 0)
			sValueKind, _ := s.GetKind(c.key)

			if sValueKind != c.kind {
				t.Errorf("kinds not equal")
			}
		})
	}
}

type benchmarkSetGet struct {
	key   string
	value string
}

func BenchmarkGet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"key1", "value1"},
		{"key2", "123"},
		{"key3", ""},
	}

	for numb, tCase := range case_BenchmarkGet {
		b.Run(strconv.Itoa(numb), func(bb *testing.B) {
			s, err := NewStorage()
			if err != nil {
				b.Errorf("new storage: %v", err)
			}
			s.Set(tCase.key, tCase.value, 0)
			bb.ResetTimer()
			for n := 0; n < b.N; n++ {
				s.Get(tCase.key)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"key1", "value1"},
		{"key2", "123"},
		{"key3", ""},
	}

	for numb, tCase := range case_BenchmarkGet {
		b.Run(strconv.Itoa(numb), func(bb *testing.B) {
			s, err := NewStorage()
			if err != nil {
				b.Errorf("new storage: %v", err)
			}
			for n := 0; n < b.N; n++ {
				s.Set(tCase.key, tCase.value, 0)
			}
		})
	}
}

func BenchmarkSetGet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"key1", "value1"},
		{"key2", "123"},
		{"key3", ""},
	}

	for numb, tCase := range case_BenchmarkGet {
		b.Run(strconv.Itoa(numb), func(bb *testing.B) {
			bb.ResetTimer()
			s, err := NewStorage()
			if err != nil {
				b.Errorf("new storage: %v", err)
			}
			for n := 0; n < b.N; n++ {
				s.Set(tCase.key, tCase.value, 0)
				s.Get(tCase.key)
			}
		})
	}
}
