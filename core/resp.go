package core

import (
	"errors"
	"fmt"
)

// reads the length typically the first integer of the string
// until hit by an non-digit byte and returns
// the integer and the delta = length + 2 (CRLF)
// TODO: Make it simpler and read until we get '\r' just like other functions
func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

// reads a RESP encoded simple string from data and returns
// the string, the delth and the error
func readSimpleString(data []byte) (string, int, error) {
	// first character +
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil

}

// reads a RESP encoded error fromd ata and returns
// the error string, the delth and the error
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// reads a RESP encoded integer from data and returns
// the integer value, the delth and the error
func readInt64(data []byte) (int64, int, error) {
	// first character :
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// reads a RESP encoded string from data and returns
// the string, the delth and the error
func readBulkString(data []byte) (string, int, error) {
	// first character $
	pos := 1

	// reading the length and forwarding the position by
	// the length of the integer + the first special character
	length, delta := readLength(data[pos:])
	pos += delta

	// reading 'length' bytes as string
	return string(data[pos : pos+length]), pos + length + 2, nil
}

// reads a RESP encoded array from data and returns
// the decoded elements, the delta and the error
func readArray(data []byte) ([]interface{}, int, error) {
	// first character *
	pos := 1

	// reading the length
	count, delta := readLength(data[pos:])
	pos += delta

	elements := make([]interface{}, count)
	for i := range elements {
		element, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elements[i] = element
		pos += delta
	}
	return elements, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	ts, ok := value.([]interface{})
	if !ok {
		return nil, errors.New("resp value is not an array")
	}
	tokens := make([]string, len(ts))
	for i := range tokens {
		tokens[i] = ts[i].(string)
	}

	return tokens, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return []byte{}
}
