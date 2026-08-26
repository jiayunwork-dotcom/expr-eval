package env

import (
	"math"
	"testing"
)

func TestNewHasConstants(t *testing.T) {
	e := New()
	if v, ok := e.Get("pi"); !ok || math.Abs(v-math.Pi) > 1e-10 {
		t.Errorf("pi = %f", v)
	}
}

func TestSetAndGet(t *testing.T) {
	e := New()
	e.Set("x", 42)
	v, ok := e.Get("x")
	if !ok || v != 42 {
		t.Errorf("x = %f, ok=%v", v, ok)
	}
}

func TestChildScope(t *testing.T) {
	parent := New()
	parent.Set("a", 1)
	child := parent.Child()
	child.Set("b", 2)

	if v, ok := child.Get("a"); !ok || v != 1 {
		t.Error("child should see parent var")
	}
	if _, ok := parent.Get("b"); ok {
		t.Error("parent should not see child var")
	}
}

func TestChildOverridesParent(t *testing.T) {
	parent := New()
	parent.Set("x", 1)
	child := parent.Child()
	child.Set("x", 99)

	v, _ := child.Get("x")
	if v != 99 {
		t.Errorf("child x = %f, want 99", v)
	}
	pv, _ := parent.Get("x")
	if pv != 1 {
		t.Errorf("parent x = %f, want 1", pv)
	}
}

func TestFlatten(t *testing.T) {
	parent := New()
	parent.Set("a", 1)
	child := parent.Child()
	child.Set("b", 2)
	child.Set("a", 10)

	flat := child.Flatten()
	if flat["a"] != 10 {
		t.Errorf("flat[a] = %f, want 10", flat["a"])
	}
	if flat["b"] != 2 {
		t.Errorf("flat[b] = %f", flat["b"])
	}
}

func TestNames(t *testing.T) {
	e := New()
	e.Set("x", 1)
	names := e.Names()
	found := false
	for _, n := range names {
		if n == "x" {
			found = true
		}
	}
	if !found {
		t.Errorf("names = %v, missing x", names)
	}
}

func TestDepth(t *testing.T) {
	e := New()
	if e.Depth() != 0 {
		t.Errorf("root depth = %d", e.Depth())
	}
	c := e.Child()
	if c.Depth() != 1 {
		t.Errorf("child depth = %d", c.Depth())
	}
}

func TestLoadConstants(t *testing.T) {
	e := Empty()
	LoadConstants(e, "physics")
	if !e.Has("c") {
		t.Error("should have speed of light")
	}
	if !e.Has("g") {
		t.Error("should have gravity")
	}
}
