package main

import (
	"key-value-store/core"
	"key-value-store/frontend"
	"key-value-store/transact"
	"log"
	"os"
)

func main() {
	tl, err := transact.NewTransactionLogger(os.Getenv("TLOG_TYPE"))
	if err != nil {
		log.Fatalln(err)
	}

	store := core.NewKeyValueStore(tl)

	fe, err := frontend.NewFrontEnd(os.Getenv("FRONTEND_TYPE"))
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(fe.Start(store))
}