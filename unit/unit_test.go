package unit

import (
	"testing"
)

func TestParseBytes(t *testing.T) {
	testCases := []struct {
		s   string
		v   uint64
		iec bool
	}{
		{"1B", 1, false},
		{"1KB", 1000, false},
		{"1MB", 1000000, false},
		{"1GB", 1000000000, false},
		{"1TB", 1000000000000, false},
		{"1PB", 1000000000000000, false},
		{"1EB", 1000000000000000000, false},
		{"1b", 1, false},
		{"1kb", 1000, false},
		{"1mb", 1000000, false},
		{"1gb", 1000000000, false},
		{"1tb", 1000000000000, false},
		{"1pb", 1000000000000000, false},
		{"1eb", 1000000000000000000, false},
		{"1K", 1000, false},
		{"1M", 1000000, false},
		{"1G", 1000000000, false},
		{"1T", 1000000000000, false},
		{"1P", 1000000000000000, false},
		{"1E", 1000000000000000000, false},
		{"1KIB", 1024, true},
		{"1MIB", 1048576, true},
		{"1GIB", 1073741824, true},
		{"1TIB", 1099511627776, true},
		{"1PIB", 1125899906842624, true},
		{"1EIB", 1152921504606846976, true},
		{"1kib", 1024, true},
		{"1mib", 1048576, true},
		{"1gib", 1073741824, true},
		{"1tib", 1099511627776, true},
		{"1pib", 1125899906842624, true},
		{"1eib", 1152921504606846976, true},
		{"1Ki", 1024, true},
		{"1Mi", 1048576, true},
		{"1Gi", 1073741824, true},
		{"1Ti", 1099511627776, true},
		{"1Pi", 1125899906842624, true},
		{"1Ei", 1152921504606846976, true},
		{"1", 1, false},
		{"1000000", 1000000, false},
		{"1,048,576", 1048576, false},
		{"1048576", 1048576, false},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			bytes, iec, err := ParseBytes(tc.s)
			if err != nil {
				t.Errorf("ParseBytes error => %v", err)
				return
			}
			if bytes != tc.v {
				t.Errorf("expect %d, but get %d", tc.v, bytes)
				return
			}
			if iec != tc.iec {
				t.Errorf("IEC expect %v, but get %v", tc.iec, iec)
				return
			}

			bytesInt, err := ParseBytesInt(tc.s)
			if err != nil {
				t.Errorf("ParseBytesInt error => %v", err)
				return
			}
			if bytesInt != int(tc.v) {
				t.Errorf("expect %d, but get %d", tc.v, bytesInt)
				return
			}
		})
	}
}

func TestParseBytesReturnError(t *testing.T) {
	testCases := []struct {
		s string
	}{
		{"1x"},
		{"a"},
		{"1MB."},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			_, _, err := ParseBytes(tc.s)
			if err == nil {
				t.Errorf("ParseBytes expect get an error, but get nil")
			}

			_, err = ParseBytesInt(tc.s)
			if err == nil {
				t.Errorf("ParseBytesInt expect get an error, but get nil")
			}
		})
	}
}

func TestBytes(t *testing.T) {
	testCases := []struct {
		s string
		v uint64
	}{
		{"10 B", 10},
		{"10 kB", 10000},
		{"10 MB", 10000000},
		{"10 GB", 10000000000},
		{"10 TB", 10000000000000},
		{"10 PB", 10000000000000000},
		{"10 EB", 10000000000000000000},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			s := Bytes(tc.v)
			if s != tc.s {
				t.Errorf("expect %s, but get %s", tc.s, s)
			}
		})
	}
}

func TestIBytes(t *testing.T) {
	testCases := []struct {
		s string
		v uint64
	}{
		{"10 B", 10},
		{"10 KiB", 10240},
		{"10 MiB", 10485760},
		{"10 GiB", 10737418240},
		{"10 TiB", 10995116277760},
		{"10 PiB", 11258999068426240},
		{"10 EiB", 11529215046068469760},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			s := IBytes(tc.v)
			if s != tc.s {
				t.Errorf("expect %s, but get %s", tc.s, s)
			}
		})
	}
}
