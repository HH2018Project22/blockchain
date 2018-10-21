#!/bin/bash

# Wipe des données existantes
killall modd
killall bloodcoin-server
rm *.db
./bloodcoin-cli dump
cp bloodcoin.db bloodcoin-server.db

# Démarrage app
./bloodcoin-server &
sleep 2

peer=http://localhost:3000/blocks/new

# Insertion d'un jeu de données
prescription=$(cat sample/valid_prescription.json)
block=$(./bloodcoin-cli -peer $peer prescription -data "$prescription" -quiet)
./bloodcoin-cli -peer $peer notification -prescription $block -type "received" -firstName "Jean" -lastName "Dupont" -service "CHU Dijon - Réanimation"
./bloodcoin-cli -peer $peer notification -prescription $block -type "packaging" -firstName "Corrine" -lastName "Touzet" -service "EFS Dijon - Stocks PSL"
./bloodcoin-cli -peer $peer notification -prescription $block -type "packaged" -firstName "Corrine" -lastName "Touzet" -service "EFS Dijon - Stocks PSL"
./bloodcoin-cli -peer $peer notification -prescription $block -type "delivering" -firstName "Speedy" -lastName "Gonzales" -service "Livraison 3000 Dijon"
./bloodcoin-cli -peer $peer notification -prescription $block -type "delivered" -firstName "Speedy" -lastName "Gonzales" -service "Livraison 3000 Dijon"
./bloodcoin-cli -peer $peer notification -prescription $block -type "transfused" -firstName "Michel" -lastName "Plancard" -service "CHU Dijon - Chirurgie"

prescription=$(cat sample/prescription_low.json)
block=$(./bloodcoin-cli -peer $peer prescription -data "$prescription" -quiet)
./bloodcoin-cli -peer $peer notification -prescription $block -type "received" -firstName "Jean" -lastName "Dupont" -service "CHU Dijon - Réanimation"
./bloodcoin-cli -peer $peer notification -prescription $block -type "packaging" -firstName "Corrine" -lastName "Touzet" -service "EFS Dijon - Stocks PSL"
./bloodcoin-cli -peer $peer notification -prescription $block -type "packaged" -firstName "Corrine" -lastName "Touzet" -service "EFS Dijon - Stocks PSL"

prescription=$(cat sample/prescription_high.json)
block=$(./bloodcoin-cli -peer $peer prescription -data "$prescription" -quiet)
./bloodcoin-cli -peer $peer notification -prescription $block -type "received" -firstName "Pierre" -lastName "Ponce" -service "CHU Besançon - Hématologie"
./bloodcoin-cli -peer $peer notification -prescription $block -type "packaging" -firstName "Martin" -lastName "Bouygues" -service "EFS Besançon - Stocks PSL"

prescription=$(cat sample/prescription_emergency.json)
block=$(./bloodcoin-cli -peer $peer prescription -data "$prescription" -quiet)
./bloodcoin-cli -peer $peer notification -prescription $block -type "received" -firstName "Pierre" -lastName "Ponce" -service "CHU Besançon - Hématologie"
