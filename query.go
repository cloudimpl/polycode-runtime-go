package runtime

import (
	"context"
	"fmt"
	"github.com/cloudimpl/polycode-sdk-go"
	"log"
)

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
		TenantId:     q.collection.tenantId,
		PartitionKey: q.collection.partitionKey,
		IsGlobal:     q.collection.isGlobal,
		Collection:   q.collection.name,
		Key:          "",
		Filter:       q.filter,
		Args:         q.args,
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
