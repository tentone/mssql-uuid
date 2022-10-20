# mssql-uuid
 - Go implementation of MS SQL Server [uniqueidentifier](https://learn.microsoft.com/en-us/sql/t-sql/data-types/uniqueidentifier-transact-sql?view=sql-server-ver16) type.
 - Implements the JSON [Marshal](https://pkg.go.dev/encoding/json#Marshaler)/[Unmarshal](https://pkg.go.dev/encoding/json#Unmarshaler) and SQL [Scan](https://pkg.go.dev/database/sql#Scanner) interfaces to ensure that data is read correctly from DB.
 - The uniqueidentifier data type is a 16-byte GUID*.
 - This data type is used as primary key alternative to identity columns.
 - uniqueidentifier is globally unique, whereas identity is unique within a table.
 - Can be used alongside ORM such as [GORM](https://gorm.io/)

## Example

```go
package main

import uuid "github.com/tentone/mssql-uuid"

func main() {
	var uid uuid.UUID = uuid.NewV4()
	print(uid.String())
}
```

## UUID Storage

 - Microsoft implementation of UUID uses litle endian for the first half the first 8 bytes, and big Endian encoding for the second 8 bytes.
 - Because of this other UUID libs fail to correctly parse the data.
```
LLLLLLLL-LLLL-LLLL-BBBB-BBBBBBBBBBBB
```

## License
 - The project is distributed under MIT license.
 - The lincese is available on the repository.
