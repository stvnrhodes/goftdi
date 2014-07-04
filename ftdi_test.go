package ftdi_test

import "testing"

func TestVersion(t *testing.T) {
	want := ftdi.Version{Major: 1}
	if v := ftdi.GetVersion(); v != want {
		t.Error("Got %v, want %v", v, want)
	}
}
