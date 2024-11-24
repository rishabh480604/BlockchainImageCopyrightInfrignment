#!/bin/sh

echo "Start  the Network"
echo "executed command : minifab netup -s couchdb -e true -i 2.4.8 -o photographer.image.com"
minifab netup -s couchdb -e true -i 2.4.8 -o photographer.image.com

sleep 5

echo "create the channel"
echo "executed command : minifab create -c autochannel"
minifab create -c autochannel

sleep 2

echo "Join the peers to the channel"
echo "executed command : minifab join -c autochannel"
minifab join -c autochannel

sleep 2

echo "Anchor update"
echo "executed command : minifab anchorupdate"
minifab anchorupdate

sleep 2

echo "Profile Generation"
echo "executed command : minifab profilegen -c autochannel"
minifab profilegen -c autochannel


echo "Provide all permission to vars"
echo "executed coomand : sudo chmod -R 777 vars/"
sudo chmod -R 777 vars/

sleep 2
echo "make chaincode/VIT-chain/go folder"
echo "executed coomand : mkdir -p vars/chaincode/VIT-chain/go"
mkdir -p vars/chaincode/VIT-chain/go

sleep 2
echo "cp chaincode files to vars created folder"
echo "executed coomand : cp -r ../Chaincode/* vars/chaincode/VIT-chain/go/"
cp -r ../Chaincode/* vars/chaincode/VIT-chain/go/

sleep 2

echo "build chaincode"
echo "executed coomand : minifab ccup -n VIT-chain -l go -v 1.0 -d false -r false"
minifab ccup -n VIT-chain -l go -v 1.0 -d false -r false


