package models

type Student struct {
	Name string
	Age  int
}

func NewStudent(n string, a int) *Student {
	return &Student{
		Name: n,
		Age:  a,
	}
}
