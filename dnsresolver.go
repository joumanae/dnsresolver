package dnsresolver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
)

const (
	TYPE_A   = 1
	CLASS_IN = 1
)

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

type DNSRecord struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32
	Data  string
}

func (h *DNSHeader) HeaderToBytesToHexadecimal(size int) string {
	headerSlice := make([]byte, DNSHeaderSize)
	var formatted []byte
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

func (q *DNSQuestion) ToBytes() []byte {
	questionSlice := []byte(q.Name)

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, q.Type_)
	binary.Write(buffer, binary.BigEndian, q.Class)

	return append(questionSlice, buffer.Bytes()...)
}

func EncodeDnsName(domainName string) string {
	var encodedDomain []byte
	parts := bytes.Split([]byte(domainName), []byte("."))
	for _, part := range parts {
		encodedDomain = append(encodedDomain, byte(len(part)))
	}
	encodedDomain = append(encodedDomain, 0x00)
	return string(encodedDomain)
}

func BuildDNSQueryHeader() DNSHeader {
	return DNSHeader{
		ID:                uint16(rand.Intn(65535) + 1),
		NumberOfQuestions: 1,
		Flags:             1 << 8,
	}
}

// max for this function: 65535
func BuildDNSQuestion(domainName string) DNSQuestion {
	return DNSQuestion{
		Name:  EncodeDnsName(domainName),
		Type_: TYPE_A,
		Class: CLASS_IN,
	}
}

func (r *DNSRecord) ParseDNSHeader([]byte) {
}

func Main() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide a url")
		return 0
	}
	question := BuildDNSQuestion(os.Args[1])
	fmt.Printf("%#v\n", question)
	// chose 5 because it is the minimal length of a DNS resolver Question
	fmt.Println("question to bytes", question.ToBytes())
	return 0
}
