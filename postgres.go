package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olehmushka/golang-toolkit/wrapped_error"
)

type client struct {
	pool              *pgxpool.Pool
	batchItemsMaxSize int
}

func New(cfg *Config) (Client, error) {
	return newClient(context.Background(), cfg)
}

func NewWriter(cfg *Config) (ClientWriter, error) {
	return newClient(context.Background(), cfg)
}

func NewReader(cfg *Config) (ClientReader, error) {
	return newClient(context.Background(), cfg)
}

func newClient(ctx context.Context, cfg *Config) (*client, error) {
	pgClient := &client{}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, wrapped_error.NewInternalServerError(err, "pgx pool connection failed")
	}

	pgClient.pool = pool

	batchItemsMaxSize := cfg.BatchItemsMaxSize
	if batchItemsMaxSize == 0 {
		batchItemsMaxSize = DefaultBatchItemsMaxSize
	}

	pgClient.batchItemsMaxSize = batchItemsMaxSize

	return pgClient, nil
}

func (c *client) ClosePool() {
	c.pool.Close()
}

// AcquireConn acquires a PgxConnection. Do not forget to use method Release!
func (c *client) AcquireConn(ctx context.Context) (PgxConn, error) {
	return c.acquireConn(ctx)
}

func (c *client) acquireConn(ctx context.Context) (PgxConn, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, wrapped_error.NewInternalServerError(err, "Acquire failed")
	}

	return conn, nil
}

func (c *client) QueryRow(ctx context.Context, query string, args ...any) (Row, error) {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).QueryRow(ctx, query, args...)
}

func (c *client) QueryRows(ctx context.Context, conn PgxConn, query string, args ...any) (Rows, error) {
	return newConn(conn, c.batchItemsMaxSize).QueryRows(ctx, conn, query, args...)
}

func (c *client) CountRows(ctx context.Context, query string, args ...any) (int, error) {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).CountRows(ctx, query, args...)
}

func (c *client) InsertRow(ctx context.Context, query string, args ...any) (uuid.UUID, error) {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).InsertRow(ctx, query, args...)
}

func (c *client) InsertRowWithStringID(ctx context.Context, query string, args ...any) (string, error) {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).InsertRowWithStringID(ctx, query, args...)
}

func (c *client) UpdateRowByID(ctx context.Context, sqlData map[string]any, schemaName, tableName, id string) error {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).UpdateRowByID(ctx, sqlData, schemaName, tableName, id)
}

func (c *client) ExecuteQuery(ctx context.Context, query string, args ...any) error {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).ExecuteQuery(ctx, query, args...)
}

func (c *client) SendBatch(ctx context.Context, batch Batch) error {
	tx, err := c.BeginTx(ctx, TxOptions{})
	if err != nil {
		return err
	}

	if err = tx.SendBatch(ctx, batch); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return tx.Commit(ctx)
}

func (c *client) BeginTx(ctx context.Context, opts TxOptions) (Transaction, error) {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.TxIsoLevel(opts.IsoLevel),
		AccessMode:     pgx.TxAccessMode(opts.AccessMode),
		DeferrableMode: pgx.TxDeferrableMode(opts.DeferrableMode),
	})

	if err != nil {
		conn.Release()
		return nil, wrapped_error.NewInternalServerError(err, "start transaction failed")
	}

	return newTransaction(tx, conn), nil
}

func (c *client) TruncateTables(ctx context.Context, tables []string) error {
	conn, err := c.acquireConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return newConn(conn, c.batchItemsMaxSize).TruncateTables(ctx, tables)
}
