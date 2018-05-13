// Package tuid implements the Time-based Unique Identifier.
package tuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// TUID represents a Time-based Unique Identifier.
type TUID uint64

// Size represents the size of a TUID in bytes.
const Size = 8

// epoch is the custom epoch (2018-01-01 00:00:00.0 UTC).
var epoch = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

var (
	// ErrIncomplete occurs when parsing fails.
	ErrIncomplete = errors.New("tuid: incomplete or truncated")
)

// New returns a new Time-based Unique Identifier.
func New() (v TUID) {
	var r [3]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	v = TUID(time.Now().Sub(epoch).Nanoseconds()/1000000) << (64 - 42)
	v |= TUID(r[0]) | TUID(r[1])<<8 | TUID(r[2]&0x3f)<<16
	return
}

// FromBytes returns the TUID converted from a raw byte slice input. It will return ErrIncomplete if the slice isn't Size bytes long.
func FromBytes(b []byte) (TUID, error) {
	if len(b) != Size {
		return 0, ErrIncomplete
	}
	return TUID(b[7]) | TUID(b[6])<<8 | TUID(b[5])<<16 | TUID(b[4])<<24 |
		TUID(b[3])<<32 | TUID(b[2])<<40 | TUID(b[1])<<48 | TUID(b[0])<<56, nil
}

// FromString returns the TUID parsed from a hexadecimal string input.
func FromString(s string) (TUID, error) {
	var b [Size]byte
	if _, err := hex.Decode(b[:], []byte(s)); err != nil {
		return 0, err
	}
	return FromBytes(b[:])
}

// Array returns the byte array representation of TUID.
func (v TUID) Array() (a [Size]byte) {
	a[0] = byte(v >> 56)
	a[1] = byte(v >> 48)
	a[2] = byte(v >> 40)
	a[3] = byte(v >> 32)
	a[4] = byte(v >> 24)
	a[5] = byte(v >> 16)
	a[6] = byte(v >> 8)
	a[7] = byte(v)
	return
}

// Bytes returns the bytes slice representation of TUID.
func (v TUID) Bytes() []byte {
	b, a := make([]byte, Size), v.Array()
	copy(b, a[:])
	return b
}

// String returns the hexadecimal string representation of TUID.
func (v TUID) String() string {
	a := v.Array()
	return hex.EncodeToString(a[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (v TUID) MarshalText() ([]byte, error) {
	a, text := v.Array(), make([]byte, Size*2)
	hex.Encode(text, a[:])
	return text, nil
}

// UnmarshalText parses a hexadecimal representation of a TUID into v.
func (v *TUID) UnmarshalText(text []byte) (err error) {
	var b [Size]byte
	if _, err = hex.Decode(b[:], text); err != nil {
		return
	}
	*v, err = FromBytes(b[:])
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (v TUID) MarshalBinary() (data []byte, err error) {
	return v.Bytes(), nil
}

// UnmarshalBinary parses a raw byte slice representation of a TUID into v.
func (v *TUID) UnmarshalBinary(data []byte) (err error) {
	*v, err = FromBytes(data)
	return
}

// Scan implements the sql.Scanner interface.
func (v *TUID) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		if len(src) == Size {
			*v, err = FromBytes(src)
		} else {
			err = v.UnmarshalText(src)
		}

	case string:
		*v, err = FromString(src)

	case int64:
		*v = TUID(src)

	default:
		err = fmt.Errorf("tuid: cannot convert %T to TUID", src)
	}

	return
}

// Value implements the driver.Valuer interface. The driver value is int64.
func (v TUID) Value() (driver.Value, error) {
	return int64(v), nil
}

// Time returns the timestamp of the TUID.
func (v TUID) Time() time.Time {
	return epoch.Add(time.Duration(v>>(64-42)) * time.Millisecond)
}
