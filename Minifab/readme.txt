Installation


./network.sh
sudo chmod -R 777 vars/
mkdir -p vars/chaincode/VIT-chain/go
cp -r ../Chaincode/* vars/chaincode/VIT-chain/go/
minifab ccup -n VIT-chain -l go -v 1.0 -d false -r false



## logo img

minifab invoke -n VIT-chain -p '"RegisterImage","https://drive.google.com/uc?export=download&id=1aP-udmiRwO9HROs1bJ_4k-F7iuNA0-7F","Nihar"'
## other img
minifab invoke -n VIT-chain -p '"RegisterImage","https://drive.google.com/uc?export=download&id=1zS86aPa74Zv22NLeOyjxuOP1fidIVkfs","Aniket"'

## logo copyright img
minifab invoke -n VIT-chain -p '"RegisterImage","https://drive.google.com/uc?export=download&id=1-hmUnOfn00Qj9vhIbfMeRF7D6gWDhmle","Aniket"'

## Return all images
minifab query -n VIT-chain -p '"GetAllImages"'

## Return getimages by owner
minifab query -n VIT-chain -p '"GetImagesByOwner","Nihar"'





