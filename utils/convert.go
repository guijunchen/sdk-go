package utils

import (
	"encoding/binary"
	"errors"
)

func U64ToBytes(i uint64) []byte {
	var b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func BytesToU64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, errors.New("invalid uint64 bytes")
	}
	return binary.LittleEndian.Uint64(b), nil
}
