package frontend

import (
	"context"
	"fmt"
	"key-value-store/core"
	kvs "key-value-store/frontend/keyvalue"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCFrontEnd struct {
	store *core.KeyValueStore
	kvs.UnimplementedKeyValueServer
}

func (gf *GRPCFrontEnd) Start(store *core.KeyValueStore) error {
	gf.store = store

	s := grpc.NewServer()
	kvs.RegisterKeyValueServer(s, &GRPCFrontEnd{})

	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Println("listening on port 50051")
	
	return s.Serve(lis)
}

func (gf *GRPCFrontEnd) Get(ctx context.Context, r *kvs.GetRequest) (*kvs.GetResponse, error) {
	log.Printf("received get key=%v\n", r.Key)

	val, err := gf.store.Get(r.Key)

	return &kvs.GetResponse{Value: val}, err
}

func (gf *GRPCFrontEnd) Put(ctx context.Context, r *kvs.PutRequest) (*kvs.PutResponse, error) {
	log.Printf("received put key=%v value=%v\n", r.Key, r.Value)

	err := gf.store.Put(r.Key, r.Value)

	return &kvs.PutResponse{}, err
}

func (gf *GRPCFrontEnd) Delete(ctx context.Context, r *kvs.DeleteRequest) (*kvs.DeleteResponse, error) {
	log.Printf("received delete key=%v\n", r.Key)

	err := gf.store.Delete(r.Key)

	return &kvs.DeleteResponse{}, err
}
