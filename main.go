/*
This package generates ean13 loyalty cards fidelity png images.
You can then transfer them to your smartphone and take them away from your purse.
*/
package main

import (
	"flag"
	"log"
)

var (
	cardNumber = flag.String("cardNumber", "1234567890123", "ean13 loyalty card number")
	cardOwner  = flag.String("cardOwner", "John Doe", "loyalty card owner")
	shopName   = flag.String("shopName", "Auchan", "shop name")
	shopLogo   = flag.String("shopLogo", "Auchan_logo.jpg", "shop jpeg logo")
)

func main() {
	flag.Parse()

	// Generate barcode image
	barCodeImg, err := getBarCode(*cardNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Transform card number into image
	codeImg, codeImgLength, err := getCodeImg(*cardNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Transform owner string into image
	ownerImg, ownerImgLength, err := getOwnerImg(*cardOwner)
	if err != nil {
		log.Fatal(err)
	}

	// Resize shop logo
	logoImgResized, err := getResizedLogo(*shopLogo)
	if err != nil {
		log.Fatal(err)
	}

	// Draw loyalty card
	loyaltyCard := drawCard(logoImgResized, barCodeImg, codeImg, ownerImg, codeImgLength, ownerImgLength)

	// Save that loyalty card image to disk.
	if err = saveloyaltyCard(loyaltyCard, *shopName); err != nil {
		log.Fatal(err)
	}
}
