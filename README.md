# mssql-uuid
 - Go implementation of MS SQL Server (uniqueidentifier)[https://learn.microsoft.com/en-us/sql/t-sql/data-types/uniqueidentifier-transact-sql?view=sql-server-ver16] type.
 - Microsoft implementation of litle endian for the first half (the first 8 bytes), and big Endian encoding for the second 8 bytes.
 - Implements the JSON and SQL go interfaces to ensure that data is read correctly from DB.
 - The UNIQUEuniqueidentifierIDENTIFIER data type is a 16-byte GUID*.
 - This data type is used as primary key alternative to identity columns.
 - uniqueidentifier is globally unique, whereas identity is unique within a table.
 - Can be used alongside ORM such as (GORM)[] or (ent.io)[]

## License
 - The project is distributed under MIT license.
 - The lincese is available on the repository.
