package blockchain

import "time"

type Patient struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UseName   string    `json:"useName"`
	BirthDate time.Time `json:"birthDate"`
}

func NewPatient(firstName string, lastName string, useName string, birthDate time.Time) *Patient {
	return &Patient{
		FirstName: firstName,
		LastName:  lastName,
		UseName:   useName,
		BirthDate: birthDate,
	}
}
