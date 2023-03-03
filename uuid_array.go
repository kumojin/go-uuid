package uuid

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type UUIDArray []UUID

func (a *UUIDArray) Scan(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*a = []UUID{}
		return nil
	case []byte:
		return json.Unmarshal(v, &a)
	case string:
		return json.Unmarshal([]byte(v), &a)
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (a UUIDArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
