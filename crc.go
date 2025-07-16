package qris

import (
	"fmt"
	"strings"
)

func UpdateCRC(data []TLV) []TLV {
	filtered := []TLV{}
	for _, tlv := range data {
		if tlv.Tag != string(TagCRC) {
			filtered = append(filtered, tlv)
		}
	}
	partial := SerializeTLV(filtered) + string(TagCRC) + "04"
	crc := CalculateCRC(partial)
	filtered = append(filtered, TLV{Tag: string(TagCRC), Len: 4, Value: crc})
	return filtered
}

func CalculateCRC(payload string) string {
	crc := uint16(0xFFFF)
	for _, c := range payload {
		crc ^= uint16(c) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
		crc &= 0xFFFF
	}
	return strings.ToUpper(fmt.Sprintf("%04X", crc))
}
