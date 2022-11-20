package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	value, ok := env.store[name]

	if !ok && env.outer != nil {
		value, ok = env.outer.Get(name)
	}

	return value, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value
	return value
}
