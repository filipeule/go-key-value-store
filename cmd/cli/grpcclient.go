package main

import (
	"context"
	kvs "key-value-store/internal/keyvalue"

	"google.golang.org/grpc"
)

type KeyValueStoreClient interface {
	Get(ctx context.Context, in *kvs.GetRequest, opts ...grpc.CallOption) (*kvs.GetResponse, error)
	Put(ctx context.Context, in *kvs.PutRequest, opts ...grpc.CallOption) (*kvs.PutResponse, error)
	Delete(ctx context.Context, in *kvs.DeleteRequest, opts ...grpc.CallOption) (*kvs.DeleteResponse, error)
}