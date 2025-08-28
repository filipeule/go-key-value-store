package frontend

import (
	"errors"
	"fmt"
	"io"
	"key-value-store/core"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var (
	certPath string
	keyPath  string
)

type RestFrontEnd struct {
	store *core.KeyValueStore
}

func (rf *RestFrontEnd) Start(store *core.KeyValueStore) error {
	rf.store = store

	r := mux.NewRouter()

	r.HandleFunc("/v1/key/{key}", rf.keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/key/{key}", rf.keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/key/{key}", rf.keyValueDeleteHandler).Methods("DELETE")

	err := initializeCertificate()
	if err != nil {
		return fmt.Errorf("failed to initialize certificates: %v", err)
	}

	return http.ListenAndServeTLS(":8080", certPath, keyPath, r)
}

func (rf *RestFrontEnd) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rf.store.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rf *RestFrontEnd) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	val, err := rf.store.Get(key)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func (rf *RestFrontEnd) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := rf.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
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

