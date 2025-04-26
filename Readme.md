## Blockchain and Image Copyright Infrignment

The project is about a system that store images and in case if a similar images with some modification is uploaded it will able to detect used algorithm is whash

The project is deployed using Hyperledger fabric cli

## Prerequisite
- Ubantu
- Docker
- Fabric
- vscode

## Installation

- Open a terminal in minifab folder
- Run below command
    - sudo chmod -R 777 vars/
    - mkdir -p vars/chaincode/VIT-chain/go
    - cp -r ../Chaincode/* vars/chaincode/VIT-chain/go/
- Now, minifab network created and running to check type docker ps ,docker containers running will showed
- To deploy chaincode 
    - minifab ccup -n VIT-chain -l go -v 1.0 -d false -r false
- Now the hyperledger is running and chaincode is deployed 
- To register an image 
    - minifab invoke -n VIT-chain -p '"RegisterImage","<image_url_downloadable_form>","<image_owner_name>"'
- To check all images registered
    - minifab query -n VIT-chain -p '"GetAllImages"'
- To check all images registered by owner_name
    - minifab query -n VIT-chain -p '"GetImagesByOwner","<owner_name>"'

image_url_format- https://drive.google.com/uc?export=download&id=<image_id_of_gdrive>


## To safely shutdown the hyperledger network

Due to some issue the version i am using have issue in saving state of hyperledger network
- minifab down
- minifab cleanup

## Hyperledger configuration 
 Minifab/spec.yaml

In case of any issue take reference from screenshot.html or mail at rishabh480604@gmail.com 


