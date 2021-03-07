package db

import (
	"github.com/sukhajata/devicetwin/pkg/errorhelper"
	"gopkg.in/couchbase/gocb.v1"
)

// CouchbaseEngine implements NoSQLEngine
type CouchbaseEngine struct {
	customerBucket *gocb.Bucket
	sharedBucket   *gocb.Bucket

	// BucketNameShared - the name of the shared bucket
	BucketNameShared string

	// BucketName - the name of the customer bucket
	BucketName string
}

// NewCouchbaseEngine factory method for creating couchbase engine
func NewCouchbaseEngine(serverAddress string, username string, password string, bucketName string, bucketNameShared string) (*CouchbaseEngine, error) {
	failCount := 0
	var cc *CouchbaseEngine

	for {
		cluster, err := gocb.Connect("couchbase://" + serverAddress)
		if err != nil {
			failCount++
			if failCount > 5 {
				return nil, err
			}
			errorhelper.StartUpError(err)
			continue
		}

		err = cluster.Authenticate(gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		})
		if err != nil {
			failCount++
			if failCount > 5 {
				return nil, err
			}
			errorhelper.StartUpError(err)
			continue
		}

		bucket, err := cluster.OpenBucket(bucketName, "")
		if err != nil {
			errorhelper.StartUpError(err)
			continue
		}

		bucketShared, err := cluster.OpenBucket(bucketNameShared, "")
		if err != nil {
			failCount++
			if failCount > 5 {
				return nil, err
			}
			errorhelper.StartUpError(err)
			continue
		}

		cc = &CouchbaseEngine{
			customerBucket:   bucket,
			sharedBucket:     bucketShared,
			BucketName:       bucketName,
			BucketNameShared: bucketNameShared,
		}

		break
	}

	return cc, nil

}

// Query - run a query
func (c *CouchbaseEngine) Query(bucketName string, queryString string, arguments []interface{}) ([]interface{}, error) {
	query := gocb.NewN1qlQuery(queryString)
	var rows gocb.QueryResults
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}

	rows, err := bucket.ExecuteN1qlQuery(query, arguments)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)

	for {
		var row interface{}
		if rows.Next(&row) {
			results = append(results, row)
		} else {
			break
		}
	}

	return results, nil
}

// Update - update a field
func (c *CouchbaseEngine) Update(bucketName string, key string, path string, value interface{}) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	_, err := bucket.
		MutateIn(key, 0, 0).
		Upsert(path, value, false).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

// Lookup - get a subdocument
func (c *CouchbaseEngine) Lookup(bucketName string, key string, path string, valuePtr interface{}) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	frag, err := bucket.LookupIn(key).Get(path).Execute()
	if err != nil {
		return err
	}
	err = frag.Content(path, &valuePtr)
	return err
}

// Get - get a document by key
func (c *CouchbaseEngine) Get(bucketName string, key string, valuePtr interface{}) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	_, err := bucket.Get(key, &valuePtr)
	return err
}

// Upsert - insert a new doc, overwriting any existing doc with same key
func (c *CouchbaseEngine) Upsert(bucketName string, key string, value interface{}) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	_, err := bucket.Upsert(key, value, 0)
	return err
}

// ArrayAppend - add an item to an array in a document
func (c *CouchbaseEngine) ArrayAppend(bucketName string, key string, path string, value interface{}) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	_, err := bucket.
		MutateIn(key, 0, 0).
		ArrayAppend(path, value, true).
		Execute()
	return err
}

// Delete - delete a document
func (c *CouchbaseEngine) Delete(bucketName string, key string) error {
	bucket := c.customerBucket
	if bucketName == c.BucketNameShared {
		bucket = c.sharedBucket
	}
	_, err := bucket.Remove(key, 0)
	return err
}
