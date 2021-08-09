Terraform Vault Backend
=======================

Terraform HTTP backend that stores the state in a Vault secret.

* Locking
* Encryption
* Versioning
* Authentication with approle

Usage
-----
Look for the example/ folder more a more detailed use.

```hcl
terraform {
  backend "http" {
    address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
    lock_address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
    unlock_address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
  }
}
```
`ref` is the path where the state and the lock are going to be stored.

`encrypt` indicates whether the encryption will be enabled or not.

Configure the backend
---------------------

The backend reads the following environment variables to set behavioral defaults.

`VAULT_ADDR`

Address of the Vault server expressed as a URL and port, for example:
https://127.0.0.1:8200/.

DEFAULT: "http://127.0.0.1:8200"

`VAULT_TOKEN`

Vault authentication token. Conceptually similar to a session token on a website, 
the VAULT_TOKEN environment variable holds the contents of the token.

MANDATORY IF VAULT_ROLE_ID and VAULT_SECRET_ID are empty

`VAULT_ROLE_ID`

Role id used for approle authentication flow.

MANDATORY IF VAULT_TOKEN is empty

`VAULT_SECRET_ID`

Secret associated to a role for approle authentication flow.

MANDATORY IF VAULT_TOKEN is empty

`BACKEND_SERVER_PORT`

Port where the backend server will be listening on.

DEFAULT: "3000"

`BACKEND_ENCRYPTION_KEY`

The encryption key used to encrypt the communication between the backend
server and the Vault server.

MANDATORY if encryption is enabled
