package postgres

type ClientReader interface {
	Reader
	Pool
}
