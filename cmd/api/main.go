package main

import (
	"fmt"
	"key-value-store/internal/transaction"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const webPort = "8080"

var (
	certPath string
	keyPath  string
)

func main() {
	err := initializeTransactionLog(transaction.FileTransaction)
	if err != nil {
		log.Fatalf("error initializing transaction log: %s\n", err)
	}

	err = initializeCertificate()
	if err != nil {
		log.Fatalf("failed to initialize certificates: %s", err)
	}

	log.Printf("listening on port %s", webPort)
	log.Fatalln(http.ListenAndServeTLS(fmt.Sprintf(":%s", webPort), certPath, keyPath, Router()))
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
