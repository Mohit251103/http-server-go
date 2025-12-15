package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func isValid(b byte) bool {
	specialChars := []byte("!#$%&'*+-.^_`|~")
	if (b < '0' || b > '9') && (b < 'a' || b > 'z') && bytes.IndexByte(specialChars, b) == -1 {
		return false
	}

	return true
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	} else if idx == 0 {
		return 0, true, nil
	}

	data = data[:idx]
	n = len(string(data) + "\r\n")
	data = bytes.TrimSpace(data)
	data = bytes.ToLower(data)

	key_value := bytes.SplitN(data, []byte(":"), 2)
	key := key_value[0]
	value := bytes.TrimSpace(key_value[1])

	for _, r := range key {
		if !isValid(r) {
			return 0, false, fmt.Errorf("invalid character present in key, %c", r)
		}
	}
	if bytes.Contains(key, []byte(" ")) {
		return 0, false, fmt.Errorf("invalid format, cannot add space after key and before semicolon")
	}

	if h[string(key)] != "" {
		h[string(key)] = h[string(key)] + ", " + string(value)
	} else {
		h[string(key)] = string(value)
	}
	return n, false, nil
}

func NewHeaders() (h Headers) {
	return Headers{}
}
