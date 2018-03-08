// Environment has variable bindings.
// Interpreter has one Environment global variable for top-level environment.
// And each let block and procedure has Environment to hold its scope's variable binding.

package scheme

type Environment struct {
	ObjectBase
}

func newEnvironment() *Environment {
	return &Environment{}
}
