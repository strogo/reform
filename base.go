package reform

import (
	"database/sql"
	"errors"
)

var (
	// ErrNoRows is returned from various methods when query produced no rows.
	ErrNoRows = sql.ErrNoRows

	// ErrNoPK is returned from various methods when primary key is required and not set.
	ErrNoPK = errors.New("reform: no primary key")
)

// View represents SQL database view or table.
type View interface {
	// Schema returns a schema name in SQL database.
	Schema() string

	// Name returns a view or table name in SQL database.
	Name() string

	// Columns returns a new slice of column names for that view or table in SQL database.
	Columns() []string

	// NewStruct makes a new struct for that view or table.
	NewStruct() Struct
}

// Table represents SQL database table with single-column primary key.
// It extends View.
type Table interface {
	View

	// NewRecord makes a new record for that table.
	NewRecord() Record

	// PKColumnIndex returns an index of primary key column for that table in SQL database.
	PKColumnIndex() uint
}

// Struct represents a row in SQL database view or table.
type Struct interface {
	// String returns a string representation of this struct or record.
	String() string

	// Values returns a slice of struct or record field values.
	// Returned interface{} values are never untyped nils.
	Values() []interface{}

	// Pointers returns a slice of pointers to struct or record fields.
	// Returned interface{} values are never untyped nils.
	Pointers() []interface{}

	// View returns View object for that struct.
	View() View
}

// Record represents a row in SQL database table with single-column primary key.
type Record interface {
	Struct

	// Table returns Table object for that record.
	Table() Table

	// PKValue returns a value of primary key for that record.
	// Returned interface{} value is never untyped nil.
	PKValue() interface{}

	// PKPointer returns a pointer to primary key field for that record.
	// Returned interface{} value is never untyped nil.
	PKPointer() interface{}

	// HasPK returns true if record has non-zero primary key set, false otherwise.
	HasPK() bool

	// SetPK sets record primary key.
	SetPK(pk interface{})
}

// BeforeInserter is an optional interface for Record which is used by Querier.Insert.
// It can be used to set record's timestamp fields, convert timezones, change data precision, etc.
// Returning error aborts operation.
type BeforeInserter interface {
	BeforeInsert() error
}

// BeforeUpdater is an optional interface for Record which is used by Querier.Update and Querier.UpdateColumns.
// It can be used to set record's timestamp fields, convert timezones, change data precision, etc.
// Returning error aborts operation.
type BeforeUpdater interface {
	BeforeUpdate() error
}

// AfterFinder is an optional interface for Record which is used by Querier's finders and selectors.
// It can be used to convert timezones, change data precision, etc.
// Returning error aborts operation.
type AfterFinder interface {
	AfterFind() error
}

// DBTX is an interface for database connection or transaction.
// It's implemented by *sql.DB, *sql.Tx, *DB, *TX and *Querier.
type DBTX interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
	QueryRow(query string, args ...interface{}) *sql.Row
}

// LastInsertIdMethod is a method of receiving primary key of last inserted row.
type LastInsertIdMethod int

const (
	// LastInsertId is method using sql.Result.LastInsertId().
	LastInsertId LastInsertIdMethod = iota

	// Returning is method using "RETURNING id" SQL syntax.
	Returning
)

// Dialect represents differences in various SQL dialects.
type Dialect interface {
	// Placeholder returns representation of placeholder parameter for given index,
	// typically "?" or "$1".
	Placeholder(index int) string

	// Placeholders returns representation of placeholder parameters for given start index and count,
	// typically []{"?", "?"} or []{"$1", "$2"}.
	Placeholders(start, count int) []string

	// QuoteIdentifier returns quoted database identifier,
	// typically "identifier" or `identifier`.
	QuoteIdentifier(identifier string) string

	// LastInsertIdMethod returns a method of receiving primary key of last inserted row.
	LastInsertIdMethod() LastInsertIdMethod
}

// check interface
var (
	_ DBTX = new(sql.DB)
	_ DBTX = new(sql.Tx)
)
