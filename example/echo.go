package main

import (
	"io"
	"os"

	"github.com/stvnrhodes/goftdi"
)

func main() {
	cfg := ftdi.Config{Vendor: 0x0403, Product: 0xFFA8, Baud: 115200}
	c, err := ftdi.Open(cfg)
	if err != nil {
		panic(err)
	}
	go io.Copy(c, os.Stdin)
	io.Copy(os.Stdout, c)
}
