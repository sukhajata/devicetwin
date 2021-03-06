package db

// NoSQLEngine represents a nosql db engine
type NoSQLEngine interface {
	Query(bucketName string, queryString string, arguments []interface{}) ([]interface{}, error)
	Update(bucketName string, key string, path string, value interface{}) error
	Lookup(bucketName string, key string, path string, valuePtr interface{}) error
	Get(bucketName string, key string, valuePtr interface{}) error
	ArrayAppend(bucketName string, key string, path string, value interface{}) error
	Upsert(bucketName string, key string, value interface{}) error
	Delete(bucketName string, key string) error
}
