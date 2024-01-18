package unit

import (
	"strings"

	"github.com/dustin/go-humanize"
)

// ParseBytes parse the string to bytes
func ParseBytes(s string) (bytes uint64, iec bool, err error) {
	bytes, err = humanize.ParseBytes(s)
	if err != nil {
		return bytes, iec, err
	}
	iec = isIEC(s)
	return bytes, iec, nil
}

// Bytes produces a human readable representation of an SI size
func Bytes(b uint64) string {
	return humanize.Bytes(b)
}

// IBytes produces a human readable representation of an IEC size
func IBytes(b uint64) string {
	return humanize.IBytes(b)
}

func isIEC(s string) bool {
	return strings.ContainsRune(strings.ToLower(s), 'i')
}
