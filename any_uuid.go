package uuid

import "database/sql/driver"

// AnyUUID For easy UUID matching for mocks and tests
type AnyUUID struct{}

// Match if the UUID is an array of 16 bytes
func (a AnyUUID) Match(v driver.Value) bool {
	b, ok := v.([]byte)
	if !ok {
		return false
	}
	return len(b) == 16
}

