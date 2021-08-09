package vault

import (
	"github.com/hashicorp/vault/api"
	"log"
	"os"
	"sync"
	"time"
)

type client struct {
	vaultCli *api.Client

	roleID          string
	secretID        string
	tokenExpiration time.Time
	m               *sync.Mutex
}

func newClient() (*client, error) {
	config := api.DefaultConfig()
	err := config.ReadEnvironment()

	if err != nil {
		return nil, err
	}
	c := &client{
		m: &sync.Mutex{},
		roleID: os.Getenv("VAULT_ROLE_ID"),
		secretID: os.Getenv("VAULT_SECRET_ID"),
	}
	c.vaultCli, err = api.NewClient(config)

	if err != nil {
		return nil, err
	}
	err = c.authenticate()



	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) read(path string) (*api.Secret, error) {
	err := c.refreshToken()

	if err != nil {
		return nil, err
	}

	return c.vaultCli.Logical().Read(path)
}

func (c *client) write(path string, data map[string]interface{}) (*api.Secret, error) {
	err := c.refreshToken()

	if err != nil {
		return nil, err
	}

	return c.vaultCli.Logical().Write(path, data)
}

func (c *client) delete(path string) (*api.Secret, error) {
	err := c.refreshToken()

	if err != nil {
		return nil, err
	}

	return c.vaultCli.Logical().Delete(path)
}

func (c *client) authenticate() (err error) {
	if c.roleID == "" && c.secretID == "" {

		return nil
	}

	options := map[string]interface{}{
		"role_id":   c.roleID,
		"secret_id": c.secretID,
	}
	var secret *api.Secret

	if secret, err = c.vaultCli.Logical().Write("auth/approle/login", options); err != nil {
		return err
	}
	c.vaultCli.SetToken(secret.Auth.ClientToken)
	c.tokenExpiration = time.Now()

	if secret.Auth.Renewable {
		c.tokenExpiration = c.tokenExpiration.Add(time.Duration(secret.Auth.LeaseDuration-60) * time.Second)
	}

	return nil
}

func (c *client) refreshToken() error {
	if c.roleID == "" && c.secretID == "" {

		return nil
	}

	c.m.Lock()
	if c.tokenExpiration.Before(time.Now()) {
		log.Println("Refreshing Vault token...")

		if err := c.authenticate(); err != nil {
			c.m.Unlock()

			return err
		}
	}
	c.m.Unlock()

	return nil
}
