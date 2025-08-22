package main

import (
	"context"
	kvs "key-value-store/internal/keyvalue"
	"log"
)

type server struct {
	kvs.UnimplementedKeyValueServer
}

func (s *server) Get(ctx context.Context, r *kvs.GetRequest) (*kvs.GetResponse, error) {
	log.Printf("received get key=%v\n", r.Key)

	val, err := Get(r.Key)

	return &kvs.GetResponse{Value: val}, err
}

func (s *server) Put(ctx context.Context, r *kvs.PutRequest) (*kvs.PutResponse, error) {
	log.Printf("received put key=%v value=%v\n", r.Key, r.Value)

	err := Put(r.Key, r.Value)

	return &kvs.PutResponse{}, err
}

func (s *server) Delete(ctx context.Context, r *kvs.DeleteRequest) (*kvs.DeleteResponse, error) {
	log.Printf("received delete key=%v\n", r.Key)

	err := Delete(r.Key)

	return &kvs.DeleteResponse{}, err
}