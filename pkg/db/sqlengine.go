package db

import pb "github.com/sukhajata/ppconnection"

// SQLEngine represents a sql db engine
type SQLEngine interface {
	Exec(sql string, arguments ...interface{}) error
	Query(sql string, arguments ...interface{}) ([]interface{}, error)
	QueryConnections(queryString string, arguments ...interface{}) ([]*pb.Connection, error)
	ScanRow(sql string, valuePtr interface{}, arguments ...interface{}) error
	Close()
}
