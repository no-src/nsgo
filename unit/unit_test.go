package unit

import "testing"

func TestParseBytes(t *testing.T) {
	testCases := []struct {
		s string
		v int
	}{
		{"1B", 1},
		{"1MB", 1000000},
		{"1GB", 1000000000},
		{"1TB", 1000000000000},
		{"1b", 1},
		{"1mb", 1000000},
		{"1gb", 1000000000},
		{"1tb", 1000000000000},
		{"1MIB", 1048576},
		{"1GIB", 1073741824},
		{"1TIB", 1099511627776},
		{"1mib", 1048576},
		{"1gib", 1073741824},
		{"1tib", 1099511627776},
		{"1", 1},
		{"1000000", 1000000},
		{"1,048,576", 1048576},
		{"1048576", 1048576},
	}
	for _, tc := range testCases {
		t.Run(tc.s, func(t *testing.T) {
			bytes, err := ParseBytes(tc.s)
			if err != nil {
				t.Errorf("parse bytes error => %v", err)
				return
			}
			if bytes != tc.v {
				t.Errorf("expect %d, but get %d", tc.v, bytes)
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
			_, err := ParseBytes(tc.s)
			if err == nil {
				t.Errorf("expect get an error, but get nil")
			}
		})
	}
}
