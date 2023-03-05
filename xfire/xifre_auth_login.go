package xfire

import (
	"bufio"
	"fmt"
	"net"
)

func GetLoginRequestMessage(stream *bufio.Reader) bool {

	msgLen, msgID, msgAttr := getMessageEssentials(stream)

	fmt.Printf("GetLoginRequestMessage: read message with ID = %d, length: %d bytes with %d attributes\n", msgID, msgLen, msgAttr)

	username, _ := getStringAttribute(stream, "name")
	fmt.Printf("%s - username\n", username)

	password, _ := getStringAttribute(stream, "password")
	fmt.Printf("%s - password (encrypted)\n", password)

	userFlags, _ := getInt32Attribute(stream, "flags")
	fmt.Printf("%x - flags\n", userFlags)

	sid, _ := getSIDAttribute(stream)
	fmt.Printf("%x - sid\n", sid)

	return true
}

func SendLoginSuccess(connection net.Conn) {
	const MSG_ID = 130

	header := makeMessageHeader(MSG_ID, 16)

	userid := makeInt32Attribute("userid", generateUserID("testingphase"))
	sid := makeSIDAttribute(generateSID("testingphase"))
	nick := makeStringAttribute("nick", "Douche Baggins")
	status := makeInt32Attribute("status", 0)

	dlSet := makeStringAttribute("dlset", "")
	p2pset := makeStringAttribute("p2pset", "")
	clntSet := makeStringAttribute("clntset", "")

	minRect := makeInt32Attribute("minrect", 1)
	maxRect := makeInt32Attribute("maxrect", 164867)

	ctry := makeInt32Attribute("ctry", 144)
	n1 := makeIPv4Attribute("n1", "204.71.190.131")
	n2 := makeIPv4Attribute("n2", "204.71.190.132")
	n3 := makeIPv4Attribute("n3", "204.71.190.133")
	pip := makeIPv4Attribute("pip", "127.0.0.1")

	salt := makeStringAttribute("salt", "babaroga")
	reason := makeStringAttribute("reason", "Mq_P8Ad3aMEUvFinw0ceu6FITnZTWXxg46XU8xHW")

	message := formMessage(header, userid, sid, nick, status,
		dlSet, p2pset, clntSet, minRect, maxRect, ctry, n1, n2, n3, pip, salt, reason)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func SendLoginFailure(connection net.Conn) {
	const MSG_ID = 129

	header := makeMessageHeader(MSG_ID, 1)
	reason := makeInt32Attribute("reason", 0)

	message := formMessage(header, reason)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func SendLoginChallengeMessage(connection net.Conn, saltValue string) {
	const MSG_ID = 0x0080

	header := makeMessageHeader(MSG_ID, 1)
	salt := makeStringAttribute("salt", saltValue)

	message := formMessage(header, salt)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func SendFriendsListMessage(connection net.Conn) {
	const MSG_ID = 131

	header := makeMessageHeader(MSG_ID, 3)
	userid := makeInt32ListAttribute("userid", []uint32{})
	friends := makeStringListAttribute("friends", []string{})
	nick := makeStringListAttribute("nick", []string{})

	message := formMessage(header, userid, friends, nick)

	message = setMessageLength(message)

	connection.Write(message)
}

func SendSessionIDListMessage(connection net.Conn) {
	const MSG_ID = 132

	header := makeMessageHeader(MSG_ID, 2)
	userid := makeInt32ListAttribute("userid", []uint32{})
	sid := makeSIDListAttribute("sid", []sid{})

	message := formMessage(header, userid, sid)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)

}

func SendGroupsList(connection net.Conn) {
	const MSG_ID = 151

	header := makeMessageHeader(MSG_ID, 2)
	groupids := makeInt32ListAttribute(string([]byte{0x19}), []uint32{})
	groupnames := makeStringListAttribute(string([]byte{0x1a}), []string{})

	message := formMessage(header, groupids, groupnames)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func SendGroupsFriendsList(connection net.Conn) {
	const MSG_ID = 152

	header := makeMessageHeader(MSG_ID, 2)
	userIDs := makeInt32ListAttribute(string([]byte{0x01}), []uint32{})
	groupIDs := makeInt32ListAttribute(string([]byte{0x19}), []uint32{})

	message := formMessage(header, userIDs, groupIDs)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func SendClientPreferences(connection net.Conn) {
	const MSG_ID = 141

	header := makeMessageHeader(MSG_ID, 1)

	dict := makeEmptyIntDict(string([]byte{byte(0x4c)}), 12)

	dict = append(dict, makeIntDictStringEntry(1, "0")...)
	dict = append(dict, makeIntDictStringEntry(4, "0")...)
	dict = append(dict, makeIntDictStringEntry(5, "0")...)
	dict = append(dict, makeIntDictStringEntry(6, "1")...)

	dict = append(dict, makeIntDictStringEntry(7, "0")...)
	dict = append(dict, makeIntDictStringEntry(8, "0")...)
	dict = append(dict, makeIntDictStringEntry(11, "0")...)
	dict = append(dict, makeIntDictStringEntry(17, "0")...)

	dict = append(dict, makeIntDictStringEntry(18, "0")...)
	dict = append(dict, makeIntDictStringEntry(19, "0")...)
	dict = append(dict, makeIntDictStringEntry(20, "0")...)
	dict = append(dict, makeIntDictStringEntry(21, "0")...)

	message := formMessage(header, dict)

	message = setMessageLength(message)

	debugPrint(message)

	connection.Write(message)
}

func GetClientConfigurationMessage(stream *bufio.Reader) {
	const MSG_ID = 16

	msgLen, msgID, msgAttr := getMessageEssentials(stream)

	fmt.Printf("GetClientConfigurationMessage: read message with ID = %d, length: %d bytes with %d attributes\n", msgID, msgLen, msgAttr)

	lang, _ := getStringAttribute(stream, "lang")
	fmt.Printf("%s - lang\n", lang)

	skin, _ := getStringAttribute(stream, "skin")
	fmt.Printf("%s - skin\n", skin)

	theme, _ := getStringAttribute(stream, "theme")
	fmt.Printf("%x - theme\n", theme)

	partner, _ := getStringAttribute(stream, "partner")
	fmt.Printf("%x - partner\n", partner)

}
