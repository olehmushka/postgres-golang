package postgres

type BatchItem interface {
	SetQuery(query string)
	SetArgs(args ...any)
	GetQuery() string
	GetArgs() []any
}
