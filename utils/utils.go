package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"regexp"
	"strings"
)

// Cookie is used to give a unique identifier to each request.
type Cookie uint64

// CreateCookie Create random 8 byte cookie
func CreateCookie() Cookie {
	var buf [8]byte
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return Cookie(binary.LittleEndian.Uint64(buf[:]))
}

// MakeResource Creae random 16 byte string resource
func MakeResource() string {
	var buf [16]byte
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return fmt.Sprintf("%x", buf)
}

// GenerateUUIDUsingRandomNumber generate UUID
func GenerateUUIDUsingRandomNumber() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

// GenerateUUID generate UUID
func GenerateUUID() (uuidStr string) {
	//uuid, err := uuid.NewV1()
	uuid := uuid.NewV4()
	return uuid.String()
}

// GenerateAlias generate user alias
func GenerateAlias() (uuidStr string) {
	return GenerateUUID()
}

// GetHostnameFromNetAddress Return host name from a provided net address
func GetHostnameFromNetAddress(netAddress net.Addr) string {
	return strings.Split(netAddress.String(), ":")[0]
}

// StripOutComments strip out comments in a string
func StripOutComments(content string) string {
	re := regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")
	newBytes := re.ReplaceAll([]byte(content), nil)
	return string(newBytes)
}

// TruncateString truncate string to max length
func TruncateString(content string, maxLength int) string {
	contentLen := len(content)
	if contentLen > maxLength {
		content = fmt.Sprintf("%s...", content[0:maxLength])
	}

	return content
}
