package blockchain

import "bytes"

type Patient struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	UseName   string `json:"useName" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required"`
	Sex       Sex    `json:"sex" validate:"required"`
}

func (p *Patient) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(p.FirstName),
		[]byte(p.LastName),
		[]byte(p.UseName),
		[]byte(p.BirthDate),
		[]byte(p.Sex),
	}, []byte{})
}

type Sex string

const (
	SexFemale Sex = "female"
	SexMale   Sex = "male"
)

func NewPatient(firstName string, lastName string, useName string, birthDate string, sex Sex) *Patient {
	return &Patient{
		FirstName: firstName,
		LastName:  lastName,
		UseName:   useName,
		BirthDate: birthDate,
	}
}
