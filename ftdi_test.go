package ftdi_test

import "testing"
import "github.com/stvnrhodes/goftdi"

// TestVersion is a bit fragile, it tests for the exact version so library
// updates will break it.
func TestVersion(t *testing.T) {
	want := ftdi.Version{Major: 1, Minor: 1, Micro: 0, Version: "1.1", Snapshot: "unknown"}
	if v := ftdi.GetVersion(); v != want {
		t.Error("Got %v, want %v", v, want)
	}
}

// TestOpen only works if there's actually a device to open
func TestOpen(t *testing.T) {
	conn, err := ftdi.Open(ftdi.Config{Vendor: 0x1234, Product: 0x5678, Baud: 9600})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Write([]byte{1, 2, 3}); err != nil {
		t.Error(err)
	}
}
