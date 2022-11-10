package postgres

type Conn interface {
	Reader
	Writer
}
