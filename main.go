package main

import (
	backend "github.com/bhoriuchi/terraform-backend-http/go"
	"github.com/frieser/terraform-vault-backend/vault"
	"log"
	"net/http"
	"os"
)

const (
	EnvBackendEncryptionKey = "BACKEND_ENCRYPTION_KEY"
	EnvBackendServerPort    = "BACKEND_SERVER_PORT"
)

const defaultBackendPort = "3000"
const pattern = "/backend"

func main() {
	encryptionKey := os.Getenv(EnvBackendEncryptionKey)

	if encryptionKey == "" {
		log.Fatalf("%s env var is mandatory",
			EnvBackendEncryptionKey)
	}
	port := os.Getenv(EnvBackendServerPort)

	if port == "" {
		port = defaultBackendPort
	}
	store, err := vault.NewStore()

	if err != nil {
		log.Fatal(err)
	}
	tfBackend := backend.NewBackend(store, &backend.Options{
		EncryptionKey: []byte(encryptionKey),
		Logger: func(level, message string, err error) {
			if err != nil {
				log.Printf("%s: %s - %v", level, message, err)

				return
			}
			log.Printf("%s: %s", level, message)
		},
		GetMetadataFunc: func(state map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"test": "metadata",
			}
		},
	})
	if err := tfBackend.Init(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "LOCK":
				tfBackend.HandleLockState(w, r)
			case "UNLOCK":
				tfBackend.HandleUnlockState(w, r)
			case http.MethodGet:
				tfBackend.HandleGetState(w, r)
			case http.MethodPost:
				tfBackend.HandleUpdateState(w, r)
			case http.MethodDelete:
				tfBackend.HandleDeleteState(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
	log.Printf("Starting test server on :%s \n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}