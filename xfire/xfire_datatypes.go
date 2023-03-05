package xfire

const (
	STR        = 0x01
	INT32      = 0x02
	SESSION_ID = 0x03
	LIST       = 0x04
	DICT       = 0x05
	DID        = 0x06
	ARRAY      = 0x09
)

const DID_LENGTH = 21

type sid [16]byte
