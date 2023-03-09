package uuid

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/gofrs/uuid/v5"
)

// UUIDSize is the size of the uuid
const UUIDSize = 16

var byteGroups = []int{8, 4, 4, 4, 12}

// OrderedUUID is an UUIDv1 with some bit shifted for database index performance
type UUID struct {
	Bytes [UUIDSize]byte
	Valid bool
}

func NewV1() UUID {
	v1, err := uuid.NewV1()
	if err != nil {
		return UUID{
			Valid: false,
		}
	}
	u := UUID{
		Valid: true,
	}
	copy(u.Bytes[0:], v1[0:])
	return u
}

func NewV4() UUID {
	v4, err := uuid.NewV4()
	if err != nil {
		return UUID{
			Valid: false,
		}
	}
	u := UUID{
		Valid: true,
	}
	copy(u.Bytes[0:], v4[0:])
	return u
}

// String returns the string representation of the UUID
func (u *UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u.Bytes[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u.Bytes[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u.Bytes[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u.Bytes[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u.Bytes[10:])

	return string(buf)
}

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	v := u.Bytes[0:16]
	return v, nil
}

// Scan implements the sql.Scanner interface.
func (u *UUID) Scan(src interface{}) error {
	u.Valid = false

	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case []byte:
		if len(src) != UUIDSize {
			return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(src))
		}
		copy(u.Bytes[0:], src[0:])
		u.Valid = true

		return nil
	}
	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//	"6ba7b8109dad11d180b400c04fd430c8"
//
// Supported UUID text representation follows:
//
//	uuid := canonical | hashlike
//	plain := canonical | hashlike
//	canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//	hashlike := 12hexoct
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 0:
		u.Valid = false
		u.Bytes = [16]byte{}
		return nil
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u.Bytes[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	u.Valid = true
	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u.Bytes[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}

	u.Valid = true
	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	if !u.Valid {
		return nil, nil
	}
	text = []byte(u.String())
	return
}

// FromString returns an OrderedUUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// FromBytes returns an OrderedUUID parsed from the byte array.
func FromBytes(input [16]byte) UUID {
	u := UUID{}
	copy(u.Bytes[0:], input[0:])

	u.Valid = true
	return u
}
