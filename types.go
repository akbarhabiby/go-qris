package qris

// QRISTag is the string type for tag identifiers in QRIS TLV format
type QRISTag string

// Core QRIS Tag Constants (EMVCo-compliant)
const (
	TagPayloadFormat     QRISTag = "00" // Payload Format Indicator
	TagPointOfInitiation QRISTag = "01" // Point of Initiation Method

	// Merchant Account Tags (range: 26–51, only a few named)
	TagMerchantAccount26 QRISTag = "26" // Merchant Account Info - 26
	TagMerchantAccount27 QRISTag = "27"
	TagMerchantAccount28 QRISTag = "28"
	TagMerchantAccount29 QRISTag = "29"
	TagMerchantAccount30 QRISTag = "30"
	TagMerchantAccount51 QRISTag = "51" // Last of merchant account tags

	TagMerchantCategory    QRISTag = "52" // MCC
	TagTransactionCurrency QRISTag = "53" // ISO 4217 (360 = IDR)
	TagTransactionAmount   QRISTag = "54" // Nominal
	TagConvenienceFee      QRISTag = "55" // Service fee (prefix tag)

	TagCountryCode  QRISTag = "58" // ISO 3166 (ID = Indonesia)
	TagMerchantName QRISTag = "59"
	TagMerchantCity QRISTag = "60"
	TagPostalCode   QRISTag = "61"

	TagAdditionalData QRISTag = "62" // e.g., bill_no, ref_no, etc.
	TagCRC            QRISTag = "63" // CRC16 (must be last)

	// Custom Service Fee Tag Prefixes (non-standard)
	TagFeeRupiah  QRISTag = "55020256" // "Rp" fee
	TagFeePercent QRISTag = "55020357" // "%"-based fee
)

// TagDescriptions maps known QRISTags to readable descriptions.
var TagDescriptions = map[QRISTag]string{
	TagPayloadFormat:       "Payload Format Indicator",
	TagPointOfInitiation:   "Point of Initiation Method",
	TagMerchantAccount26:   "Merchant Account Info (26)",
	TagMerchantCategory:    "Merchant Category Code (MCC)",
	TagTransactionCurrency: "Transaction Currency (ISO 4217)",
	TagTransactionAmount:   "Transaction Amount",
	TagConvenienceFee:      "Convenience Fee Indicator",
	TagCountryCode:         "Country Code (ISO 3166)",
	TagMerchantName:        "Merchant Name",
	TagMerchantCity:        "Merchant City",
	TagPostalCode:          "Postal Code",
	TagAdditionalData:      "Additional Data Field Template",
	TagCRC:                 "CRC-16 Checksum",
	TagFeeRupiah:           "Service Fee (Rupiah)",
	TagFeePercent:          "Service Fee (Percent)",
}

// POIMType defines values for Point of Initiation Method (Tag "01")
type POIMType string

const (
	POIMStatic  POIMType = "11" // Static QR (amount entered manually)
	POIMDynamic POIMType = "12" // Dynamic QR (amount embedded)
)

var POIMDescriptions = map[POIMType]string{
	POIMStatic:  "Static QRIS",
	POIMDynamic: "Dynamic QRIS",
}

type QRISFeeType int

const (
	QRISFeeNone QRISFeeType = iota
	QRISFeeRupiah
	QRISFeePercent
)

type QRISAmountOptions struct {
	Amount   int
	FeeType  QRISFeeType
	FeeValue float64
}
