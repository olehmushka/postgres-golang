package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/olehmushka/golang-toolkit/wrapped_error"
)

type connection struct {
	pgxConn           pgxConn
	batchItemsMaxSize int
}

func newConn(tx pgxConn, batchItemsMaxSize int) Conn {
	return &connection{
		pgxConn:           tx,
		batchItemsMaxSize: batchItemsMaxSize,
	}
}

func (c *connection) QueryRow(ctx context.Context, query string, args ...interface{}) (Row, error) {
	return c.pgxConn.QueryRow(ctx, query, args...), nil
}

func (c *connection) QueryRows(ctx context.Context, conn PgxConn, query string, args ...interface{}) (Rows, error) {
	return conn.Query(ctx, query, args...)
}

func (c *connection) CountRows(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	return count, c.pgxConn.QueryRow(ctx, query, args...).Scan(&count)
}

func (c *connection) InsertRow(ctx context.Context, query string, args ...interface{}) (uuid.UUID, error) {
	var id uuid.UUID
	return id, c.pgxConn.QueryRow(ctx, query, args...).Scan(&id)
}

func (c *connection) InsertRowWithStringID(ctx context.Context, query string, args ...interface{}) (string, error) {
	var id string
	return id, c.pgxConn.QueryRow(ctx, query, args...).Scan(&id)
}

func (c *connection) UpdateRowByID(ctx context.Context, sqlData map[string]interface{}, schemaName, tableName, id string) error {
	var (
		args []interface{}
		cols string
	)

	count := 0
	for key, value := range sqlData {
		count++
		cols += fmt.Sprintf("%s=$%d", key, count)

		if count != len(sqlData) {
			cols += ","
		}
		args = append(args, value)
	}

	if schemaName == "" {
		schemaName = DefaultSchemaName
	}

	// id is the last argument in query
	args = append(args, id)
	count++

	query := fmt.Sprintf("UPDATE %s.%s SET %s WHERE id=$%d", schemaName, tableName, cols, count)

	_, err := c.pgxConn.Exec(ctx, query, args...)

	return err
}

func (c *connection) ExecuteQuery(ctx context.Context, query string, args ...interface{}) error {
	if _, err := c.pgxConn.Exec(ctx, query, args...); err != nil {
		return wrapped_error.NewInternalServerError(err, "can not exec query")
	}
	return nil
}

func (c *connection) SendBatch(ctx context.Context, batch Batch) error {
	pgxBatch := &pgx.Batch{}

	num := 0
	for _, item := range batch.GetItems() {
		num++
		pgxBatch.Queue(item.GetQuery(), item.GetArgs()...)

		if num == batch.Len() {
			return c.sendBatch(ctx, pgxBatch)
		}

		if pgxBatch.Len() == c.batchItemsMaxSize {
			if err := c.sendBatch(ctx, pgxBatch); err != nil {
				return err
			}
			pgxBatch = &pgx.Batch{}
		}
	}

	return nil
}

func (c *connection) sendBatch(ctx context.Context, pgxBatch *pgx.Batch) error {
	batchResults := c.pgxConn.SendBatch(ctx, pgxBatch)

	var (
		queryErr error
		rows     pgx.Rows
	)

	for queryErr == nil {
		rows, queryErr = batchResults.Query()
		rows.Close()
	}
	return nil
}

func (c *connection) TruncateTables(ctx context.Context, tables []string) error {
	for _, table := range tables {
		if _, err := c.pgxConn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)); err != nil {
			return wrapped_error.NewInternalServerError(err, "can not exec query")
		}
	}

	return nil
}
