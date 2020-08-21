package tests

import (
	"fmt"
)

type TestI interface {
	Test()
}

type Test struct {
	Name string
}

func NewTest() *Test {
	return &Test{Name: "tests"}
}

func (t *Test) Test() {
	fmt.Println(t.Name)
}
