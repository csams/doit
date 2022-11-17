package set

import (
	"testing"
)

func TestAdd(t *testing.T) {
	a := New[string]()
	a.Add("test")
	if !a.Has("test") {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	a := New[string]()
	a.Add("test")
	if !a.Has("test") {
		t.Fail()
	}
	a.Remove("test")
	if a.Has("test") {
		t.Fail()
	}
}

func TestUnion(t *testing.T) {
	a := New[string]()
	b := New[string]()

	a.Add("a", "b", "c")
	b.Add("b", "c", "d")

	i := a.Union(b)
	pass := i.Has("a") && i.Has("b") && i.Has("c") && i.Has("d")
	if !pass {
		t.Fail()
	}
}
func TestIntersection(t *testing.T) {
	a := New[string]()
	b := New[string]()

	a.Add("a", "b", "c")
	b.Add("b", "c", "d")

	i := a.Intersection(b)
	pass := !i.Has("a") && i.Has("b") && i.Has("c") && !i.Has("d")
	if !pass {
		t.Fail()
	}
}
