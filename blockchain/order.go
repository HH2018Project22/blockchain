package blockchain

import (
	"bytes"
	"fmt"
	"time"
)

type Order struct {
	ID                  string    `json:"id" validate:"required"`
	Amount              int       `json:"amount" validate:"required,min=1"`
	BloodType           string    `json:"bloodType" validate:"required"`
	TransfusionProtocol string    `json:"transfusionProtocol" validate:"required"`
	TransfusionTime     time.Time `json:"transfusionTime" validate:"required"`
}

func (o *Order) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(o.ID),
		[]byte(fmt.Sprintf("%d", o.Amount)),
		[]byte(o.BloodType),
		[]byte(o.TransfusionProtocol),
		[]byte(o.TransfusionTime.String()),
	}, []byte{})
}

func NewOrder(id string, amount int, bloodType string, transfusionProtocol string, transfusionTime time.Time) *Order {
	return &Order{id, amount, bloodType, transfusionProtocol, transfusionTime}
}
