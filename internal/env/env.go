package env

import (
	"fmt"
	"math"
	"sort"
)

type Env struct {
	parent *Env
	vars   map[string]float64
}

func New() *Env {
	e := &Env{vars: map[string]float64{
		"pi":  math.Pi,
		"e":   math.E,
		"tau": 2 * math.Pi,
		"inf": math.Inf(1),
		"nan": math.NaN(),
	}}
	return e
}

func Empty() *Env {
	return &Env{vars: map[string]float64{}}
}

func (e *Env) Child() *Env {
	return &Env{parent: e, vars: map[string]float64{}}
}

func (e *Env) Set(name string, val float64) {
	e.vars[name] = val
}

func (e *Env) SetAll(vars map[string]float64) {
	for k, v := range vars {
		e.vars[k] = v
	}
}

func (e *Env) Get(name string) (float64, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return 0, false
}

func (e *Env) MustGet(name string) (float64, error) {
	v, ok := e.Get(name)
	if !ok {
		return 0, fmt.Errorf("undefined variable %q", name)
	}
	return v, nil
}

func (e *Env) Has(name string) bool {
	_, ok := e.Get(name)
	return ok
}

func (e *Env) Delete(name string) {
	delete(e.vars, name)
}

func (e *Env) Names() []string {
	seen := map[string]bool{}
	var names []string
	for cur := e; cur != nil; cur = cur.parent {
		for k := range cur.vars {
			if !seen[k] {
				seen[k] = true
				names = append(names, k)
			}
		}
	}
	sort.Strings(names)
	return names
}

func (e *Env) LocalNames() []string {
	names := make([]string, 0, len(e.vars))
	for k := range e.vars {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func (e *Env) Depth() int {
	d := 0
	for cur := e.parent; cur != nil; cur = cur.parent {
		d++
	}
	return d
}

func (e *Env) Flatten() map[string]float64 {
	result := map[string]float64{}
	var chain []*Env
	for cur := e; cur != nil; cur = cur.parent {
		chain = append(chain, cur)
	}
	for i := len(chain) - 1; i >= 0; i-- {
		for k, v := range chain[i].vars {
			result[k] = v
		}
	}
	return result
}

func (e *Env) Merge(other *Env) {
	for k, v := range other.vars {
		e.vars[k] = v
	}
}

func (e *Env) Size() int {
	return len(e.vars)
}
