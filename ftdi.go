package ftdi

import "io"

// #include <ftdi.h>
import "C"

// Version represents a version of libftdi.
type Version struct {
	Major, Minor, Micro int
	Version             string // Version as (static) string
	Snapshot            string //Git snapshot version if known. Otherwise "unknown" or empty string.
}

// GetVersion gets the current libftdi version.
func GetVersion() Version {
	info := C.ftdi_get_library_version()
	return Version{
		Major:    info.major,
		Minor:    info.minor,
		Micro:    info.micro,
		Version:  C.GoString(info.version_str),
		Snapshot: C.GoString(info.snapshot_str),
	}
}

type Config struct {
}
type conn struct {
	ctx *C.struct_ftdi_context
}

func Open(c *Config) (io.ReadWriteCloser, error) {
	ptr := C.ftdi_new()
	if ptr == nil {
		return
	}
	C.ftdi_usb_open()
}
