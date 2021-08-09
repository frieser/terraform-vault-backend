// docker run -it -p 8200:8200 vault
// annotate the token

// export VAULT_ADDR = "http://localhost:8200"
// export VAULT_TOKEN = paste_the_previous_token
// export BACKEND_ENCRYPTION_KEY="myscret"

// go run main.go

terraform {
  backend "http" {
    address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
    lock_address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
    unlock_address = "http://localhost:3000/backend?ref=secret/data/test&encrypt=true"
  }
}

resource "local_file" "testfile" {
  content = "foobar"
  filename = "${path.module}/test.json"
}