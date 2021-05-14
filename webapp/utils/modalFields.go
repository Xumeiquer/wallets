package utils

import (
	"errors"
)

type ModalFieldStr struct {
	Value *string
}

func NewModalField(s string) ModalFieldStr {
	return ModalFieldStr{
		Value: &s,
	}
}

func (s ModalFieldStr) StringValue() string {
	if s.Value == nil {
		panic(errors.New("ModalFieldStr must not have a nil pointer"))
	}
	return *s.Value

}

func (s ModalFieldStr) SetStringValue(v string) {
	if s.Value == nil {
		panic(errors.New("ModalFieldStr must not have a nil pointer"))
	}
	*s.Value = v
}
