package qris

import (
	"fmt"
	"strconv"
	"strings"
)

type TLV struct {
	Tag   string
	Len   int
	Value string
	Sub   []TLV
}

func ParseTLV(s string) ([]TLV, error) {
	var result []TLV
	for i := 0; i < len(s); {
		if i+4 > len(s) {
			break
		}
		tag := s[i : i+2]
		length, err := strconv.Atoi(s[i+2 : i+4])
		if err != nil || i+4+length > len(s) {
			return nil, fmt.Errorf("invalid TLV segment at %d", i)
		}
		value := s[i+4 : i+4+length]
		sub := []TLV{}
		if tag >= "26" && tag <= "51" {
			sub, _ = ParseTLV(value)
		}
		result = append(result, TLV{Tag: tag, Len: length, Value: value, Sub: sub})
		i += 4 + length
	}
	return result, nil
}

func SerializeTLV(data []TLV) string {
	var builder strings.Builder
	for _, tlv := range data {
		lengthStr := fmt.Sprintf("%02d", len(tlv.Value))
		builder.WriteString(tlv.Tag)
		builder.WriteString(lengthStr)
		builder.WriteString(tlv.Value)
	}
	return builder.String()
}

func ReplaceTLVValue(data []TLV, targetTag string, newValue string) []TLV {
	for i := range data {
		if data[i].Tag == targetTag {
			data[i].Value = newValue
			data[i].Len = len(newValue)
			data[i].Sub = nil
			break
		}
	}
	return UpdateCRC(data)
}

func RemoveTLV(data []TLV, tag string) []TLV {
	result := []TLV{}
	for _, tlv := range data {
		if tlv.Tag != tag {
			result = append(result, tlv)
		}
	}
	return result
}

func RemoveTLVPrefix(data []TLV, prefix string) []TLV {
	result := []TLV{}
	for _, tlv := range data {
		if !strings.HasPrefix(tlv.Tag, prefix) {
			result = append(result, tlv)
		}
	}
	return result
}
