package db

import (
	"context"
	"encoding/json"
	"github.com/sukhajata/devicetwin.git/pkg/loggerhelper"
	pb "github.com/sukhajata/ppconnection"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TimescaleEngine implements SQLEngine
type PostgresEngine struct {
	pool        *pgxpool.Pool
	FailureChan chan error
}

// NewTimescaleEngine factory method for creating timescale engine
func NewTimescaleEngine(psqlURL string) (*PostgresEngine, error) {
	var err error
	var pool *pgxpool.Pool
	retries := 0
	for {
		pool, err = pgxpool.Connect(context.Background(), psqlURL)
		if err != nil {
			retries++
			loggerhelper.WriteToLog(err)
			if retries > 5 {
				return nil, err
			}
			time.Sleep(time.Second * 2)
			continue
		}

		break
	}

	return &PostgresEngine{
		pool:        pool,
		FailureChan: make(chan error, 1),
	}, nil

}

// Query - get array
func (t *PostgresEngine) Query(queryString string, arguments ...interface{}) ([]interface{}, error) {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), queryString, arguments...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]interface{}, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return results, err
		}
		results = append(results, values)
	}

	return results, err
}

// QueryConnections - get array of Connection structs
func (t *PostgresEngine) QueryConnections(queryString string, arguments ...interface{}) ([]*pb.Connection, error) {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), queryString, arguments...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	connections := make([]*pb.Connection, 0)
	for rows.Next() {
		var data string
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		var connection *pb.Connection
		err := json.Unmarshal([]byte(data), &connection)
		if err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}

	return connections, nil
}

// Exec - run a query without return
func (t *PostgresEngine) Exec(queryString string, arguments ...interface{}) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer conn.Release()
	_, err = conn.Exec(context.Background(), queryString, arguments...)

	return err
}

// ScanRow - query a row and scan into the value pointer
func (t *PostgresEngine) ScanRow(queryString string, valuePtr interface{}, arguments ...interface{}) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return err
	}

	defer conn.Release()
	err = conn.QueryRow(context.Background(), queryString, arguments...).Scan(valuePtr)

	return err
}

// Close the pool
func (t *PostgresEngine) Close() {
	t.pool.Close()
}
