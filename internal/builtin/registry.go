package builtin

type Registry struct {
	funcs map[string]Func
}

func NewRegistry() *Registry {
	r := &Registry{funcs: make(map[string]Func)}
	r.RegisterAll(MathFuncs())
	r.RegisterAll(StringFuncs())
	r.RegisterAll(LogicFuncs())
	return r
}

func EmptyRegistry() *Registry {
	return &Registry{funcs: make(map[string]Func)}
}

func (r *Registry) Register(name string, fn Func) {
	r.funcs[name] = fn
}

func (r *Registry) RegisterAll(funcs map[string]Func) {
	for name, fn := range funcs {
		r.funcs[name] = fn
	}
}

func (r *Registry) Get(name string) Func {
	return r.funcs[name]
}

func (r *Registry) Has(name string) bool {
	_, ok := r.funcs[name]
	return ok
}

func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.funcs))
	for name := range r.funcs {
		out = append(out, name)
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.funcs)
}

func (r *Registry) Remove(name string) {
	delete(r.funcs, name)
}

func (r *Registry) Clone() *Registry {
	cp := &Registry{funcs: make(map[string]Func, len(r.funcs))}
	for k, v := range r.funcs {
		cp.funcs[k] = v
	}
	return cp
}
