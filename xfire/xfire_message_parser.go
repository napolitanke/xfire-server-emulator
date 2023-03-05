package xfire

import (
	"bufio"
	"io"
)

func parseMessage(stream *bufio.Reader, messageLength uint16) {
	message := make([]byte, messageLength)
	io.ReadFull(stream, message)

	switch messageID := LittleEndianToUInt16([2]byte{message[0], message[1]}); messageID {
	case 0x01:
		getLoginRequest(message[2:])
	}

}

func fetchMessage(stream *bufio.Reader, firstByte byte) {
	var singleByteBuffer [1]byte
	io.ReadFull(stream, singleByteBuffer[:])

	messageLength := LittleEndianToUInt16([2]byte{firstByte, singleByteBuffer[0]}) - 2

	parseMessage(stream, messageLength)

}

func GetMessageFromClient(stream *bufio.Reader) {

	var singleByteBuffer [1]byte

	for messageDetected := false; messageDetected; {
		io.ReadFull(stream, singleByteBuffer[:])
		if singleByteBuffer[0] != 0x0 {
			messageDetected = true
			fetchMessage(stream, singleByteBuffer[0])
		}

	}

}
