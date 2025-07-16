package main

import (
	"fmt"
	"os"

	"github.com/akbarhabiby/go-qris"
)

func main() {
	// 1. Sample QRIS file
	raw := "./testdata/qris.png"

	// 2. Parse QRIS string
	q, err := qris.NewQRISFromImage(raw)
	if err != nil {
		fmt.Println("Failed to parse QRIS:", err)
		os.Exit(1)
	}

	// 3. Update merchant info
	q.SetMerchantName("Bitcoin Bu Sugeng")
	q.SetMerchantCityAndPostalCode("Kab. Demak", "59567")

	// 4. Set nominal amount and fee
	q.SetAmountWithOptions(qris.QRISAmountOptions{
		Amount:   15000,
		FeeType:  qris.QRISFeeRupiah,
		FeeValue: 2000,
	})

	// 5. Print serialized QRIS string
	fmt.Println("QRIS Payload:")
	fmt.Println(q.Serialize())

	// 6. Print QR code to terminal
	fmt.Println("\nScan this QR (Terminal View):")
	_ = q.PrintToTerminal()

	// 7. Save QR code as PNG
	err = q.SaveQRISAsImage("qris.png", 300)
	if err != nil {
		fmt.Println("Failed to save QR image:", err)
	} else {
		fmt.Println("Saved QRIS to qris.png ✅")
	}
}
