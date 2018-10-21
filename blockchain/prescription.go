package blockchain

import (
	"bytes"
)

const (
	UrgencyLow       = "low"
	UrgencyHigh      = "high"
	UrgencyEmergency = "emergency"
)

type HashedPrescription struct {
	Hash         string        `json:"hash" validate:"required"`
	Prescription *Prescription `json:"prescription" validate:"dive,required"`
}

type Prescription struct {
	Patient     *Patient  `json:"patient" validate:"dive,required"`
	Prescriptor *Operator `json:"prescriptor" validate:"dive,required"`
	Order       *Order    `json:"order" validate:"dive,required"`
	Urgency     string    `json:"urgency" validate:"required"`
}

func (p *Prescription) Hash() []byte {
	return bytes.Join([][]byte{
		p.Patient.Hash(),
		p.Prescriptor.Hash(),
		p.Order.Hash(),
		[]byte(p.Urgency),
	}, []byte{})
}

func NewPrescription(patient *Patient) *Prescription {
	return &Prescription{
		Patient: patient,
	}
}
