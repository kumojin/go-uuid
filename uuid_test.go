package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
const nilUUID = "00000000-0000-0000-0000-000000000000"

func TestUUIDValue(t *testing.T) {
	u := NewV1()
	val, err := u.Value()

	assert.Nil(t, err)
	assert.Equal(t, u.Bytes[:], val)
}

func TestUUIDValueNil(t *testing.T) {
	u := UUID{}
	val, err := u.Value()

	assert.Nil(t, err)
	assert.Nil(t, val)
}

func TestUUIDScan(t *testing.T) {
	u := UUID{}
	val := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	err := u.Scan(val)

	var expectedVal [16]byte
	copy(expectedVal[0:], val)

	assert.Nil(t, err)
	assert.Equal(t, expectedVal, u.Bytes)
	assert.True(t, u.Valid)
}

func TestUUIDScanNilValue(t *testing.T) {
	u := UUID{}
	err := u.Scan(nil)

	assert.Nil(t, err)
	assert.Equal(t, nilUUID, u.String())
	assert.False(t, u.Valid)
}

func TestUUIDScanInvalidType(t *testing.T) {
	u := UUID{}
	err := u.Scan("string")

	assert.NotNil(t, err)
	assert.Equal(t, "uuid: cannot convert string to UUID", err.Error())
	assert.Equal(t, nilUUID, u.String())
	assert.False(t, u.Valid)
}

func TestUUIDScanInvalidNumberOfBytes(t *testing.T) {
	val := []byte{0x6b, 0xa7, 0xb8, 0x10}

	u := UUID{}
	err := u.Scan(val)

	assert.NotNil(t, err)
	assert.Equal(t, "uuid: UUID must be exactly 16 bytes long, got 4 bytes", err.Error())
	assert.Equal(t, nilUUID, u.String())
	assert.False(t, u.Valid)

	val = []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8, 0xc8}
	err = u.Scan(val)
	assert.NotNil(t, err)
	assert.Equal(t, "uuid: UUID must be exactly 16 bytes long, got 17 bytes", err.Error())
	assert.Equal(t, nilUUID, u.String())
	assert.False(t, u.Valid)
}

func TestFromString(t *testing.T) {
	b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	emptyBytes := [16]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	// Valid uuid
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	s2 := "6ba7b8109dad11d180b400c04fd430c8"

	// With non valid hex char
	s3 := "6ba7b810zzzzzzzz80b400c04fd430c8"
	s4 := "6ba7b810-zzzz-zzzz-80b4-00c04fd430c8"

	// Without the dash or not at the right place
	s5 := "6ba7b81019dad111d1180b4100c04fd430c8"
	s6 := "6ba7b8109d-ad-11d1-80b4-00c04fd430c8"

	// Invalid string length
	s7 := "6ba7b810"

	u0, err := FromString("")
	assert.Nil(t, err)
	assert.Equal(t, emptyBytes, u0.Bytes)
	assert.False(t, u0.Valid)

	u1, err := FromString(s1)
	assert.Nil(t, err)
	assert.Equal(t, b, u1.Bytes)
	assert.True(t, u1.Valid)

	u2, err := FromString(s2)
	assert.Nil(t, err)
	assert.Equal(t, b, u2.Bytes)
	assert.True(t, u2.Valid)

	_, err = FromString(s3)
	assert.NotNil(t, err)

	_, err = FromString(s4)
	assert.NotNil(t, err)

	_, err = FromString(s5)
	assert.NotNil(t, err)

	_, err = FromString(s6)
	assert.NotNil(t, err)

	_, err = FromString(s7)
	assert.NotNil(t, err)
}

func TestMarshalText(t *testing.T) {
	b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	u := UUID{Bytes: b, Valid: true}

	expected := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	jsonString, err := u.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t, expected, jsonString)
}

func TestMarshalTextWithNotValidUUID(t *testing.T) {
	u := UUID{}

	expected := []byte(nil)

	jsonString, err := u.MarshalText()
	assert.Nil(t, err)
	assert.Equal(t, expected, jsonString)
}

func TestFromBytes(t *testing.T) {
	b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	u := FromBytes(b)

	assert.Equal(t, b, u.Bytes)
	assert.True(t, u.Valid)
}

func TestNewV1(t *testing.T) {
	u := NewV1()

	assert.True(t, u.Valid)
	assert.NotEqual(t, u.String(), nilUUID)
}

func TestNewV4(t *testing.T) {
	u := NewV4()

	assert.True(t, u.Valid)
	assert.NotEqual(t, u.String(), nilUUID)
}