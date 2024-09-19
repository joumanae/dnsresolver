package dnsresolver_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"slices"
	"testing"

	dnsresolver "github.com/joumanae/dsnresolver"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestHeaderToBytesReturnsCorrectFormat(t *testing.T) {

	got := dnsresolver.BuildDNSQueryHeader().ToBytes()
	want := []byte{0x1, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	if !slices.Equal(want, got[2:]) {
		t.Fatalf("the dns resolver gave back an incorrect header. want %v, got %v", want, got)
	}
}

func TestQueryToBytes(t *testing.T) {
	// Create a DNSQuery instance
	question := dnsresolver.DNSQuery{
		Name:  []byte("example.com"),
		Type_: 1,
		Class: 1,
	}

	// Call the QueryToBytes function

	got := question.ToBytes()
	// Create the expected byte slice
	var expected bytes.Buffer
	binary.Write(&expected, binary.BigEndian, []byte("example.com"))
	binary.Write(&expected, binary.BigEndian, uint16(1))
	binary.Write(&expected, binary.BigEndian, uint16(1))

	// Compare the expected and got byte slices
	if !bytes.Equal(got, expected.Bytes()) {
		t.Errorf("QueryToBytes returned unexpected result. got %v, want %v", got, expected.Bytes())
	}
}

func TestEncodeDnsName(t *testing.T) {
	// Create a domain name
	domainName := "example.com"
	got := dnsresolver.EncodeDnsName(domainName)
	want := []byte{0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00}
	if !slices.Equal(got, want) {
		t.Errorf("EncodeDnsName returned unexpected result. got %#v, want %#v", got, want)
	}
}

func TestBuildDNSQueryHeader_ReturnsCorrectHeader(t *testing.T) {
	t.Parallel()
	header := dnsresolver.BuildDNSQueryHeader()
	if header.ID == 0 {
		t.Fatal("header ID should not be zero")
	}
	if header.NumberOfQuerys != 1 {
		t.Fatalf("number of questions should be 1, got %d", header.NumberOfQuerys)
	}
	if header.Flags&0b0000000_100000000 == 0 {
		t.Fatalf("recursionDesired flag should be set")
	}
}

func TestBuildDNSQuery_ReturnsCorrectString(t *testing.T) {
	//
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dnsresolver": dnsresolver.Main,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}
