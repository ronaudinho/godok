package tests

import (
	"errors"
	"fmt"
)

var ErrNoName = errors.New("no name")

type Simple struct {
	Name string
}

func SimpleF() {
	fmt.Println("simple")
}

func (s *Simple) Test() {
	fmt.Println(s.Name)
}

func (s *Simple) String() string {
	return s.Name
}

func (s *Simple) HasName() bool {
	return s.Name != ""
}

func (s *Simple) Replace(str string) error {
	if !s.HasName() {
		return ErrNoName
	}
	s.Name = str
	return nil
}
