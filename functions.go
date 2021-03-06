/*
All functions used to generate a loyalty card
*/
package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/codabar"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/twooffive"
	"github.com/nfnt/resize"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var width, height int = 470, 320
var (
	fontBytes []byte
	c         *freetype.Context
	f         *truetype.Font
)

func init() {
	// Read the font data.
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Dir(ex)

	fontBytes, err := ioutil.ReadFile(dir + "/FreeSans.ttf")
	if err != nil {
		log.Fatal(err)
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize freetype context
	c = freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)
}

func getBarCode() (barcode.Barcode, error) {
	var (
		codeEncoded barcode.Barcode
		err         error
	)

	// Create the barcode
	switch *codeType {
	case "codabar":
		codeEncoded, err = codabar.Encode(*cardNumber)
		if err != nil {
			return nil, err
		}
	case "code128":
		codeEncoded, err = code128.Encode(*cardNumber)
		if err != nil {
			return nil, err
		}
	case "code39":
		codeEncoded, err = code39.Encode(*cardNumber, false, false)
		if err != nil {
			return nil, err
		}
	case "code39FullAscii":
		codeEncoded, err = code39.Encode(*cardNumber, false, true)
		if err != nil {
			return nil, err
		}
	case "itf":
		codeEncoded, err = twooffive.Encode(*cardNumber, true)
		if err != nil {
			return nil, err
		}
	default:
		codeEncoded, err = ean.Encode(*cardNumber)
		if err != nil {
			return nil, err
		}
	}

	// Scale the barcode to 200x100 pixels if possible
	minWidth := codeEncoded.Bounds().Max.X - codeEncoded.Bounds().Min.X
	if 200 > minWidth {
		return barcode.Scale(codeEncoded, 200, 100)
	}
	return barcode.Scale(codeEncoded, minWidth, 100)
}

func getCodeImg() (image.Image, int, error) {
	// Format code
	var str string
	switch *codeType {
	case "codabar":
		str = (*cardNumber)[:1] + " " + (*cardNumber)[1:len(*cardNumber)-1] + " " + (*cardNumber)[len(*cardNumber)-1:]
	case "code128":
		str = *cardNumber
	case "code39", "code39FullAscii", "itf":
		str = (*cardNumber)[:]
	default:
		str = (*cardNumber)[:1] + " " + (*cardNumber)[1:7] + " " + (*cardNumber)[7:]
	}

	codeImg := image.NewRGBA(image.Rectangle{image.ZP, image.Pt(width, 100)})
	draw.Draw(codeImg, codeImg.Bounds(), image.White, image.ZP, draw.Src)
	c.SetFontSize(18)
	c.SetClip(codeImg.Bounds())
	c.SetDst(codeImg)

	// Draw the text.
	pt := freetype.Pt(0, int(c.PointToFixed(18)>>6))
	lastPoint, err := c.DrawString(str, pt)
	if err != nil {
		return nil, 0, err
	}

	textLength := fixed.Int26_6.Ceil(lastPoint.X)

	return codeImg, textLength, nil
}

func getCodeTypeImg() (image.Image, int, error) {
	// if *codeType is ean, override *codeType
	if *codeType == "ean" {
		if len(*cardNumber) == 12 || len(*cardNumber) == 13 {
			*codeType = "ean13"
		} else if len(*cardNumber) == 7 || len(*cardNumber) == 8 {
			*codeType = "ean8"
		} else {

			return nil, 0, errors.New("ean code length is " + string(strconv.Itoa(len(*cardNumber))) + ". It should be 8 or 13.")
		}
	}

	codeTypeImg := image.NewRGBA(image.Rectangle{image.ZP, image.Pt(width, 100)})
	draw.Draw(codeTypeImg, codeTypeImg.Bounds(), image.White, image.ZP, draw.Src)
	c.SetFontSize(14)
	c.SetClip(codeTypeImg.Bounds())
	c.SetDst(codeTypeImg)
	//c.SetSrc(image.NewUniform(color.RGBA{0, 128, 0, 255}))

	// Draw the text.
	pt := freetype.Pt(0, int(c.PointToFixed(14)>>6))
	lastPoint, err := c.DrawString("Code : "+strings.ToUpper(*codeType), pt)
	if err != nil {
		return nil, 0, err
	}

	textLength := fixed.Int26_6.Ceil(lastPoint.X)

	return codeTypeImg, textLength, nil
}

func getOwnerImg() (image.Image, int, error) {
	ownerImg := image.NewRGBA(image.Rectangle{image.ZP, image.Pt(width, 100)})
	draw.Draw(ownerImg, ownerImg.Bounds(), image.White, image.ZP, draw.Src)
	c.SetFontSize(14)
	c.SetClip(ownerImg.Bounds())
	c.SetDst(ownerImg)
	//c.SetSrc(image.NewUniform(color.RGBA{0, 128, 0, 255}))

	// Draw the text.
	pt := freetype.Pt(0, int(c.PointToFixed(14)>>6))
	lastPoint, err := c.DrawString(*cardOwner, pt)
	if err != nil {
		return nil, 0, err
	}

	textLength := fixed.Int26_6.Ceil(lastPoint.X)

	return ownerImg, textLength, nil
}

func getResizedLogo() (image.Image, error) {
	handler, err := os.Open(*shopLogo)
	if err != nil {
		return nil, err
	}
	defer handler.Close()

	logoImg, _, err := image.Decode(handler)
	if err != nil {
		return nil, err
	}

	logoImgResized := resize.Thumbnail(400, 50, logoImg, resize.NearestNeighbor)

	return logoImgResized, nil
}

func drawCard(logo, barCode, code, codeType, owner image.Image, codeLength, codeTypeLength, ownerLength int) image.Image {
	card := image.NewRGBA(image.Rectangle{image.ZP, image.Pt(width, height)})
	draw.Draw(card, card.Bounds(), image.White, image.ZP, draw.Src)
	draw.Draw(card, card.Bounds(), logo, image.Pt((int(logo.Bounds().Dx()-width)/2), -50), draw.Src)
	draw.Draw(card, card.Bounds(), barCode, image.Pt(int(barCode.Bounds().Dx()-width)/2, -122), draw.Src)
	draw.Draw(card, card.Bounds(), code, image.Pt(int(codeLength-width)/2, -222), draw.Src)
	draw.Draw(card, card.Bounds(), codeType, image.Pt(-50, -260), draw.Src)
	draw.Draw(card, card.Bounds(), owner, image.Pt(-width+ownerLength+50, -260), draw.Src)

	return card
}

func saveloyaltyCard(image image.Image) error {
	filename := *shopName + "_" + *cardOwner + ".jpg"
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)

	opt := jpeg.Options{90}
	err = jpeg.Encode(b, image, &opt)

	fmt.Println("Wrote " + filename + ".")
	return err
}
