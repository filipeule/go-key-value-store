package main

import (
	"context"
	kvs "key-value-store/internal/keyvalue"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cli := kvs.NewKeyValueClient(conn)

	var action, key, value string

	if len(os.Args) > 2 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch action {
	case "get":
		r, err := cli.Get(ctx, &kvs.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}
		log.Printf("get %s returns %s", key, r.Value)
	case "put":
		_, err := cli.Put(ctx, &kvs.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("could not put key %s: %v\n", key, err)
		}
		log.Printf("put %s", key)
	case "delete":
		_, err := cli.Delete(ctx, &kvs.DeleteRequest{Key: key})
		if err != nil {
			log.Fatalf("could not delete key %s: %v\n", key, err)
		}
		log.Printf("delete %s", key)
	}
}