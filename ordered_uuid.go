package uuid

import (
	uuid "github.com/satori/go.uuid"
)

// NewOrdered creates a new UUIDv1 and converts it to an OrderedUUID
//  v1 UUID: aaaabbbb-cccc-dddd-1234-567890123456
//  transposed: ddddcccc-aaaa-bbbb-1234-567890123456
func NewOrdered() UUID {
	u := uuid.NewV1()

	return FromUUIDv1(u)
}

// DEPRECATED, use NewOrdered instead
func NewOrderedUUID() UUID {
	return NewOrdered()
}

// FromUUIDv1 populates Byte from uuidV1 and converts it
func FromUUIDv1(uuid uuid.UUID) UUID {
	buf := make([]byte, UUIDSize)

	copy(buf[0:2], uuid[6:8])
	copy(buf[2:4], uuid[4:6])
	copy(buf[4:6], uuid[0:2])
	copy(buf[6:8], uuid[2:4])
	copy(buf[8:], uuid[8:])

	u := UUID{}
	copy(u.Bytes[0:], buf[0:])
	u.Valid = true
	return u
}
