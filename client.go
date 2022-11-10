package postgres

type Client interface {
	ClientReader
	ClientWriter
}
