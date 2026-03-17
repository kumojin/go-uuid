# Migration guide: go-uuid to gofrs/uuid/v5

This guide explains how to replace `github.com/kumojin/go-uuid/v2` with [
`github.com/gofrs/uuid/v5`](https://github.com/gofrs/uuid).

## Why migrate?

This library is a thin wrapper around `gofrs/uuid/v5`. Its main value was `OrderedUUID` — a custom byte-reordering of
UUID v1 for database index performance. [RFC 9562](https://www.rfc-editor.org/rfc/rfc9562) (published May 2024)
standardized UUID v7, which solves the same problem in an interoperable way. `gofrs/uuid/v5` already supports UUID v7
natively.

## API mapping

| go-uuid              | gofrs/uuid/v5                 |
|----------------------|-------------------------------|
| `uuid.NewV1()`       | `uuid.NewV1()`                |
| `uuid.NewV4()`       | `uuid.NewV4()`                |
| `uuid.NewOrdered()`  | `uuid.NewV7()` (for new data) |
| `uuid.FromString(s)` | `uuid.FromString(s)`          |
| `uuid.FromBytes(b)`  | `uuid.FromBytes(b[:])`        |
| `u.String()`         | `u.String()`                  |
| `u.MarshalText()`    | `u.MarshalText()`             |
| `u.UnmarshalText(t)` | `u.UnmarshalText(t)`          |
| `u.Scan(src)`        | `u.Scan(src)`                 |
| `u.Value()`          | `u.Value()`                   |

### Key differences

- `gofrs/uuid.UUID` is a `[16]byte` directly, not a struct with `Bytes` and `Valid` fields. There is no built-in
  nullability — use `*uuid.UUID` or `uuid.NullUUID` instead.
- `gofrs/uuid.FromBytes` takes a `[]byte` slice, not a `[16]byte` array.
- `gofrs/uuid.NewV7()` returns `(uuid.UUID, error)` — handle the error explicitly instead of relying on the `Valid`
  field.

## Handling existing OrderedUUID data

OrderedUUID uses a **non-standard byte layout**: it rearranges the timestamp fields of a UUID v1 so that the most
significant time bits come first. This format is incompatible with UUID v7.

If you have OrderedUUID values stored in a database, you have two options.

### Option 1: Migrate existing data (recommended)

Convert all stored OrderedUUIDs to standard UUID v4 or leave them as-is (they are valid 16-byte identifiers — they just
won't sort chronologically under a standard UUID interpretation). Then use UUID v7 for all new records.

This works well when:

- Chronological ordering of old records is not important
- You can run a data migration

In practice, this means:

1. Switch all code that **generates** UUIDs to use `uuid.NewV7()`
2. Keep reading existing values as raw `uuid.UUID` — they'll parse fine since they're still valid 16-byte values, they
   just won't sort by creation time
3. Over time, old records with OrderedUUID values get replaced by new records with UUID v7 values

### Option 2: Support both formats during transition

If you need to preserve the chronological ordering semantics of existing OrderedUUIDs (e.g., for range queries on old
data), keep a small conversion utility alongside your new code.

```go
package uuidutil

import "github.com/gofrs/uuid/v5"

// FromOrderedUUID reverses the byte reordering applied by go-uuid's
// NewOrdered(). This recovers the original UUID v1 from an OrderedUUID,
// preserving its timestamp for sorting or comparison purposes.
//
// OrderedUUID layout:  ddddcccc-aaaa-bbbb-1234-567890123456
// Original v1 layout:  aaaabbbb-cccc-dddd-1234-567890123456
func FromOrderedUUID(ordered uuid.UUID) uuid.UUID {
	var original uuid.UUID
	copy(original[0:2], ordered[4:6])
	copy(original[2:4], ordered[6:8])
	copy(original[4:6], ordered[2:4])
	copy(original[6:8], ordered[0:2])
	copy(original[8:], ordered[8:])
	return original
}

// ToOrderedUUID applies the go-uuid byte reordering to a UUID v1.
// Use this only if you need to write values compatible with existing
// OrderedUUID data.
func ToOrderedUUID(v1 uuid.UUID) uuid.UUID {
	var ordered uuid.UUID
	copy(ordered[0:2], v1[6:8])
	copy(ordered[2:4], v1[4:6])
	copy(ordered[4:6], v1[0:2])
	copy(ordered[6:8], v1[2:4])
	copy(ordered[8:], v1[8:])
	return ordered
}
```

Then plan a migration to Option 1 when feasible.

## Handling UUIDArray

`UUIDArray` serializes/deserializes a `[]UUID` as JSON for database storage. Replace it with a small custom type:

```go
package uuidutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(val any) error {
	switch v := val.(type) {
	case nil:
		*a = []uuid.UUID{}
		return nil
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (a UUIDArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
```

## Handling nullability

go-uuid uses a `Valid` field to represent null UUIDs. With `gofrs/uuid/v5`, use `uuid.NullUUID`:

```go
// Before (go-uuid)
type User struct {
ID       uuid.UUID // Valid == false means NULL
}

// After (gofrs/uuid/v5)
type User struct {
ID       uuid.NullUUID
}

// Or simply use a pointer
type User struct {
ID       *uuid.UUID
}
```
