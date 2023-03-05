package xfire

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func getMessageEssentials(stream *bufio.Reader) (messageLength uint16, messageID uint16, numberOfAttributes uint8) {

	var twoByteBuffer [2]byte
	var singleByteBuffer [1]byte

	io.ReadFull(stream, twoByteBuffer[:])

	messageLength = LittleEndianToUInt16(twoByteBuffer)
	fmt.Printf("getMessageEssentials: messageLength = %d (0x%x)\n", messageLength, twoByteBuffer)

	io.ReadFull(stream, twoByteBuffer[:])

	messageID = LittleEndianToUInt16(twoByteBuffer)
	fmt.Printf("getMessageEssentials: messageID = %d (0x%x)\n", messageID, twoByteBuffer)

	io.ReadFull(stream, singleByteBuffer[:])

	numberOfAttributes = uint8(singleByteBuffer[0])
	fmt.Printf("getMessageEssentials: numberOfAttributes = %d (0x%x)\n", numberOfAttributes, singleByteBuffer)

	return messageLength, messageID, numberOfAttributes
}

func getMessageAll(stream *bufio.Reader) []byte {
	var twoByteBuffer [2]byte
	io.ReadFull(stream, twoByteBuffer[:])
	messageLength := LittleEndianToUInt16(twoByteBuffer)

	message := make([]byte, messageLength-2)
	io.ReadFull(stream, message)

	var output []byte
	output = append(output, twoByteBuffer[:]...)
	output = append(output, message...)

	return output
}

func getStringAttribute(stream *bufio.Reader, attributeName string) (string, error) {

	var attrLen [1]byte
	io.ReadFull(stream, attrLen[:])

	if int(attrLen[0]) == len(attributeName) {
		attrName := make([]byte, len(attributeName))
		io.ReadFull(stream, attrName[:])
		if string(attrName[:]) == attributeName {
			var valType [1]byte
			io.ReadFull(stream, valType[:])
			if valType[0] == STR {
				var valLen [2]byte
				io.ReadFull(stream, valLen[:])
				valueLength := LittleEndianToUInt16(valLen)

				rawName := make([]byte, valueLength)
				io.ReadFull(stream, rawName[:])

				return string(rawName[:]), nil
			}

			return "", errors.New("getStringAttribute: Wrong value type!")
		}

		return "", errors.New("getStringAttribute: Wrong argument in buffer!")

	}

	return "", errors.New("getStringAttribute: Wrong argument lenght!")
}

func getInt32Attribute(stream *bufio.Reader, attributeName string) (uint32, error) {

	var attrLen [1]byte
	io.ReadFull(stream, attrLen[:])

	if int(attrLen[0]) == len(attributeName) {
		attrName := make([]byte, len(attributeName))
		io.ReadFull(stream, attrName[:])
		if string(attrName[:]) == attributeName {
			var valType [1]byte
			io.ReadFull(stream, valType[:])
			if valType[0] == INT32 {
				var val [4]byte
				io.ReadFull(stream, val[:])
				value := LittleEndianToUInt32(val)

				return value, nil
			}

			return 0, errors.New("getInt32Attribute: Wrong value type!")
		}

		return 0, errors.New("getInt32Attribute: Wrong argument in buffer!")

	}

	return 0, errors.New("getInt32Attribute: Wrong argument lenght!")
}

func getSIDAttribute(stream *bufio.Reader) ([]byte, error) {

	var attrLen [1]byte
	io.ReadFull(stream, attrLen[:])
	attributeName := "sid"

	if int(attrLen[0]) == len(attributeName) {
		attrName := make([]byte, len(attributeName))
		io.ReadFull(stream, attrName[:])
		if string(attrName[:]) == attributeName {
			var valType [1]byte
			io.ReadFull(stream, valType[:])
			if valType[0] == SESSION_ID {
				var val [16]byte
				io.ReadFull(stream, val[:])

				return val[:], nil
			}

			return nil, errors.New("getSIDAttribute: Wrong value type!")
		}

		return nil, errors.New("getSIDAttribute: Wrong argument in buffer!")

	}

	return nil, errors.New("getSIDAttribute: Wrong argument lenght!")
}

func makeMessageHeader(messageID uint16, args uint8) []byte {

	var output []byte

	// 01..LENGTH
	output = append(output, 0x00)
	output = append(output, 0x00)

	//02..ID
	r := UInt16ToLittleEndian(uint16(messageID))
	for i := 0; i < len(r); i++ {
		output = append(output, r[i])
	}

	//03..NUMBER OF ATTRIBUTES
	output = append(output, args)

	return output
}

func setMessageLength(message []byte) []byte {
	length := uint16(len(message))

	r := UInt16ToLittleEndian(length)
	for i := 0; i < len(r); i++ {
		message[i] = r[i]
	}

	return message
}

func makeStringAttribute(attrName string, value string) []byte {

	var output []byte

	output = append(output, byte(len(attrName)))

	name := []byte(attrName)
	output = append(output, name...)

	output = append(output, STR)
	valLen := UInt16ToLittleEndian(uint16(len(value)))
	output = append(output, valLen...)

	val := []byte(value)
	output = append(output, val...)

	debugPrint(output)

	return output
}

func makeInt32Attribute(attrName string, value uint32) []byte {

	var output []byte

	output = append(output, byte(len(attrName)))

	name := []byte(attrName)
	output = append(output, name...)

	output = append(output, INT32)
	output = append(output, UInt32ToLittleEndian(value)...)

	return output
}

func makeIPv4Attribute(attrName string, value string) []byte {

	var output []byte

	output = append(output, byte(len(attrName)))

	name := []byte(attrName)
	output = append(output, name...)

	output = append(output, INT32)
	addressStr := strings.Split(value, ".")

	for _, element := range addressStr {
		x, _ := strconv.Atoi(element)
		output = append(output, byte(x))
	}

	return output
}

func makeSIDAttribute(value []byte) []byte {

	attrName := "sid"
	attrNameLen := []byte{byte(len(attrName))}

	attrNameArray := []byte(attrName)

	attrType := []byte{byte(SESSION_ID)}

	output := formMessage(attrNameLen, attrNameArray, attrType, value)

	debugPrint(output)
	return output
}

func formMessage(parts ...[]byte) []byte {
	var output []byte

	for _, part := range parts {
		output = append(output, part...)
	}

	return output
}

func makeInt32ListAttribute(name string, attributes []uint32) []byte {

	nameLen := []byte{byte(len(name))}
	nameBytes := []byte(name)

	msgType := []byte{byte(LIST), byte(INT32)}

	numAttributes := []byte(UInt16ToLittleEndian(uint16(len(attributes))))

	message := formMessage(nameLen, nameBytes, msgType, numAttributes)

	for _, attribute := range attributes {
		message = append(message, UInt32ToLittleEndian(attribute)...)
	}

	return message
}

func makeSIDListAttribute(name string, attributes []sid) []byte {

	nameLen := []byte{byte(len(name))}
	nameBytes := []byte(name)

	msgType := []byte{byte(LIST), byte(SESSION_ID)}

	numAttributes := []byte(UInt16ToLittleEndian(uint16(len(attributes))))

	message := formMessage(nameLen, nameBytes, msgType, numAttributes)

	for _, attribute := range attributes {
		message = append(message, attribute[:]...)
	}

	return message
}

func makeStringListAttribute(name string, attributes []string) []byte {

	nameLen := []byte{byte(len(name))}
	nameBytes := []byte(name)

	msgType := []byte{byte(LIST), byte(STR)}

	numAttributes := []byte(UInt16ToLittleEndian(uint16(len(attributes))))

	message := formMessage(nameLen, nameBytes, msgType, numAttributes)

	for _, attribute := range attributes {
		attrLen := UInt16ToLittleEndian(uint16(len(attribute)))
		message = append(message, attrLen...)
		message = append(message, []byte(attribute)...)
	}

	return message
}

func makeEmptyIntDict(name string, numberOfPairs uint8) []byte {

	nameLength := []byte{byte(len(name))}

	numberOfEntries := []byte{byte(numberOfPairs)}

	message := formMessage(nameLength, []byte(name), []byte{byte(ARRAY)}, numberOfEntries)

	return message
}

func makeIntDictStringEntry(key uint8, value string) []byte {

	entryKey := []byte{(byte(key))}
	entryType := []byte{byte(STR)}
	entryValueLength := UInt16ToLittleEndian(uint16(len(value)))
	entryValue := []byte(value)

	message := formMessage(entryKey, entryType, entryValueLength, entryValue)

	return message
}
