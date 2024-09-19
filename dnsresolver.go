package dnsresolver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	Name  []byte
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
		Flags:          1<<8 | 1<<7, // Set recursion desired (RD) bit
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

	r.Name = q.Name
	r.TTL = 0

}

func ParseRecord(r io.Reader) (DNSRecord, error) {
	rec := DNSRecord{}
	var nameLen byte
	for {
		if err := binary.Read(r, binary.BigEndian, &nameLen); err != nil {
			return rec, err
		}
		if nameLen == 0 {
			break
		}
		name := make([]byte, nameLen)
		_, err := io.ReadFull(r, name)
		if err != nil {
			return rec, err
		}
		rec.Name = append(rec.Name, name...)
		rec.Name = append(rec.Name, byte('.'))
	}
	if err := binary.Read(r, binary.BigEndian, &rec.Type); err != nil {
		return rec, err
	}
	if err := binary.Read(r, binary.BigEndian, &rec.Class); err != nil {
		return rec, err
	}
	if err := binary.Read(r, binary.BigEndian, &rec.TTL); err != nil {
		return rec, err
	}
	if err := binary.Read(r, binary.BigEndian, &rec.RDLen); err != nil {
		return rec, err
	}
	var addr [4]byte
	if err := binary.Read(r, binary.BigEndian, &addr); err != nil {
		return rec, err
	}
	rec.RData = addr[:]
	return rec, nil
}

func Main() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please provide a url")
		return 0
	}
	header := BuildDNSQueryHeader()
	fmt.Printf("Query header: %#v\n", header)
	query := BuildDNSQuery(os.Args[1])
	fmt.Printf("Query:\n\tName: %q\n\tType: %x\n\tClass: %x\n", query.Name, query.Type_, query.Class)
	message := append(header.ToBytes(), query.ToBytes()...)
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer conn.Close()
	buf := make([]byte, 104)
	w, err := conn.Write(message)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("wrote %d bytes\n", w)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("data: %#v\n", buf)
	resp := bytes.NewReader(buf)
	header, err = ParseHeader(resp)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("header: %#v\n", header)
	record, err := ParseRecord(resp)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("name: %q\n", record.Name)
	fmt.Printf("record: %#v\n", record)
	return 0
}
