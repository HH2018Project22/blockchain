package blockchain

import (
	"bytes"
	"encoding/json"
)

const Prescription EventType = "prescription"

type PrescriptionEvent struct {
	Patient *Patient `json:"patient"`
}

func (e *PrescriptionEvent) Type() EventType {
	return Prescription
}

func (e *PrescriptionEvent) Validate(bc *Blockchain) bool {
	return true
}

func (e *PrescriptionEvent) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(Prescription),
		[]byte(e.Patient.FirstName),
		[]byte(e.Patient.LastName),
		[]byte(e.Patient.UseName),
		[]byte(e.Patient.BirthDate.String()),
	}, []byte{})
}

func (e *PrescriptionEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    EventType `json:"typ"`
		Patient *Patient  `json:"patient"`
	}{Prescription, e.Patient})
}

func NewPrescriptionEvent(patient *Patient) Event {
	return &PrescriptionEvent{patient}
}
