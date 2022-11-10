package postgres

type Rows interface {
	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() error

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	// Scan reads the values from the current row into dest values positionally.
	// dest can include pointers to core types, values implementing the Scanner
	// interface, and nil. nil will skip the value entirely.
	Scan(dest ...any) error

	// Values returns the decoded row values.
	Values() ([]any, error)

	// RawValues returns the unparsed bytes of the row values. The returned [][]byte is only valid until the Next
	// call or the Rows is closed. However, the underlying byte data is safe to retain a reference to and mutate.
	RawValues() [][]byte
}
