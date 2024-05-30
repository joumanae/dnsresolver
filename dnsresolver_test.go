package dnsresolver_test

import (
	"testing"

	dnsresolver "github.com/joumanae/dsnresolver"
)

func TestHeaderToBytesReturnsCorrectFormat(t *testing.T) {
	var dnsheader dnsresolver.DNSHeader
	got := dnsheader.HeaderToBytesToHexadecimal(12)
	want := "\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00"

	if want != got {
		t.Fatalf("the dns resolver gave back an incorrect header. want %v, got %v", want, got)
	}
}
