package builtin

// Registry holds all available built-in functions indexed by name.
type Registry struct {
	funcs map[string]Func
}

// NewRegistry creates a registry with all standard functions registered.
func NewRegistry() *Registry {
	r := &Registry{funcs: make(map[string]Func)}
	r.RegisterAll(MathFuncs())
	r.RegisterAll(StringFuncs())
	r.RegisterAll(LogicFuncs())
	return r
}

// EmptyRegistry creates a registry with no functions.
func EmptyRegistry() *Registry {
	return &Registry{funcs: make(map[string]Func)}
}

// Register adds a single function.
func (r *Registry) Register(name string, fn Func) {
	r.funcs[name] = fn
}

// RegisterAll adds multiple functions from a map.
func (r *Registry) RegisterAll(funcs map[string]Func) {
	for name, fn := range funcs {
		r.funcs[name] = fn
	}
}

// Get returns the function for the given name, or nil if not found.
func (r *Registry) Get(name string) Func {
	return r.funcs[name]
}

// Has returns true if a function with the given name is registered.
func (r *Registry) Has(name string) bool {
	_, ok := r.funcs[name]
	return ok
}

// Names returns all registered function names in no particular order.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.funcs))
	for name := range r.funcs {
		out = append(out, name)
	}
	return out
}

// Count returns the number of registered functions.
func (r *Registry) Count() int {
	return len(r.funcs)
}

// Remove deregisters a function by name.
func (r *Registry) Remove(name string) {
	delete(r.funcs, name)
}

// Clone creates a copy of the registry.
func (r *Registry) Clone() *Registry {
	cp := &Registry{funcs: make(map[string]Func, len(r.funcs))}
	for k, v := range r.funcs {
		cp.funcs[k] = v
	}
	return cp
}
