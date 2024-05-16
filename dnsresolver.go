package dnsresolver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const TYPE_A = 1
const CLASS_IN = 1

type DNSHeader struct {
	ID                  uint16
	Flags               uint16
	NumberOfQuestions   uint16
	NumberOfAuthorities uint16
	NumberOfAdditionals uint16
	NumberOfAnswers     uint16
}

const DNSHeaderSize = 12

type DNSQuestion struct {
	Name  string
	Type_ uint16
	Class uint16
}

func (h *DNSHeader) HeaderToBytes(size int) string {
	headerSlice := make([]byte, DNSHeaderSize)
	var formatted []byte
	// is it better to write v >> 8 ??
	// BigEndian saves the most significant piece of data at the lowest in memory
	// Endian is useful when data isn't single-byte. In our case, each element is multiple bytes
	binary.BigEndian.PutUint16(headerSlice[0:2], h.ID)
	binary.BigEndian.PutUint16(headerSlice[2:4], h.Flags)
	binary.BigEndian.PutUint16(headerSlice[4:6], h.NumberOfQuestions)
	binary.BigEndian.PutUint16(headerSlice[6:8], h.NumberOfAuthorities)
	binary.BigEndian.PutUint16(headerSlice[8:10], h.NumberOfAdditionals)
	binary.BigEndian.PutUint16(headerSlice[10:12], h.NumberOfAnswers)

	for _, b := range headerSlice {
		formatted = append(formatted, fmt.Sprintf("\\x%02x", b)...)
	}
	return string(formatted)
}

// What should be the format of QuestionToBytes?
func (q *DNSQuestion) QuestionToBytes(size int) []byte {

	questionSlice := []byte(q.Name)

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, q.Type_)
	binary.Write(buffer, binary.BigEndian, q.Class)

	return append(questionSlice, buffer.Bytes()...)
}

func EncodeDnsName(DomainName string) []byte {
	var encodedDomain []byte
	parts := bytes.Split([]byte(DomainName), []byte("."))
	for _, part := range parts {
		encodedDomain = append(encodedDomain, byte(len(part)))
	}
	encodedDomain = append(encodedDomain, 0x00)
	return encodedDomain
}

// max for this function: 65535
func BuildDNSQuery(DomaineName, RecordType string) string {
	name := EncodeDnsName(DomaineName)
	var id uint16

	for n := 0; n < 65535; n++ {
		id = uint16(rand.Intn(65535))
	}
	recursionDesired := 1 << 8
	header := DNSHeader{
		ID:                id,
		NumberOfQuestions: 1,
		Flags:             uint16(recursionDesired),
	}
	question := DNSQuestion{
		Name:  string(name),
		Type_: TYPE_A,
		Class: CLASS_IN,
	}
	return header.HeaderToBytes(DNSHeaderSize) + string(question.QuestionToBytes(12))

}

func Main() int {
	h := DNSHeader{
		ID:                  0x1314,
		Flags:               0,
		NumberOfQuestions:   1,
		NumberOfAnswers:     0,
		NumberOfAuthorities: 0,
		NumberOfAdditionals: 0,
	}
	fmt.Println(h.HeaderToBytes(DNSHeaderSize))
	// Example usage
	q := DNSQuestion{
		Name:  "example.com",
		Type_: 1, // A record
		Class: 1, // IN class
	}
	// chose 5 because it is the minimal length of a DNS resolver Question
	fmt.Println("question to bytes", q.QuestionToBytes(5))

	fmt.Println("Encoded name", EncodeDnsName("example.com"))
	fmt.Println("Building the query", BuildDNSQuery("example.com", "1"))

	return 0
}
