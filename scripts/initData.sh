#!/bin/bash

# Wipe des données existantes
rm *.db
./bloodcoin-cli dump
cp bloodcoin.db bloodcoin-server.db

# Démarrage app
./bloodcoin-server &
sleep 2

# Insertion d'un jeu de données
prescription=$(cat sample/valid_prescription.json)
echo $prescription
block=$(./bloodcoin-cli -peer http://localhost:3000/blocks/new prescription -data "$prescription" | sed 's/^.*Block: \(.*\){34}/\1/')
echo $block
#./bloodcoin-cli prescription -data "$(cat sample/valid_prescription.json)"
