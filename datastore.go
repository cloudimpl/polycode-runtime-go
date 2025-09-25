package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
	"log"
	"time"
)

type ReadOnlyDataStoreBuilder struct {
	client    ServiceClient
	sessionId string
	tenantId  string
}

func (f ReadOnlyDataStoreBuilder) WithTenantId(tenantId string) polycode.ReadOnlyDataStoreBuilder {
	f.tenantId = tenantId
	return f
}

func (f ReadOnlyDataStoreBuilder) Get() polycode.ReadOnlyDataStore {
	fmt.Printf("getting unsafe db for tenant id = %s and partition key = %s", f.tenantId, f.partitionKey)
	return ReadOnlyDataStore{
		client:    f.client,
		sessionId: f.sessionId,
		tenantId:  f.tenantId,
	}
}

type ReadOnlyDataStore struct {
	client    ServiceClient
	sessionId string
	tenantId  string
}

func (r ReadOnlyDataStore) Collection(name string) polycode.ReadOnlyCollection {
	return Collection{
		client:    r.client,
		sessionId: r.sessionId,
		tenantId:  r.tenantId,
		name:      name,
	}
}

func (r ReadOnlyDataStore) GlobalCollection(name string) polycode.ReadOnlyCollection {
	return Collection{
		client:    r.client,
		sessionId: r.sessionId,
		tenantId:  r.tenantId,
		name:      name,
		isGlobal:  true,
	}
}

type DataStoreBuilder struct {
	client    ServiceClient
	sessionId string
	tenantId  string
}

func (f DataStoreBuilder) WithTenantId(tenantId string) polycode.DataStoreBuilder {
	f.tenantId = tenantId
	return f
}

func (f DataStoreBuilder) Get() polycode.DataStore {
	fmt.Printf("getting db for tenant id = %s", f.tenantId)
	return DataStore{
		client:    f.client,
		sessionId: f.sessionId,
		tenantId:  f.tenantId,
	}
}

type DataStore struct {
	client    ServiceClient
	sessionId string
	tenantId  string
}

func (d DataStore) Collection(name string) polycode.Collection {
	return Collection{
		client:    d.client,
		sessionId: d.sessionId,
		tenantId:  d.tenantId,
		name:      name,
	}
}

func (d DataStore) GlobalCollection(name string) polycode.Collection {
	return Collection{
		client:    d.client,
		sessionId: d.sessionId,
		tenantId:  d.tenantId,
		name:      name,
		isGlobal:  true,
	}
}

type ReadOnlyDoc struct {
	parent polycode.ReadOnlyDoc
	path   string
	id     string
	val    string
}

func (r ReadOnlyDoc) Unmarshal(item interface{}) error {
	return json.Unmarshal([]byte(d.val), item)
}

func (r ReadOnlyDoc) ExpireIn(expireIn time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r ReadOnlyDoc) Update(item interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r ReadOnlyDoc) Delete() error {
	//TODO implement me
	panic("implement me")
}

func (r ReadOnlyDoc) Collection(name string) polycode.ReadOnlyCollection {
	//TODO implement me
	panic("implement me")
}

type Doc struct {
	parent polycode.Doc
	path   string
	id     string
	val    string
}

func (d Doc) Unmarshal(item interface{}) error {
	return json.Unmarshal([]byte(d.val), item)
}

func (d Doc) ExpireIn(expireIn time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (d Doc) Update(item interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d Doc) Delete() error {
	//TODO implement me
	panic("implement me")
}

func (d Doc) Collection(name string) polycode.Collection {
	//TODO implement me
	panic("implement me")
}

type Collection struct {
	client    ServiceClient
	sessionId string
	name      string
	isGlobal  bool
	tenantId  string
}

func (c Collection) GetOne(id string) (polycode.Doc, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c Collection) InsertOne(id string, item interface{}) (polycode.Doc, error) {
	//TODO implement me
	panic("implement me")
}

func (c Collection) UpdateOne(id string, item interface{}) (polycode.Doc, error) {
	//TODO implement me
	panic("implement me")
}

func (c Collection) UpsertOne(id string, item interface{}) (polycode.Doc, error) {
	//TODO implement me
	panic("implement me")
}

func (c Collection) DeleteOne(id string) (polycode.Doc, error) {
	//TODO implement me
	panic("implement me")
}

func (c Collection) InsertOne(id string, item interface{}) error {
	return c.InsertOneWithTTL(item, -1)
}

func (c Collection) InsertOneWithTTL(id string, item interface{}, expireIn time.Duration) error {
	var ttl int64
	if expireIn == -1 {
		ttl = -1
	} else {
		ttl = time.Now().Unix() + int64(expireIn.Seconds())
	}

	req := PutRequest{
		TenantId:   c.tenantId,
		Action:     "insert",
		IsGlobal:   c.isGlobal,
		Collection: c.name,
		Key:        id,
		Item:       item,
		TTL:        ttl,
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

type ReadOnlyQuery struct {
	collection polycode.ReadOnlyCollection
	filter     string
	args       []any
	limit      int
}

func (q Query) Filter(expr string, args ...interface{}) polycode.Query {
	q.filter = expr
	q.args = args
	return q
}

func (q Query) Limit(limit int) polycode.Query {
	q.limit = limit
	return q
}

func (q Query) One(ctx context.Context, ret interface{}) (bool, error) {
	req := QueryRequest{
		TenantId:   q.collection.tenantId,
		IsGlobal:   q.collection.isGlobal,
		Collection: q.collection.name,
		Key:        "",
		Filter:     q.filter,
		Args:       q.args,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		fmt.Printf("client: error query item %s\n", err.Error())
		return false, err
	}

	if len(r) == 0 {
		return false, nil
	}

	e := r[0]
	err = ConvertType(e, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
	}

	return true, nil
}

func (q Query) All(ctx context.Context, ret interface{}) error {
	req := QueryRequest{
		TenantId:     q.collection.tenantId,
		PartitionKey: q.collection.partitionKey,
		IsGlobal:     q.collection.isGlobal,
		Collection:   q.collection.name,
		Key:          "",
		Filter:       q.filter,
		Args:         q.args,
		Limit:        q.limit,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		log.Println("client: error query item ", err.Error())
		return err
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return err
	}

	return nil
}

type Query struct {
	collection *Collection
	filter     string
	args       []any
	limit      int
}

func (q Query) Filter(expr string, args ...interface{}) polycode.Query {
	q.filter = expr
	q.args = args
	return q
}

func (q Query) Limit(limit int) polycode.Query {
	q.limit = limit
	return q
}

func (q Query) One(ctx context.Context, ret interface{}) (bool, error) {
	req := QueryRequest{
		TenantId:   q.collection.tenantId,
		IsGlobal:   q.collection.isGlobal,
		Collection: q.collection.name,
		Key:        "",
		Filter:     q.filter,
		Args:       q.args,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		fmt.Printf("client: error query item %s\n", err.Error())
		return false, err
	}

	if len(r) == 0 {
		return false, nil
	}

	e := r[0]
	err = ConvertType(e, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return false, err
	}

	return true, nil
}

func (q Query) All(ctx context.Context, ret interface{}) error {
	req := QueryRequest{
		TenantId:     q.collection.tenantId,
		PartitionKey: q.collection.partitionKey,
		IsGlobal:     q.collection.isGlobal,
		Collection:   q.collection.name,
		Key:          "",
		Filter:       q.filter,
		Args:         q.args,
		Limit:        q.limit,
	}

	r, err := q.collection.client.QueryItems(q.collection.sessionId, req)
	if err != nil {
		log.Println("client: error query item ", err.Error())
		return err
	}

	err = ConvertType(r, ret)
	if err != nil {
		fmt.Printf("failed to convert type: %s\n", err.Error())
		return err
	}

	return nil
}
