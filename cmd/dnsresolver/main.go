package main

import (
	"os"

	dnsresolver "github.com/joumanae/dsnresolver"
)

func main() {
	os.Exit(dnsresolver.Main())
}
