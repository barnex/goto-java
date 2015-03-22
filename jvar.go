package gotojava

type JVar struct {
	Name string
}

func (j *JVar) jExpr() {}

func NewJVar(name string) *JVar {
	return &JVar{Name: name}
}
