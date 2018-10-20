package blockchain

import (
	"bytes"
	"encoding/json"

	"gopkg.in/go-playground/validator.v9"
)

const PrescriptionEventType EventType = "prescription"

type PrescriptionEvent struct {
	Prescription *Prescription `json:"prescription" validate:"dive,required"`
}

func (e *PrescriptionEvent) Type() EventType {
	return PrescriptionEventType
}

func (e *PrescriptionEvent) Validate(bc *Blockchain) error {
	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		return err
	}
	return nil
}

func (e *PrescriptionEvent) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(PrescriptionEventType),
		e.Prescription.Hash(),
	}, []byte{})
}

func (e *PrescriptionEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         EventType     `json:"type"`
		Prescription *Prescription `json:"prescription"`
	}{PrescriptionEventType, e.Prescription})
}

func NewPrescriptionEvent(prescription *Prescription) Event {
	return &PrescriptionEvent{prescription}
}
