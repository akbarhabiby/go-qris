package qris

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	qrcgen "github.com/skip2/go-qrcode"
)

type QRIS struct {
	Raw  string
	TLVs []TLV
}

// NewQRISFromString parses a QRIS string into a QRIS struct.
func NewQRISFromString(raw string) (*QRIS, error) {
	tlvs, err := ParseTLV(raw)
	if err != nil {
		return nil, err
	}
	return &QRIS{Raw: raw, TLVs: tlvs}, nil
}

// NewQRISFromImage parses a QRIS image (PNG/JPG) into a QRIS struct.
func NewQRISFromImage(path string) (*QRIS, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to binary bitmap: %v", err)
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode QR code: %v", err)
	}

	raw := strings.TrimSpace(result.String())
	return NewQRISFromString(raw)
}

func (q *QRIS) Serialize() string {
	return SerializeTLV(q.TLVs)
}

func (q *QRIS) MapToStruct() *QRISData {
	data := &QRISData{
		MerchantAccounts: make(map[string]string),
		AdditionalData:   make(map[string]string),
		Unmapped:         make(map[string]string),
	}

	for _, tlv := range q.TLVs {
		tag := QRISTag(tlv.Tag)
		switch {
		case tag == TagPayloadFormat:
			data.PayloadFormat = tlv.Value
		case tag == TagPointOfInitiation:
			data.PointOfInitiation = tlv.Value
		case tag >= "26" && tag <= "51":
			data.MerchantAccounts[tlv.Tag] = tlv.Value
		case tag == TagMerchantCategory:
			data.MerchantCategoryCode = tlv.Value
		case tag == TagTransactionCurrency:
			data.TransactionCurrency = tlv.Value
		case tag == TagTransactionAmount:
			data.TransactionAmount = tlv.Value
		case strings.HasPrefix(tlv.Tag, "55"):
			data.TipOrConvenience = tlv.Value
			switch tlv.Tag {
			case string(TagFeeRupiah):
				data.FeeType = "rupiah"
				data.FeeValue = tlv.Value
			case string(TagFeePercent):
				data.FeeType = "percent"
				data.FeeValue = tlv.Value
			}
		case tag == TagCountryCode:
			data.CountryCode = tlv.Value
		case tag == TagMerchantName:
			data.MerchantName = tlv.Value
		case tag == TagMerchantCity:
			data.MerchantCity = tlv.Value
		case tag == TagPostalCode:
			data.PostalCode = tlv.Value
		case tag == TagAdditionalData:
			// Parse sub-TLVs from tag 62
			if subTLVs, err := ParseTLV(tlv.Value); err == nil {
				for _, sub := range subTLVs {
					data.AdditionalData[sub.Tag] = sub.Value
				}
			}
		case tag == TagCRC:
			data.CRC = tlv.Value
		default:
			data.Unmapped[tlv.Tag] = tlv.Value
		}
	}

	return data
}

func (q *QRIS) Get(tag QRISTag) string {
	for _, tlv := range q.TLVs {
		if tlv.Tag == string(tag) {
			return tlv.Value
		}
	}
	return ""
}

func (q *QRIS) Replace(tag string, newValue string) {
	q.TLVs = ReplaceTLVValue(q.TLVs, tag, newValue)
}

func (q *QRIS) IsDynamic() bool {
	return q.Get(TagPointOfInitiation) == string(POIMDynamic)
}

func (q *QRIS) IsStatic() bool {
	return q.Get(TagPointOfInitiation) == string(POIMStatic)
}

// SetAmountWithOptions adds amount and optional fee to QRIS and makes it dynamic.
func (q *QRIS) SetAmountWithOptions(opts QRISAmountOptions) {
	if opts.Amount <= 0 {
		return
	}

	q.TLVs = ReplaceTLVValue(q.TLVs, string(TagPointOfInitiation), string(POIMDynamic))
	q.TLVs = RemoveTLV(q.TLVs, string(TagTransactionAmount))
	q.TLVs = RemoveTLVPrefix(q.TLVs, "55")

	amountStr := strconv.Itoa(opts.Amount)
	amountTLV := TLV{
		Tag:   string(TagTransactionAmount),
		Len:   len(amountStr),
		Value: amountStr,
	}

	var feeTLV *TLV
	switch opts.FeeType {
	case QRISFeeRupiah:
		feeStr := strconv.Itoa(int(opts.FeeValue))
		feeTLV = &TLV{
			Tag:   string(TagFeeRupiah),
			Len:   len(feeStr),
			Value: feeStr,
		}
	case QRISFeePercent:
		feeStr := fmt.Sprintf("%.2f", opts.FeeValue)
		feeTLV = &TLV{
			Tag:   string(TagFeePercent),
			Len:   len(feeStr),
			Value: feeStr,
		}
	}

	// Insert amount & fee before TagCountryCode (Tag 58)
	result := []TLV{}
	inserted := false
	for _, tlv := range q.TLVs {
		if tlv.Tag == string(TagCountryCode) && !inserted {
			result = append(result, amountTLV)
			if feeTLV != nil {
				result = append(result, *feeTLV)
			}
			inserted = true
		}
		result = append(result, tlv)
	}
	q.TLVs = UpdateCRC(result)
}

func (q *QRIS) SetMerchantName(name string) {
	q.TLVs = ReplaceTLVValue(q.TLVs, string(TagMerchantName), name)
}

func (q *QRIS) SetMerchantCityAndPostalCode(city string, postalCode string) {
	q.TLVs = ReplaceTLVValue(q.TLVs, string(TagMerchantCity), city)
	q.TLVs = ReplaceTLVValue(q.TLVs, string(TagPostalCode), postalCode)
}

func (q *QRIS) GenerateQRISImage(size int) ([]byte, error) {
	if q == nil {
		return nil, fmt.Errorf("QRIS is nil")
	}
	data := q.Serialize()
	return qrcgen.Encode(data, qrcgen.Medium, size)
}

func (q *QRIS) SaveQRISAsImage(path string, size int) error {
	if q == nil {
		return fmt.Errorf("QRIS is nil")
	}
	data := q.Serialize()
	return qrcgen.WriteFile(data, qrcgen.Medium, size, path)
}

func (q *QRIS) PrintToTerminal() error {
	if q == nil {
		return fmt.Errorf("QRIS is nil")
	}
	qr, err := qrcgen.New(q.Serialize(), qrcgen.Medium)
	if err != nil {
		return err
	}
	fmt.Println(qr.ToSmallString(false))
	return nil
}
