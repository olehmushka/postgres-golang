package postgres

type Config struct {
	Username          string
	Password          string
	DBName            string
	Host              string
	Port              int
	BatchItemsMaxSize int
}
