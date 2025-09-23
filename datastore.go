package runtime

import (
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
	"time"
)

type ReadOnlyDataStoreBuilder struct {
	client       ServiceClient
	sessionId    string
	tenantId     string
	partitionKey string
}

func (f ReadOnlyDataStoreBuilder) WithTenantId(tenantId string) polycode.ReadOnlyDataStoreBuilder {
	f.tenantId = tenantId
	return f
}

func (f ReadOnlyDataStoreBuilder) WithPartitionKey(partitionKey string) polycode.ReadOnlyDataStoreBuilder {
	f.partitionKey = partitionKey
	return f
}

func (f ReadOnlyDataStoreBuilder) Get() polycode.ReadOnlyDataStore {
	fmt.Printf("getting unsafe db for tenant id = %s and partition key = %s", f.tenantId, f.partitionKey)
	return ReadOnlyDataStore{
		client:       f.client,
		sessionId:    f.sessionId,
		tenantId:     f.tenantId,
		partitionKey: f.partitionKey,
	}
}

type ReadOnlyDataStore struct {
	client       ServiceClient
	sessionId    string
	tenantId     string
	partitionKey string
}

func (r ReadOnlyDataStore) Collection(name string) polycode.ReadOnlyCollection {
	return Collection{
		client:       r.client,
		sessionId:    r.sessionId,
		tenantId:     r.tenantId,
		partitionKey: r.partitionKey,
		name:         name,
	}
}

func (r ReadOnlyDataStore) GlobalCollection(name string) polycode.ReadOnlyCollection {
	return Collection{
		client:       r.client,
		sessionId:    r.sessionId,
		tenantId:     r.tenantId,
		partitionKey: r.partitionKey,
		name:         name,
		isGlobal:     true,
	}
}

type DataStoreBuilder struct {
	client       ServiceClient
	sessionId    string
	tenantId     string
	partitionKey string
}

func (f DataStoreBuilder) WithTenantId(tenantId string) polycode.DataStoreBuilder {
	f.tenantId = tenantId
	return f
}

func (f DataStoreBuilder) WithPartitionKey(partitionKey string) polycode.DataStoreBuilder {
	f.partitionKey = partitionKey
	return f
}

func (f DataStoreBuilder) Get() polycode.DataStore {
	fmt.Printf("getting unsafe db for tenant id = %s and partition key = %s", f.tenantId, f.partitionKey)
	return DataStore{
		client:       f.client,
		sessionId:    f.sessionId,
		tenantId:     f.tenantId,
		partitionKey: f.partitionKey,
	}
}

type DataStore struct {
	client       ServiceClient
	sessionId    string
	tenantId     string
	partitionKey string
}

func (d DataStore) Collection(name string) polycode.Collection {
	return Collection{
		client:       d.client,
		sessionId:    d.sessionId,
		tenantId:     d.tenantId,
		partitionKey: d.partitionKey,
		name:         name,
	}
}

func (d DataStore) GlobalCollection(name string) polycode.Collection {
	return Collection{
		client:       d.client,
		sessionId:    d.sessionId,
		tenantId:     d.tenantId,
		partitionKey: d.partitionKey,
		name:         name,
		isGlobal:     true,
	}
}

type Collection struct {
	client       ServiceClient
	sessionId    string
	name         string
	isGlobal     bool
	tenantId     string
	partitionKey string
}

func (c Collection) InsertOne(item interface{}) error {
	return c.InsertOneWithTTL(item, -1)
}

func (c Collection) InsertOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		TenantId:     c.tenantId,
		PartitionKey: c.partitionKey,
		Action:       "insert",
		IsGlobal:     c.isGlobal,
		Collection:   c.name,
		Key:          id,
		Item:         item,
		TTL:          ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) UpdateOne(item interface{}) error {
	return c.UpdateOneWithTTL(item, -1)
}

func (c Collection) UpdateOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		TenantId:     c.tenantId,
		PartitionKey: c.partitionKey,
		Action:       "update",
		IsGlobal:     c.isGlobal,
		Collection:   c.name,
		Key:          id,
		Item:         item,
		TTL:          ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) UpsertOne(item interface{}) error {
	return c.UpsertOneWithTTL(item, -1)
}

func (c Collection) UpsertOneWithTTL(item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

	id, err := GetId(item)
	if err != nil {
		fmt.Printf("failed to get id: %s\n", err.Error())
		return err
	}

	req := PutRequest{
		TenantId:     c.tenantId,
		PartitionKey: c.partitionKey,
		Action:       "upsert",
		IsGlobal:     c.isGlobal,
		Collection:   c.name,
		Key:          id,
		Item:         item,
		TTL:          ttl,
	}

	err = c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) DeleteOne(key string) error {
	req := PutRequest{
		TenantId:     c.tenantId,
		PartitionKey: c.partitionKey,
		Action:       "delete",
		IsGlobal:     c.isGlobal,
		Collection:   c.name,
		Key:          key,
	}

	err := c.client.PutItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to put item: %s\n", err.Error())
		return err
	}

	return nil
}

func (c Collection) GetOne(key string, ret interface{}) (bool, error) {
	req := QueryRequest{
		TenantId:     c.tenantId,
		PartitionKey: c.partitionKey,
		IsGlobal:     c.isGlobal,
		Collection:   c.name,
		Key:          key,
		Filter:       "",
		Args:         nil,
	}

	r, err := c.client.GetItem(c.sessionId, req)
	if err != nil {
		fmt.Printf("failed to get item: %s\n", err.Error())
		return false, err
	}

	if r == nil {
		println("item not found")
		return false, nil
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
	}

	return true, nil
}

func (c Collection) Query() polycode.Query {
	return Query{
		collection: &c,
	}
}
