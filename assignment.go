package assignment

import (
	"errors"
	"fmt"
)

// DataElement define an interface representing the allowed types.
type DataElement interface {
	isDataElement()
}

func (Data) isDataElement()       {}
func (DataInt32) isDataElement()  {}
func (DataString) isDataElement() {}

type (
	Data       []DataElement
	DataInt32  int32
	DataString string
	// Type identifier
	Type int8
)

const (
	IntegerType Type = 0
	StringType  Type = 1
	ArrayType   Type = 2
)

func Encode(toSend Data) (string, error) {
	var encoded []byte
	for _, el := range toSend {
		encodedElement, err := encodeElement(el)
		if err != nil {
			return "", err
		}
		encoded = append(encoded, encodedElement...)
	}
	return string(encoded), nil
}

// Helper function to recursively encode Data type
func encodeElement(element DataElement) ([]byte, error) {
	var encoded []byte
	switch el := element.(type) {
	case DataInt32:
		encoded = append(encoded, []byte{byte(IntegerType)}...)
		encoded = append(encoded, encodeInt(uint32(el))...)
		return encoded, nil
	case DataString:
		encodedStr := []byte(el)
		length := encodeInt(uint32(len(encodedStr)))
		bytes := []byte{byte(StringType)}
		bytes = append(bytes, length...)
		bytes = append(bytes, encodedStr...)
		encoded = append(encoded, bytes...)
		return encoded, nil
	case Data:
		var encodedElements []byte
		for _, element := range el {
			encodedElement, err := encodeElement(element)
			if err != nil {
				return nil, err
			}
			encodedElements = append(encodedElements, encodedElement...)
		}
		length := encodeInt(uint32(len(el)))
		encoded = append(encoded, []byte{byte(ArrayType)}...)
		encoded = append(encoded, length...)
		encoded = append(encoded, encodedElements...)
		return encoded, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown type: %T", el))
	}
}

func encodeInt(value uint32) []byte {
	var encoded []byte
	return append(encoded,
		byte(value>>0),
		byte(value>>8),
		byte(value>>16),
		byte(value>>24))
}

func decodeInt(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, errors.New("wrong int size")
	}
	return uint32(b[0])<<0 | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

func Decode(received string) (Data, error) {
	dataBytes := []byte(received)

	var elements Data
	index := 0
	for index < len(received) {
		element, newIndex, err := decodeElement(dataBytes, index)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
		index = newIndex
	}

	return elements, nil
}

// Helper function to recursively decode Data type
func decodeElement(data []byte, index int) (DataElement, int, error) {
	if index >= len(data) {
		return nil, index, errors.New("unexpected end of data")
	}

	typeByte := data[index]
	index++

	switch Type(typeByte) {
	case IntegerType:
		if index+4 > len(data) {
			return nil, index, errors.New("invalid integer data")
		}

		intValue, err := decodeInt(data[index : index+4])
		if err != nil {
			return nil, index, err
		}

		index += 4
		return DataInt32(intValue), index, nil
	case StringType:
		if index+4 > len(data) {
			return nil, index, errors.New("invalid string length")
		}

		length, err := decodeInt(data[index : index+4])
		if err != nil {
			return nil, index, err
		}

		index += 4
		if index+int(length) > len(data) {
			return nil, index, errors.New("invalid string data")
		}

		encodedStr := data[index : index+int(length)]
		index += int(length)
		strValue := utf8Decode(encodedStr)
		return DataString(strValue), index, nil
	case ArrayType:
		if index+4 > len(data) {
			return nil, index, errors.New("invalid array length")
		}

		length, err := decodeInt(data[index : index+4])
		if err != nil {
			return nil, index, err
		}

		index += 4
		var array Data
		for i := uint32(0); i < length; i++ {
			el, newIndex, err := decodeElement(data, index)
			if err != nil {
				return nil, index, err
			}
			array = append(array, el)
			index = newIndex
		}
		return array, index, nil
	default:
		return nil, index, errors.New("invalid type indicator")
	}
}

func utf8Decode(b []byte) string {
	// Assume valid UTF-8
	return string(b)
}
