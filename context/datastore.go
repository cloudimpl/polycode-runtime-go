package context

import "time"

type Collection interface {
	InsertOne(item interface{}) error
	InsertOneWithTTL(item interface{}, expireIn time.Duration) error
	UpdateOne(item interface{}) error
	UpdateOneWithTTL(item interface{}, expireIn time.Duration) error
	UpsertOne(item interface{}) error
	UpsertOneWithTTL(item interface{}, expireIn time.Duration) error
	DeleteOne(key string) error
	GetOne(key string, ret interface{}) (bool, error)
	Query() Query
}

type DataStore interface {
	Collection(name string) Collection
	GlobalCollection(name string) Collection
}

type DataStoreBuilder interface {
	WithTenantId(tenantId string) DataStoreBuilder
	WithPartitionKey(partitionKey string) DataStoreBuilder
	Get() DataStore
}
