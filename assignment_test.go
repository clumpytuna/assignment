package assignment

import (
	"testing"
)

// Test cases for encoding and decoding
func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name string
		data Data
	}{
		{
			name: "Simple Integer",
			data: Data{DataInt32(42)},
		},
		{
			name: "Negative Integer",
			data: Data{DataInt32(-123456)},
		},
		{
			name: "Simple String",
			data: Data{DataString("hello")},
		},
		{
			name: "UTF-8 String",
			data: Data{DataString("こんにちは")}, // "Hello" in Japanese
		},
		{
			name: "Empty String",
			data: Data{DataString("")},
		},
		{
			name: "Simple Array",
			data: Data{
				DataInt32(1),
				DataInt32(2),
				DataInt32(3),
			},
		},
		{
			name: "Nested Arrays",
			data: Data{
				DataInt32(1),
				Data{
					DataInt32(2),
					Data{
						DataInt32(3),
					},
				},
			},
		},
		{
			name: "Mixed Types",
			data: Data{
				DataString("foo"),
				DataInt32(42),
				Data{
					DataString("bar"),
					DataInt32(-42),
					Data{
						DataString("baz"),
					},
				},
			},
		},
		{
			name: "Empty Array",
			data: Data{},
		},
		{
			name: "Max Int32",
			data: Data{DataInt32(2147483647)},
		},
		{
			name: "Min Int32",
			data: Data{DataInt32(-2147483648)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := Encode(tc.data)
			if err != nil {
				t.Fatalf("Encoding failed: %v", err)
			}

			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decoding failed: %v", err)
			}
			if !compareDataElements(tc.data, decoded) {
				t.Errorf("Decoded data does not match original.\nOriginal: %v\nDecoded: %v", tc.data, decoded)
			}
		})
	}
}

// Additional test for error handling
func TestDecodeInvalidData(t *testing.T) {
	invalidData := []byte{0xFF} // Invalid type indicator
	_, err := Decode(string(invalidData))
	if err == nil {
		t.Error("Expected error when decoding invalid data, but got none")
	}
}

// compareDataElements function
func compareDataElements(a, b DataElement) bool {
	switch va := a.(type) {
	case DataInt32:
		vb, ok := b.(DataInt32)
		return ok && va == vb
	case DataString:
		vb, ok := b.(DataString)
		return ok && va == vb
	case Data:
		vb, ok := b.(Data)
		if !ok || len(va) != len(vb) {
			return false
		}
		for i := range va {
			if !compareDataElements(va[i], vb[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
