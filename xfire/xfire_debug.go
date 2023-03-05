package xfire

import "fmt"

func debugPrint(message []byte) {
	fmt.Printf("message: %x | LEN: %d | calcLEN : %d\n", message, message[0], len(message))
	fmt.Printf("message: %s\n", message)
}
