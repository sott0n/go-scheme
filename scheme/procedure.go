// Procedure is a type for scheme procedure, which is expressed
// by lambda syntax form, like (lambda (x) x)
// when procedure has free variable, free variable must be binded when
// procedure is generated.
// So all Procedures have variable binding by Environment type (when there is
// no free variable, Procedure has Environmnet which is empty).

package scheme

type Procedure struct {
	ObjectBase
}

func NewProcedure(func(Object) Object) *Procedure {
	return &Procedure{}
}
