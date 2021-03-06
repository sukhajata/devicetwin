package db

import "testing"

func TestPostgresEngine_ImplementsInterface(t *testing.T) {
	var _ SQLEngine = (*PostgresEngine)(nil)
}
