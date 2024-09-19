package dnsresolver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
)

const (
	TYPE_A   = 1
	CLASS_IN = 1
)

type DNSHeader struct {
	ID                  uint16
	Flags               uint16
	NumberOfQuerys      uint16
	NumberOfAuthorities uint16
	NumberOfAdditionals uint16
	NumberOfAnswers     uint16
}

const DNSHeaderSize = 12

type DNSQuery struct {
	Name  []byte
	Type_ uint16
	Class uint16
}

type DNSRecord struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32 // time to live
	RDLen uint16
	RData []byte
}

type DNSPacket struct {
	Header      DNSHeader
	Querys      []DNSQuery
	Answers     []DNSRecord
	Authorities []DNSRecord
	Additionals []DNSRecord
}

func (h DNSHeader) ToBytes() []byte {
	headerSlice := make([]byte, DNSHeaderSize)

	// BigEndian saves the most significant piece of data at the lowest in memory
	// Endian is useful when data isn't single-byte. In our case, each element is multiple bytes
	binary.BigEndian.PutUint16(headerSlice[0:2], h.ID)
	binary.BigEndian.PutUint16(headerSlice[2:4], h.Flags)
	binary.BigEndian.PutUint16(headerSlice[4:6], h.NumberOfQuerys)
	binary.BigEndian.PutUint16(headerSlice[6:8], h.NumberOfAuthorities)
	binary.BigEndian.PutUint16(headerSlice[8:10], h.NumberOfAdditionals)
	binary.BigEndian.PutUint16(headerSlice[10:12], h.NumberOfAnswers)
	return headerSlice
}

func (q *DNSQuery) ToBytes() []byte {
	questionSlice := []byte(q.Name)

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, q.Type_)
	binary.Write(buffer, binary.BigEndian, q.Class)

	return append(questionSlice, buffer.Bytes()...)
}

func EncodeDnsName(domainName string) []byte {
	var encodedDomain []byte
	parts := bytes.Split([]byte(domainName), []byte("."))
	for _, part := range parts {
		encodedDomain = append(encodedDomain, byte(len(part)))
		encodedDomain = append(encodedDomain, part...)
	}
	encodedDomain = append(encodedDomain, 0x00)
	return encodedDomain
}

func BuildDNSQueryHeader() DNSHeader {
	return DNSHeader{
		ID:             uint16(rand.Intn(65535) + 1),
		NumberOfQuerys: 1,
		Flags:          1 << 8,
	}
}

// max for this function: 65535
func BuildDNSQuery(domainName string) DNSQuery {
	return DNSQuery{
		Name:  EncodeDnsName(domainName),
		Type_: TYPE_A,
		Class: CLASS_IN,
	}
}

func ParseHeader(r *bytes.Reader) (DNSHeader, error) {
	h := DNSHeader{}
	DNSFields := []interface{}{
		&h.ID,
		&h.Flags,
		&h.NumberOfQuerys,
		&h.NumberOfAuthorities,
		&h.NumberOfAdditionals,
		&h.NumberOfAnswers,
	}
	for _, field := range DNSFields {
		if err := binary.Read(r, binary.BigEndian, field); err != nil {
			return h, err
		}
	}
	return h, nil
}

// TODO: Need to parse and print the query, then parse and print the reord
func (r *DNSRecord) ParseDNSQuery([]byte) {
	q := DNSQuery{}
	r.Type = q.Type_
	r.Class = q.Class

	r.Name = string(q.Name)
	r.TTL = 0

}

// ParseDNSRecord takes a byte slice and parses it into a DNSRecord.
// It then copies the fields from the parsed DNSRecord into the receiver.
func (r *DNSRecord) ParseDNSRecord(b []byte) {
	rec := DNSRecord{}
	// The following line is a placeholder until we implement the actual parsing of the DNSRecord
	fmt.Println(rec)
}

func Main() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide a url")
		return 0
	}
	header := BuildDNSQueryHeader()
	query := BuildDNSQuery(os.Args[1])
	message := append(header.ToBytes(), query.ToBytes()...)
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer conn.Close()
	buf := make([]byte, 12)
	w, err := conn.Write(message)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("wrote %d bytes\n", w)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	header, err = ParseHeader(bytes.NewReader(buf))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("header: %#v\n", header)

	fmt.Printf("read %d bytes, bytes:%#v\n", n, buf)
	return 0
}
