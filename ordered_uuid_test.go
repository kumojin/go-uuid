package uuid

import (
	"encoding/hex"
	"github.com/gofrs/uuid/v5"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderedUUIDShouldReturnANotNullUUID(t *testing.T) {
	u := NewOrdered()

	assert.NotEqual(t, nilUUID, u.String())
}

func TestOrderedUUIDFromUUIDv1ShouldConvertItCorrectly(t *testing.T) {
	u := uuid.UUID{}
	data, _ := hex.DecodeString("aaaabbbbccccdddd1234567890123456")
	copy(u[0:], data)

	o := FromUUIDv1(u)

	assert.Equal(t, "aaaabbbb-cccc-dddd-1234-567890123456", u.String())
	assert.Equal(t, "ddddcccc-aaaa-bbbb-1234-567890123456", o.String())
}
