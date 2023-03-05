package xfire

import (
	"fmt"
	"math/rand"
	"time"
)

func LittleEndianToUInt16(buf [2]byte) uint16 {
	output := uint16(buf[0]) + (uint16(buf[1]) << 8)
	return output
}

func UInt16ToLittleEndian(number uint16) []byte {
	var output []byte

	output = append(output, byte(number%256))
	output = append(output, byte(number/256))
	return output
}

func LittleEndianToUInt32(buf [4]byte) uint32 {
	output := uint32(buf[0]) + (uint32(buf[1]) << 8) + (uint32(buf[2]) << 16) + (uint32(buf[3]) << 24)
	return output
}

func LittleEndianToUInt32Slice(buf []byte) uint32 {
	output := uint32(buf[0]) + (uint32(buf[1]) << 8) + (uint32(buf[2]) << 16) + (uint32(buf[3]) << 24)
	return output
}

func UInt32ToLittleEndian(number uint32) []byte {
	var output []byte

	for i := 0; i < 4; i++ {
		output = append(output, byte(number%256))
		number = number / 256
	}

	return output
}

func getHash(value string) uint32 {

	sum := uint32(0)

	for _, element := range value {
		sum += uint32(element)
	}

	return sum
}

func generateUserID(username string) uint32 {

	base := uint32(time.Now().Unix())
	mask := uint32(base & 0x000000FF)
	hash := getHash(username) * mask

	return ((base + hash) & 0xFFFFFF00)
}

func generateSID(username string) []byte {
	base := UInt32ToLittleEndian(uint32(time.Now().Unix()))
	rand := UInt32ToLittleEndian(uint32(rand.Int31()))

	user := []byte(username)

	fmt.Printf("LEN USER: %d", len(user))

	if len(user) < 9 {
		for i := 0; i < 8-len(user); i++ {
			user = append(user, 0)
		}
	} else {
		user = user[0:8]
	}

	sid := formMessage(base, rand, user)

	fmt.Printf("base: %x\nrand: %x\nuser: %x\n", base, rand, user)
	fmt.Printf("SID: %x -- %d\n", sid, len(sid))
	return sid
}
