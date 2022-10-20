package mssql

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Size of a UUID in bytes.
const Size = 16

// Size of a canonical representation of UUID in bytes.
const CanonicalSize = 36

// UUID representation compliant with specification described in RFC 4122.
type UUID [Size]byte

// UUID versions
const (
	_ byte = iota
	V4
)

// UUID layout variants.
const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// Nil is special form of UUID that is specified to have all 128 bits set to zero.
var NilUUID = UUID{}

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 UUID, u2 UUID) bool {
	return bytes.Equal(u1[:], u2[:])
}

// Version returns algorithm version used to generate UUID.
func (u UUID) Version() byte {
	return u[6] >> 4
}

// Variant returns UUID layout variant.
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VariantNCS
	case (u[8] >> 6) == 0x02:
		return VariantRFC4122
	case (u[8] >> 5) == 0x06:
		return VariantMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VariantFuture
	}
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u *UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// SetVersion sets version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits.
func (u *UUID) SetVariant(v byte) {
	switch v {
	case VariantNCS:
		u[8] = u[8]&(0xff>>1) | (0x00 << 7)
	case VariantRFC4122:
		u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	case VariantMicrosoft:
		u[8] = u[8]&(0xff>>3) | (0x06 << 5)
	case VariantFuture:
		fallthrough
	default:
		u[8] = u[8]&(0xff>>3) | (0x07 << 5)
	}
}

// Can be used to split slices of UUIDs into multiple subslices with a defined size.
func Batch(uuids []UUID, size int, onBatch func(uuids []UUID) error) error {
	for i := 0; i < len(uuids); i += size {
		var end int = i + size
		if end > len(uuids) {
			end = len(uuids)
		}

		var err = onBatch(uuids[i:end])
		if err != nil {
			return err
		}
	}

	return nil
}

// FromBytes returns UUID converted from raw byte slice input.
//
// It will return error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

// FromBytesOrNil returns UUID converted from raw byte slice input.
//
// Same behavior as FromBytes, but returns a Nil UUID on error.
func FromBytesOrNil(input []byte) UUID {
	var uuid, err = FromBytes(input)

	if err != nil {
		return NilUUID
	}

	return uuid
}

// FromString returns UUID parsed from string input.
//
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// FromStringOrNil returns UUID parsed from string input.
//
// Same behavior as FromString but returns a Nil UUID on error.
func FromStringOrNil(input string) UUID {
	var uuid UUID

	var err error
	uuid, err = FromString(input)
	if err != nil {
		return NilUUID
	}

	return uuid
}

// MarshalText implements the encoding.TextMarshaler interface.
//
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
//
// Following format is supported "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
func (u *UUID) UnmarshalText(t []byte) (err error) {
	var byteGroups = []int{8, 4, 4, 4, 12}

	if len(t) == CanonicalSize {
		if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
			return fmt.Errorf("Incorrect UUID format %s", t)
		}

		var src = t[:]
		var dst = u[:]

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
	} else {
		return fmt.Errorf("Incorrect UUID length: %s", t)
	}

	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
//
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != Size {
		return fmt.Errorf("UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}

	copy(u[:], data)

	return nil
}

// Implements the json.Unmarshaler interface.
//
// Unmarshal JSON value from canonical format
func (u *UUID) UnmarshalJSON(body []byte) error {
	var value string
	var err error = json.Unmarshal(body, &value)

	if err != nil {
		return err
	}

	return u.UnmarshalText([]byte(value))
}

// Implements the json.Marshaler interface.
//
// Marshal JSON value as string using the canonical format
func (u *UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// Value implements the driver.Valuer interface.
//
// Exports the value as a canonical string that most SQL servers supports.
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan implements the sql.Scanner interface with a solution specific for MSSQL.
//
// A 16-byte slice is handled by UnmarshalBinary, while a longer byte slice or a string is handled by UnmarshalText.
func (u *UUID) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		if len(src) == Size {
			// Assume that is MSSQL uniqueidentifier
			return u.SetFromMSSQLBytes(src)
		} else {
			// Canonical format
			return u.UnmarshalText(src)
		}
	case string:
		// Canonical format
		return u.UnmarshalText([]byte(src))
	}

	return fmt.Errorf("cannot convert %T to UUID", src)
}

// Set UUID from data fetched from a MSSQL Server as unique identifier.
//
// SQL server stores UUID in a different way some conversion is required.
func (u *UUID) SetFromMSSQLBytes(data []byte) error {
	var a = binary.LittleEndian.Uint32(data[0:])
	var b = binary.LittleEndian.Uint16(data[4:])
	var c = binary.LittleEndian.Uint16(data[6:])

	var d = binary.BigEndian.Uint16(data[8:])
	var e = binary.BigEndian.Uint16(data[10:])
	var f = binary.BigEndian.Uint32(data[12:])

	var uid = fmt.Sprintf("%08x-%04x-%04x-%04x-%04x%08x", a, b, c, d, e, f)

	return u.UnmarshalText([]byte(uid))
}
