package xfire

import (
	"bufio"
	"fmt"
	"io"
)

func StartHandshake(stream *bufio.Reader) []byte {

	var message []byte
	var buffer [4]byte
	io.ReadFull(stream, buffer[:])
	message = buffer[:]
	fmt.Printf("%s\n", message)
	return message
}

func GetClientInformationMessage(stream *bufio.Reader) bool {
	var twoByteBuffer [2]byte

	io.ReadFull(stream, twoByteBuffer[:])
	messageLength := LittleEndianToUInt16(twoByteBuffer)

	io.ReadFull(stream, twoByteBuffer[:])
	messageID := LittleEndianToUInt16(twoByteBuffer)

	fmt.Printf("Lenght: %d (%x) --- ID: %x", messageLength, messageLength, messageID)

	messageBuffer := make([]byte, messageLength-4)
	io.ReadFull(stream, messageBuffer)

	offset := 0

	var dataString string
	var version uint32
	for offset < int(messageLength)-4 {
		switch valueType := messageBuffer[offset]; valueType {
		case STR:
			offset++
			dataLength := uint8(messageBuffer[offset])
			fmt.Printf("datalength %d\n", dataLength)
			offset++
			dataString = string(messageBuffer[offset : offset+int(dataLength)])
			offset += int(dataLength)
		case INT32:
			offset++
			version = LittleEndianToUInt32Slice(messageBuffer[offset:(offset + 4)])
			offset += 4
		}

	}
	fmt.Printf("%s : %d\n", dataString, version)
	if version == 133 {
		return true
	} else {
		return false
	}
}
