package dnsresolver

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
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
	return header.HeaderToBytesToHexadecimal(DNSHeaderSize) + string(question.QuestionToBytes(12))

}

func (r *DNSRecord) ParseDNSHeader([]byte) {

}

func Main() int {
	url := flag.String("url", "", "url to resolve")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Please provide a url")
		return 0
	}
	h := DNSHeader{
		ID:                  0x1314,
		Flags:               0,
		NumberOfQuestions:   1,
		NumberOfAnswers:     0,
		NumberOfAuthorities: 0,
		NumberOfAdditionals: 0,
	}
	h.HeaderToBytesToHexadecimal(DNSHeaderSize)
	// Example usage
	q := DNSQuestion{
		Name:  *url,
		Type_: 1, // A record
		Class: 1, // IN class
	}
	// chose 5 because it is the minimal length of a DNS resolver Question
	fmt.Println("question to bytes", q.QuestionToBytes(5))
	fmt.Println("Building the query", BuildDNSQuery(*url, "1"), len(BuildDNSQuery("example.com", "1")))

	return 0
}
