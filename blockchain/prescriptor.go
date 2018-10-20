package blockchain

import (
	"bytes"
)

type Prescriptor struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Service   string `json:"service" validate:"required"`
}

func (p *Prescriptor) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(p.FirstName),
		[]byte(p.LastName),
		[]byte(p.Service),
	}, []byte{})
}

func NewPrescriptor(firstName string, lastName string, service string) *Prescriptor {
	return &Prescriptor{firstName, lastName, service}
}
