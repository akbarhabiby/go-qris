# 📦 QRIS Go Package

[![Go Reference](https://pkg.go.dev/badge/github.com/akbarhabiby/go-qris.svg)](https://pkg.go.dev/github.com/akbarhabiby/go-qris)


A modern and developer-friendly Go package for parsing, generating, and modifying **QRIS (Quick Response Code Indonesian Standard)** data, including support for static and dynamic QR codes, fee insertion, and QR image handling.

---

## ✨ Features

- ✅ Parse and serialize QRIS from raw string or image (PNG/JPG)
- ✅ Modify merchant info, amount, city, postal code, and fees
- ✅ Convert static QRIS to dynamic
- ✅ Support service fee (rupiah or percent)
- ✅ Export as PNG or decode from image
- ✅ Map QRIS to structured Go type (`QRISData`)
- ✅ No external C dependencies

---

## 📦 Installation

```bash
go get github.com/akbarhabiby/go-qris
````

---

## 🛠 Usage

### Parse from string

```go
raw := "0002010102112657..." // your QRIS raw string
q, err := qris.NewQRISFromString(raw)
if err != nil {
	log.Fatal(err)
}
fmt.Println("Merchant:", q.Get(qris.TAG_MERCHANT_NAME))
```

### Convert to dynamic QR with amount and service fee

```go
q.SetAmountWithOptions(qris.QRISAmountOptions{
	Amount:   10000,
	FeeType:  qris.QRISFeeRupiah,
	FeeValue: 2000,
})
fmt.Println(q.Serialize())
```

### Replace merchant info

```go
q.SetMerchantName("Toko Baru")
q.SetMerchantCityAndPostalCode("Jakarta", "12345")
```

---

## 🖼 QR Image Support

### Load from image

```go
q, err := qris.NewQRISFromImage("qris.png")
```

### Generate QR image (as bytes)

```go
imgBytes, err := q.GenerateQRISImage(256)
```

### Save as PNG

```go
err := q.SaveQRISAsImage("output.png", 256)
```

---

## 🧩 Map to Struct

```go
data := q.MapToStruct()
fmt.Printf("Amount: %s\n", data.TransactionAmount)
fmt.Printf("Merchant: %s\n", data.MerchantName)
```

---

## 🧪 Testing

```bash
# Run all tests with verbose output
go test ./... -v

# Run all benchmarks and show memory allocations
go test -bench=. -benchmem
```
