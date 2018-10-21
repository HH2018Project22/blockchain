package blockchain

import (
	"bytes"
)

type Operator struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Service   string `json:"service" validate:"required"`
}

func (o *Operator) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(o.FirstName),
		[]byte(o.LastName),
		[]byte(o.Service),
	}, []byte{})
}

func NewOperator(firstName string, lastName string, service string) *Operator {
	return &Operator{firstName, lastName, service}
}
