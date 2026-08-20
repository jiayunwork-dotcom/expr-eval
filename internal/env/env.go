// Package env provides a scoped variable environment for expression evaluation.
// It supports nested scopes (let-bindings), built-in constants (pi, e, tau),
// and variable registration with default values.
package env

import (
	"fmt"
	"math"
	"sort"
)

// Env holds variable bindings with support for nested scopes.
type Env struct {
	parent *Env
	vars   map[string]float64
}

// New creates a root environment with built-in constants.
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

// Empty creates an environment with no predefined bindings.
func Empty() *Env {
	return &Env{vars: map[string]float64{}}
}

// Child creates a new scope nested under this environment.
func (e *Env) Child() *Env {
	return &Env{parent: e, vars: map[string]float64{}}
}

// Set binds a variable in the current scope.
func (e *Env) Set(name string, val float64) {
	e.vars[name] = val
}

// SetAll binds multiple variables in the current scope.
func (e *Env) SetAll(vars map[string]float64) {
	for k, v := range vars {
		e.vars[k] = v
	}
}

// Get looks up a variable, searching parent scopes if not found locally.
func (e *Env) Get(name string) (float64, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return 0, false
}

// MustGet looks up a variable or returns an error.
func (e *Env) MustGet(name string) (float64, error) {
	v, ok := e.Get(name)
	if !ok {
		return 0, fmt.Errorf("undefined variable %q", name)
	}
	return v, nil
}

// Has returns true if the variable is defined in any scope.
func (e *Env) Has(name string) bool {
	_, ok := e.Get(name)
	return ok
}

// Delete removes a variable from the current scope only.
func (e *Env) Delete(name string) {
	delete(e.vars, name)
}

// Names returns all variable names visible from this scope.
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

// LocalNames returns variable names defined in this scope only (not parent).
func (e *Env) LocalNames() []string {
	names := make([]string, 0, len(e.vars))
	for k := range e.vars {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Depth returns the nesting depth of this scope (0 = root).
func (e *Env) Depth() int {
	d := 0
	for cur := e.parent; cur != nil; cur = cur.parent {
		d++
	}
	return d
}

// Flatten returns a flat map of all visible bindings (child overrides parent).
func (e *Env) Flatten() map[string]float64 {
	result := map[string]float64{}
	// collect from root to leaf, so leaf overrides
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

// Merge copies all bindings from other into this env (current scope).
func (e *Env) Merge(other *Env) {
	for k, v := range other.vars {
		e.vars[k] = v
	}
}

// Size returns the number of bindings in the current scope.
func (e *Env) Size() int {
	return len(e.vars)
}
