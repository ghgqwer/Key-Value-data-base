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
		{"second", "key2", "int_value"},
		{"third", "key3", ""},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)
			sValue := s.Get(c.key)

			if sValue != c.value {
				t.Errorf("values not equal")
			}
		})
	}
}
