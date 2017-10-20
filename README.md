# loyaltycard
Generate ean13 loyalty cards

## Why this software
You don't want have all yours loyalty cards in your purse and you don't want use a proprietary software (privacy, unwanted advertisements...) to store them in your smartphone.

One solution is to create your own digital loyalty cards with *loyaltycard* software and store them as a picture in your smartphone.


## Installation
```
go get github.com/Geek2France/loyaltycard
cd ${GOPATH}/src/github.com/Geek2France/loyaltycard
go install
cp FreeSans.ttf ${GOPATH}/bin
```
The file FreeSans.ttf and the executable file loyaltycard should always be in the same directory.

## Usage
You should set the following parameters :
* `-cardNumber`: Ean13 loyalty card code. Example:  "1234567890128"
* `-cardOwner` : Loyalty card owner. Example: "John DOE"
* `-shopName`  : Shop that provided the loyalty card. This will be the name of loyalty card image. Example: "Décathlon"
* `-shopLogo`  : Image representating the shop. Png and jpeg image formats are supported. Example: "/data/Logos/decathlonlogo.png"
```
${GOPATH}/bin/loyaltycard -cardNumber "1234567890128" -cardOwner "John DOE" -shopName "Décathlon" -shopLogo "/data/Logos/decathlonlogo.png"
```
<br />
This will generate a loyalty card image named D&eacute;cathlon.jpg:
<table>
<tr>
<th>
<img src="https://github.com/Geek2France/loyaltycard/blob/master/blob/master/img/D%C3%A9cathlon.jpg" alt="digital loyalty card example">
</th>
</tr>
</table>
<br/>
Now transfer this image to your smartphone and leave your loyalty card at home :grinning:
