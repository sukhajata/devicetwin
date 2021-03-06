package db

import "testing"

func TestCouchbaseEngine_ImplementsInterface(t *testing.T) {
	var _ NoSQLEngine = (*CouchbaseEngine)(nil)
}
