package main

import (
	"flag"
	"log"

	"github.com/HH2018Project22/bloodcoin/blockchain"
	"github.com/btcsuite/btcutil/base58"
)

var (
	notificationCommand = flag.NewFlagSet("notification", flag.ExitOnError)
	prescriptionHash    string
	notificationType    string
	firstName           string
	lastName            string
	service             string
)

func init() {
	notificationCommand.StringVar(&prescriptionHash, "prescription", prescriptionHash, "Hash de la prescription")
	notificationCommand.StringVar(&notificationType, "type", notificationType, "Type de la notification")
	notificationCommand.StringVar(&firstName, "firstName", firstName, "Prénom")
	notificationCommand.StringVar(&lastName, "lastName", lastName, "Nom")
	notificationCommand.StringVar(&service, "service", service, "Service")
}

func doNotification(args []string) {

	if err := notificationCommand.Parse(args); err != nil {
		panic(err)
	}

	bc := getBlockchain()

	log.Println("adding notification")

	prescriptionHashData := base58.Decode(prescriptionHash)

	notification := blockchain.NewNotificationEvent(
		prescriptionHashData,
		blockchain.NotificationType(notificationType),
		blockchain.NewOperator(firstName, lastName, service),
	)

	if _, err := bc.AddEvent(notification); err != nil {
		panic(err)
	}

	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		panic(err)
	}

}
