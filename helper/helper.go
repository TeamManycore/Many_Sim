package helper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

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

func N(content uint16) bool {
	return (content & 0x8000) != 0
}

func Flag(yesno bool) uint16 {
	if yesno {
		return 0xFFFF
	} else {
		return 0
	}
}

func Signed(data uint16) int64 {
	if (data & 0x8000) == 0 {
		return int64(data)
	} else {
		return int64(0 - ((^(data) + 1) & 0xFFFF))
	}
}

func ByteMalformed(b uint16) uint16 {
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

func Word2hex(n uint16) string {
	s := fmt.Sprintf("%04x", n)
	return s
}

func Hex2word(data []byte) uint16 {
	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return uint16(decoded[0])<<8 |
		uint16(decoded[1])
}
