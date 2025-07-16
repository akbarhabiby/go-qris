package qris

import (
	"os"
	"testing"
)

const sampleQRISRaw = "00020101021126570011ID.DANA.WWW037491823004928374019283740192837401928UKC51440014ID.CO.QRIS.WWW84017629301574892036417UKC5204481453033605802ID5908Toko 8166013Jakarta Pusat610510330630468FE" // Not Scannable

func TestNewQRISFromString(t *testing.T) {
	q, err := NewQRISFromString(sampleQRISRaw)
	if err != nil {
		t.Fatalf("Failed to parse QRIS from string: %v", err)
	}
	if !q.IsStatic() || q.IsDynamic() {
		t.Error("Expected static QRIS")
	}
}

func TestGet(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	got := q.Get(TagMerchantName)
	if got != "Toko 816" {
		t.Errorf("Expected 'Toko 816', got '%s'", got)
	}
}

func TestReplace(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	q.Replace(string(TagMerchantName), "WarungA")
	if q.Get(TagMerchantName) != "WarungA" {
		t.Error("Replace failed")
	}
}

func TestSetMerchantCityAndPostalCode(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	q.SetMerchantCityAndPostalCode("Jakarta", "12345")
	if q.Get(TagMerchantCity) != "Jakarta" || q.Get(TagPostalCode) != "12345" {
		t.Error("Failed to update city/postal")
	}
}

func TestSetAmountWithOptions_Rupiah(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	q.SetAmountWithOptions(QRISAmountOptions{
		Amount:   15000,
		FeeType:  QRISFeeRupiah,
		FeeValue: 2000,
	})

	out := q.Serialize()

	if out[:12] != "000201010212" {
		t.Error("QRIS not converted to dynamic")
	}
	if !contains(out, "540515000") {
		t.Error("Amount tag 54 not present or incorrect")
	}
	if !contains(out, string(TagFeeRupiah)) {
		t.Error("Fee tag for Rupiah not present")
	}
}

func TestSetAmountWithOptions_Percent(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	q.SetAmountWithOptions(QRISAmountOptions{
		Amount:   5000,
		FeeType:  QRISFeePercent,
		FeeValue: 2.5,
	})

	out := q.Serialize()
	if !contains(out, string(TagFeePercent)) {
		t.Error("Fee tag for Percent not present")
	}
}

func TestGenerateQRISImage(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	data, err := q.GenerateQRISImage(300)
	if err != nil {
		t.Fatalf("Failed to generate QR image: %v", err)
	}
	if len(data) < 100 {
		t.Error("QR image is too small")
	}
}

func TestSaveQRISAsImage(t *testing.T) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	path := "test-qris.png"

	err := q.SaveQRISAsImage(path, 300)
	if err != nil {
		t.Fatalf("Save image failed: %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		t.Error("Saved image file is missing or empty")
	}
}

func TestPrintToTerminal(t *testing.T) {
	q, err := NewQRISFromString(sampleQRISRaw)
	if err != nil {
		t.Fatalf("Failed to create QRIS: %v", err)
	}

	// Suppress terminal output
	nullOut, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("Failed to open os.DevNull: %v", err)
	}
	defer nullOut.Close()

	originalStdout := os.Stdout
	os.Stdout = nullOut
	defer func() { os.Stdout = originalStdout }()

	err = q.PrintToTerminal()
	if err != nil {
		t.Errorf("PrintToTerminal returned error: %v", err)
	}
}

func TestNewQRISFromImage(t *testing.T) {
	if _, err := os.Stat("testdata/qris.png"); os.IsNotExist(err) {
		t.Skip("Image test skipped (no testdata/qris.png)")
	}

	q, err := NewQRISFromImage("testdata/qris.png")
	if err != nil {
		t.Fatalf("Failed to decode QR image: %v", err)
	}
	if len(q.TLVs) == 0 {
		t.Error("TLV not parsed from image")
	}
}

func BenchmarkParseQRIS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewQRISFromString(sampleQRISRaw)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSerializeQRIS(b *testing.B) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Serialize()
	}
}

func BenchmarkSetAmount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q, _ := NewQRISFromString(sampleQRISRaw)
		q.SetAmountWithOptions(QRISAmountOptions{
			Amount:   10000,
			FeeType:  QRISFeePercent,
			FeeValue: 2.5,
		})
	}
}

func BenchmarkGenerateQRISImage(b *testing.B) {
	q, _ := NewQRISFromString(sampleQRISRaw)
	q.SetAmountWithOptions(QRISAmountOptions{
		Amount:   10000,
		FeeType:  QRISFeeRupiah,
		FeeValue: 2000,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := q.GenerateQRISImage(256)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func contains(s, substr string) bool {
	return stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
