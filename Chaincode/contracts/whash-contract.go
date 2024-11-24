package contracts

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/loupax/imagehash"
)

// ImageRegistry smart contract
type ImageRegistry struct {
	contractapi.Contract
}

// Image struct stores hash and owner
type Image struct {
	Hash  string `json:"hash"`
	Owner string `json:"owner"`
}

// downloadImageFromURL downloads an image from a URL and saves it to a temporary path
func downloadImageFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "downloaded_image_*.jpg")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}

	return tmpFile.Name(), nil
}

// RegisterImage registers a new image's hash on the ledger by downloading from a URL
func (c *ImageRegistry) RegisterImage(ctx contractapi.TransactionContextInterface, imageUrl string, owner string) (string, error) {
	imagePath, err := downloadImageFromURL(imageUrl)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %v", err)
	}
	defer os.Remove(imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("error opening downloaded image: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	var hashObj imagehash.Imagehash
	if err := hashObj.Whash(img, 64); err != nil { // 64x64 hash size
		return "", fmt.Errorf("error calculating wHash: %v", err)
	}
	hash := hashObj.String()

	existingImages, err := c.GetAllImages(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get existing images: %v", err)
	}

	threshold := 100
	for _, existingHash := range existingImages {
		var existingHashObj imagehash.Imagehash
		if err := existingHashObj.FromString(existingHash); err != nil {
			continue
		}

		distance, err := hashObj.Distance(existingHashObj)
		if err != nil {
			return "", fmt.Errorf("error calculating distance: %v", err)
		}

		if distance <= threshold {
			return fmt.Sprintf("Duplicate image detected. Owner: %s, Hamming Distance: %d",
				owner, distance), nil
		}
	}

	image := Image{
		Hash:  hash,
		Owner: owner,
	}
	imageJSON, err := json.Marshal(image)
	if err != nil {
		return "", fmt.Errorf("error marshalling image data: %v", err)
	}

	err = ctx.GetStub().PutState(hash, imageJSON)
	if err != nil {
		return "", fmt.Errorf("failed to register image: %v", err)
	}

	return fmt.Sprintf("Image registered successfully. Owner: %s", owner), nil
}

// GetAllImages retrieves all the registered image hashes from the world state
func (c *ImageRegistry) GetAllImages(ctx contractapi.TransactionContextInterface) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var hashes []string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var image Image
		err = json.Unmarshal(queryResponse.Value, &image)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, image.Hash)
	}

	return hashes, nil
}

// GetImageDetailByOwner retrieves hashes of images associated with a specific owner
func (c *ImageRegistry) GetImageDetailByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var hashes []string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var image Image
		err = json.Unmarshal(queryResponse.Value, &image)
		if err != nil {
			return nil, err
		}

		if image.Owner == owner {
			hashes = append(hashes, image.Hash)
		}
	}

	return hashes, nil
}

// VerifyImage checks if an image already exists in the ledger by comparing its hash
func (c *ImageRegistry) VerifyImage(ctx contractapi.TransactionContextInterface, imageUrl string) (string, error) {
	imagePath, err := downloadImageFromURL(imageUrl)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %v", err)
	}
	defer os.Remove(imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("error opening downloaded image: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Compute the Wavelet Hash (wHash) of the input image
	var hashObj imagehash.Imagehash
	if err := hashObj.Whash(img, 64); err != nil {
		return "", fmt.Errorf("error calculating wHash: %v", err)
	}

	// Retrieve all existing images from the ledger
	existingImages, err := c.GetAllImages(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get existing images: %v", err)
	}

	threshold := 100 // Updated similarity threshold
	for _, existingHash := range existingImages {
		var existingHashObj imagehash.Imagehash
		if err := existingHashObj.FromString(existingHash); err != nil {
			continue
		}

		// Calculate Hamming distance
		distance, err := hashObj.Distance(existingHashObj)
		if err != nil {
			return "", fmt.Errorf("error calculating distance: %v", err)
		}

		// If similarity is found
		if distance <= threshold {
			return fmt.Sprintf("Image found. Owner: %s, Hamming Distance: %d",
				existingHash, distance), nil
		}
	}

	// No matching image found
	return "No matching image found.", nil
}

