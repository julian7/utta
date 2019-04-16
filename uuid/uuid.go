package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/pkg/errors"
)

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122.
type UUID [16]byte

// Nil empty UUID, all zeroes
var Nil UUID

// rander is the random function
var rander = rand.Reader

// New generates a UUID from using RNG
func New() (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(rander, uuid[:])
	if err != nil {
		return Nil, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid, nil
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// FromString fills up UUID from string representation
func (uuid *UUID) FromString(str string) error {
	b := []byte(str)
	if len(str) != 36 {
		return errors.New("invalid length")
	}
	if b[8] != '-' ||
		b[13] != '-' ||
		b[18] != '-' ||
		b[23] != '-' {
		return errors.Errorf("invalid UUID format: %v", str)
	}
	for i, x := range [5]struct {
		size  uint8
		index uint8
		uidx  uint8
	}{
		{8, 0, 0},
		{4, 9, 4},
		{4, 14, 6},
		{4, 19, 8},
		{12, 24, 10},
	} {
		_, err := hex.Decode(uuid[x.uidx:x.uidx+(x.size/2)], b[x.index:x.index+x.size])
		if err != nil {
			return errors.Wrapf(err, "invalid UUID format at %d", i)
		}
	}
	return nil
}
