package set

import "encoding/json"

type Set[T comparable] struct {
	back map[T]bool
}

func New[T comparable](ts ...T) *Set[T] {
	s := &Set[T]{
		back: make(map[T]bool),
	}
	return s.Add(ts...)
}

func (s *Set[T]) Add(ts ...T) *Set[T] {
	for _, t := range ts {
		s.back[t] = true
	}
	return s
}

func (s *Set[T]) Remove(ts ...T) *Set[T] {
	for _, t := range ts {
		delete(s.back, t)
	}
	return s
}

func (s *Set[T]) Has(ts ...T) bool {
	for _, t := range ts {
		_, e := s.back[t]
		if !e {
			return false
		}
	}
	return true
}

func (s *Set[T]) Union(o *Set[T]) *Set[T] {
	res := New[T]()
	for key, value := range s.back {
		res.back[key] = value
	}

	for key, value := range o.back {
		res.back[key] = value
	}
	return res
}

func (s *Set[T]) Intersection(o *Set[T]) *Set[T] {
	a, b := s, o
	if len(b.back) < len(a.back) {
		a, b = b, a
	}

	res := New[T]()
	for key, value := range a.back {
		if _, exists := b.back[key]; exists {
			res.back[key] = value
		}
	}
	return res
}

func (s *Set[T]) Size() int {
	return len(s.back)
}

func FromList[T comparable](l []T) *Set[T] {
	res := New[T]()
	if l == nil {
		return res
	}
	for _, v := range l {
		res.back[v] = true
	}
	return res
}

func (s *Set[T]) ToList() []T {
	keys := make([]T, 0, len(s.back))
	for k := range s.back {
		keys = append(keys, k)
	}
	return keys
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var ts []T
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}
	s.back = FromList(ts).back
	return nil
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToList())
}
