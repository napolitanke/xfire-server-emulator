package main

import (
	"bufio"
	"fmt"
	"net"

	"./xfire"
)

func invertBuffer(buffer []byte) []byte {
	return buffer[:len(buffer)-1]
}

func compareBuffer(buffer []byte, item []byte) bool {

	if !(len(buffer) == len(item)) {
		return false
	}
	for i := 0; i < len(buffer); i++ {
		if buffer[i] != item[i] {
			return false
		}
	}
	return true
}

func handleConnection(connection net.Conn) {

	UA01 := []byte{0x55, 0x41, 0x30, 0x31}

	stream := bufio.NewReader(connection)
	user := connection.RemoteAddr().String()

	if !compareBuffer(xfire.StartHandshake(stream), UA01) {
		fmt.Printf("Error 1 - wrong handshake!\n")
		return
	}
	fmt.Printf("Client (%s) connected!\n", user)

	if xfire.GetClientInformationMessage(stream) {
		fmt.Printf("Client (%s) is of correct version. Proceeding...\n", user)
		//connection.Close()
		xfire.SendLoginChallengeMessage(connection, "babaroga")
		//xfire.SendLoginChallengeMessageOLD(connection)
		if xfire.GetLoginRequestMessage(stream) {
			xfire.SendLoginSuccess(connection)
			//time.Sleep(1 * time.Second)
			xfire.GetClientConfigurationMessage(stream)
			xfire.SendClientPreferences(connection)
			xfire.SendGroupsList(connection)
			xfire.SendGroupsFriendsList(connection)
			xfire.SendFriendsListMessage(connection)
			xfire.SendSessionIDListMessage(connection)
		} else {
			xfire.SendLoginFailure(connection)
		}

	}

	//handleConnection(connection)
}

func main() {
	const XFIRE_PORT = ":25999"

	server, _ := net.Listen("tcp", XFIRE_PORT)

	for {
		connection, _ := server.Accept()
		go handleConnection(connection)
	}

}
