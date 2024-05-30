package dnsresolver_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	dnsresolver "github.com/joumanae/dsnresolver"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestHeaderToBytesReturnsCorrectFormat(t *testing.T) {
	var dnsheader dnsresolver.DNSHeader
	got := dnsheader.HeaderToBytesToHexadecimal(12)
	want := "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00"

	if want != got {
		t.Fatalf("the dns resolver gave back an incorrect header. want %v, got %v", want, got)
	}
}
func TestQuestionToBytes(t *testing.T) {
	// Create a DNSQuestion instance
	question := dnsresolver.DNSQuestion{
		Name:  "example.com",
		Type_: 1,
		Class: 1,
	}

	// Call the QuestionToBytes function
	got := question.QuestionToBytes(len(question.Name) + 4)

	// Create the expected byte slice
	var expected bytes.Buffer
	binary.Write(&expected, binary.BigEndian, []byte("example.com"))
	binary.Write(&expected, binary.BigEndian, uint16(1))
	binary.Write(&expected, binary.BigEndian, uint16(1))

	// Compare the expected and got byte slices
	if !bytes.Equal(got, expected.Bytes()) {
		t.Errorf("QuestionToBytes returned unexpected result. got %v, want %v", got, expected.Bytes())
	}
}

func TestEncodeDnsName(t *testing.T) {
	// Create a domain name
	domainName := "example.com"
	got := dnsresolver.EncodeDnsName(domainName)
	want := []byte{7, 3, 0}
	if !bytes.Equal(got, want) {
		t.Errorf("EncodeDnsName returned unexpected result. got %v, want %v", got, want)
	}
}

func TestBuildDNSQuery(t *testing.T) {
	// Create a domain name
	domainName := "example.com"
	// Create a record type
	recordType := "1"
	got := len(dnsresolver.BuildDNSQuery(domainName, recordType))
	want := 55
	if got != want {
		t.Errorf("BuildDNSQuery returned unexpected result. got %v, want %v", got, want)
	}
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
