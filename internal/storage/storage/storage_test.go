package storage

import "testing"

type testCase struct {
	name  string
	key   string
	value string
}

func TestSetGet(t *testing.T) {
	cases := []testCase{
		{"first", "key1", "string_value"},
		{"second", "key2", "123"},
		{"third", "key3", ""},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)
			sValue, _ := s.Get(c.key)

			if sValue != c.value {
				t.Errorf("values not equal")
			}
		})
	}
}

type TestCaseGetKind struct {
	name  string
	key   string
	value interface{}
	kind  string
}

func TestGetKind(t *testing.T) {
	cases := []TestCaseGetKind{
		{"first", "key1", "string_value", "S"},
		{"second", "key2", 23, "D"},
		{"third", "key3", "", "S"},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)
			sValueKind, _ := s.GetKind(c.key)

			if sValueKind != c.kind {
				t.Errorf("kinds not equal")
			}
		})
	}
}

type benchmarkSetGet struct {
	name  string
	key   string
	value string
}

func BenchmarkGet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"first", "key1", "value1"},
		{"second", "key2", "123"},
		{"third", "key3", ""},
	}
	s, err := NewStorage()
	if err != nil {
		b.Errorf("new storage: %v", err)
	}

	for _, tCase := range case_BenchmarkGet {
		b.Run(tCase.name, func(bb *testing.B) {
			s.Set(tCase.key, tCase.value)
			bb.ResetTimer()
			for n := 0; n < b.N; n++ {
				s.Get(tCase.key)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"first", "key1", "value1"},
		{"second", "key2", "123"},
		{"third", "key3", ""},
	}

	for _, tCase := range case_BenchmarkGet {
		b.Run(tCase.name, func(bb *testing.B) {
			s, err := NewStorage()
			if err != nil {
				b.Errorf("new storage: %v", err)
			}
			for n := 0; n < b.N; n++ {
				s.Set(tCase.key, tCase.value)
			}
		})
	}
}

func BenchmarkSetGet(b *testing.B) {
	case_BenchmarkGet := []benchmarkSetGet{
		{"first", "key1", "value1"},
		{"second", "key2", "123"},
		{"third", "key3", ""},
	}
	s, err := NewStorage()
	if err != nil {
		b.Errorf("new storage: %v", err)
	}

	for _, tCase := range case_BenchmarkGet {
		b.Run(tCase.name, func(bb *testing.B) {
			bb.ResetTimer()
			for n := 0; n < b.N; n++ {
				s.Set(tCase.key, tCase.value)
				s.Get(tCase.key)
			}
		})
	}
}
