package qris

type QRISData struct {
	// Required Fields
	PayloadFormat     string `json:"payload_format"`      // Tag 00
	PointOfInitiation string `json:"point_of_initiation"` // Tag 01

	// Merchant Account Information (Tag 26–51)
	MerchantAccounts map[string]string `json:"merchant_accounts"` // ex: {"26": "...", "51": "..."}

	// Merchant & Transaction Info
	MerchantCategoryCode string `json:"merchant_category_code"` // Tag 52
	TransactionCurrency  string `json:"transaction_currency"`   // Tag 53
	TransactionAmount    string `json:"transaction_amount"`     // Tag 54
	TipOrConvenience     string `json:"tip_or_convenience"`     // Tag 55 or 55020256/55020357

	// Location Info
	CountryCode  string `json:"country_code"`  // Tag 58
	MerchantName string `json:"merchant_name"` // Tag 59
	MerchantCity string `json:"merchant_city"` // Tag 60
	PostalCode   string `json:"postal_code"`   // Tag 61

	// Optional additional fields
	AdditionalData map[string]string `json:"additional_data"` // Tag 62
	Unmapped       map[string]string `json:"unmapped"`        // Any unknown tag (fallback)

	// Technical Fields
	CRC string `json:"crc"` // Tag 63

	// Derived info from fee tag (if 55020256 or 55020357 is found)
	FeeType  string `json:"fee_type,omitempty"`  // "rupiah", "percent"
	FeeValue string `json:"fee_value,omitempty"` // "2000", "2.5", etc.
}
