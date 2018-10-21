package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	validator "gopkg.in/go-playground/validator.v9"
)

const NotificationEventType EventType = "notification"

const (
	Received   NotificationType = "received"
	Packaging  NotificationType = "packaging"
	Packaged   NotificationType = "packaged"
	Delivering NotificationType = "delivering"
	Delivered  NotificationType = "delivered"
	Transfused NotificationType = "transfused"
)

var notificationsSequence = []NotificationType{
	Received, Packaging, Packaged,
	Delivering, Delivered, Transfused,
}

func getParentNotificationType(nt NotificationType) NotificationType {
	var previous NotificationType
	for _, snt := range notificationsSequence {
		if snt == nt {
			return previous
		}
		previous = snt
	}
	return ""
}

type NotificationType string

type NotificationEvent struct {
	PrescriptionHash []byte           `json:"prescription" validate:"required"`
	NotificationType NotificationType `json:"notification" validate:"required"`
}

func (e *NotificationEvent) Type() EventType {
	return NotificationEventType
}

func (e *NotificationEvent) Validate(bc *Blockchain) error {

	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		return err
	}

	block := bc.FindPrescriptionBlock(e.PrescriptionHash)
	if block == nil {
		return errors.New("prescription does not exist")
	}

	events := bc.FindPrescriptionNotificationEvents(block.Hash)
	for _, ee := range events {
		if e.NotificationType == ee.NotificationType {
			return errors.New("notification already exists")
		}
	}

	if e.NotificationType == Received {
		return nil
	}

	parentNotificationType := getParentNotificationType(e.NotificationType)

	for _, e := range events {
		if e.NotificationType == parentNotificationType {
			return nil
		}
	}

	return errors.New("invalid notification event")
}

func (e *NotificationEvent) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(NotificationEventType),
		e.PrescriptionHash,
		[]byte(e.NotificationType),
	}, []byte{})
}

func (e *NotificationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type             EventType        `json:"type"`
		PrescriptionHash string           `json:"prescription"`
		NotificationType NotificationType `json:"notification"`
	}{NotificationEventType, base58.Encode(e.PrescriptionHash), e.NotificationType})
}

func (e *NotificationEvent) UnmarshalJSON(data []byte) error {
	var rawNotification map[string]*json.RawMessage
	if err := json.Unmarshal(data, &rawNotification); err != nil {
		return err
	}
	if err := json.Unmarshal(*rawNotification["notification"], &e.NotificationType); err != nil {
		return err
	}
	var prescriptionHash string
	if err := json.Unmarshal(*rawNotification["prescription"], &prescriptionHash); err != nil {
		return err
	}
	e.PrescriptionHash = base58.Decode(prescriptionHash)
	return nil
}

func NewNotificationEvent(prescriptionHash []byte, notificationType NotificationType) Event {
	return &NotificationEvent{prescriptionHash, notificationType}
}
