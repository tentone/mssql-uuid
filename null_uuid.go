package mssql

import (
	"database/sql/driver"
	"encoding/json"
)

// NullUUID can be used with the standard sql package to represent a UUID value that can be NULL in the database.
//
// Contains a Valid parameter using to check if the UUID value is null or not.
type NullUUID struct {
	UUID

	Valid bool
}

// Value implements the driver.Valuer interface.
//
// Used to provide a value to the SQL server for storage.
func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	// Delegate to UUID Value function
	return u.UUID.Value()
}

// Scan implements the sql.Scanner interface.
//
// Used to read the value provided by the SQL server.
func (u *NullUUID) Scan(src any) error {
	if src == nil {
		u.UUID, u.Valid = NilUUID, false
		return nil
	}

	// Delegate to UUID Scan function
	u.Valid = true
	return u.UUID.Scan(src)
}

// Implements the json.Unmarshaler interface.
//
// Unmarshal JSON value from canonical format
func (u *NullUUID) UnmarshalJSON(body []byte) error {
	var value string
	var err error = json.Unmarshal(body, &value)
	if err != nil {
		return err
	}

	if len(value) < CanonicalSize {
		u.Valid = false
		return nil
	}

	u.Valid = true
	return u.UUID.UnmarshalText([]byte(value))
}

// Implements the json.Marshaler interface.
//
// Marshal JSON value as string using the canonical format
func (u *NullUUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(u.UUID.String())
}
