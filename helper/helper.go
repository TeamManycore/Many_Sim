package helper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// LoadFile does essentially what it says it does. It loads a file and transforms it into a
// two dimensional array of blocks of four bytes, without whitespaces. This is needed for further
// processing
func LoadFile(file string) [][]byte {
	content, err := ioutil.ReadFile(string(file))
	if err != nil {
		panic(err)
	}

	content = []byte(strings.Join(strings.Fields(string(content)), ""))
	newContent := bytes.ToLower(content)
	var divided [][]byte

	chunkSize := 4

	for i := 0; i < len(newContent); i += chunkSize {
		end := i + chunkSize

		if end > len(newContent) {
			end = len(newContent)
		}

		divided = append(divided, newContent[i:end])
	}
	return divided
}

// Flag creates a word from a boolean. This is used for certain ALU operations, as they would usually output a bool.
// Since the CPU uses only uint16 the boolean value has to be tranfromed appropriately
func Flag(yesno bool) uint16 {
	if yesno {
		return 0x0001
	} else {
		return 0
	}
}

// Signed transforms a uint16 into a signed integer. This is needed for the same reason as the function Flag.
// When a signed integer is needed, this will convert a unsigned integer
func Signed(data uint16) int16 {
	if (data & 0x8000) == 0 {
		return int16(data)
	} else {
		return int16(0 - ((^(data) + 1) & 0xFFFF))
	}
}

// WordMalformed is needed, because certain parts of go are just bollocks.
// Go should definitly not be used in embedded systems
func WordMalformed(b uint16) uint16 {
	var val uint16 = b
	if (b & 0xFFF0) == 0 {
		val = b << 12
	} else if (b & 0xFF00) == 0 {
		val = b << 8
	} else if (b & 0xF000) == 0 {
		val = b << 4
	}
	return val
}

// Word2hex takes a word and transforms it into a four character string, that will be needed for creating an image
func Word2hex(n uint16) string {
	s := fmt.Sprintf("%04x", n)
	return s
}

// Hex2word takes a bytearray, which consists of ASCII values, and transforms it into a word that can be processed by the cpu
func Hex2word(data []byte) uint16 {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return uint16(decoded[0])<<8 |
		uint16(decoded[1])
}
