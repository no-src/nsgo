package randutil

import (
	"crypto/rand"
	"errors"
	"testing"
)

func TestRandomString(t *testing.T) {
	randLen := 10
	randStr := RandomString(randLen)
	if len(randStr) != randLen {
		t.Errorf("test RandomString error, expect len:%d, actual:%s", randLen, randStr)
	}
}

func TestRandomString_MoreThanMaxLength(t *testing.T) {
	randLen := 30
	maxLen := 20
	randStr := RandomString(randLen)
	if len(randStr) != maxLen {
		t.Errorf("test RandomString error, expect len:%d, actual:%s", randLen, randStr)
	}
}

func TestRandomString_WithReadError(t *testing.T) {
	read = func(b []byte) (n int, err error) {
		return 0, errors.New("read error test")
	}
	defer func() {
		read = rand.Read
	}()
	randLen := 10
	randStr := RandomString(randLen)
	if len(randStr) != randLen {
		t.Errorf("test RandomString error, expect len:%d, actual:%s", randLen, randStr)
	}

	randStr2 := RandomString(randLen)
	if len(randStr2) != randLen {
		t.Errorf("test RandomString error, expect len:%d, actual:%s", randLen, randStr2)
	}

	if randStr == randStr2 {
		t.Errorf("test RandomString error, the first random string (%s) is the same as the second (%s)", randStr, randStr2)
	}
}
