package main

import (
	"fmt"
	kvs "key-value-store/internal/keyvalue"
	"key-value-store/internal/transaction"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

const (
	webPort = "8080"
	grpcPort = "50051"
)

var (
	certPath string
	keyPath  string
)

func main() {
	err := initializeTransactionLog(transaction.FileTransaction)
	if err != nil {
		log.Fatalf("error initializing transaction log: %s\n", err)
	}

	// err = initializeCertificate()
	// if err != nil {
	// 	log.Fatalf("failed to initialize certificates: %s", err)
	// }

	// log.Printf("listening on port %s", webPort)
	// log.Fatalln(http.ListenAndServeTLS(fmt.Sprintf(":%s", webPort), certPath, keyPath, Router()))

	s := grpc.NewServer()
	kvs.RegisterKeyValueServer(s, &server{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	log.Printf("listening on port %s\n", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initializeCertificate() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	execDir := filepath.Dir(execPath)

	certDir := filepath.Join(execDir, "cert")

	certPath = filepath.Join(certDir, "cert.pem")
	keyPath = filepath.Join(certDir, "key.pem")

	return nil
}
