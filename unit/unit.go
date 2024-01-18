package unit

import "github.com/dustin/go-humanize"

// ParseBytes parse the string to bytes
func ParseBytes(s string) (int, error) {
	b, err := humanize.ParseBytes(s)
	if err != nil {
		return 0, err
	}
	return int(b), nil
}
