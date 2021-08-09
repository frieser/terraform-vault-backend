package vault

import (
	"github.com/bhoriuchi/terraform-backend-http/go/store"
	"github.com/bhoriuchi/terraform-backend-http/go/types"
	"github.com/mitchellh/mapstructure"
)

const LockPath = "/lock"
const statePath = "/state"

type Store struct {
	c *client
}

func NewStore() (store.Store, error) {
	var err error
	s := Store{}

	s.c, err = newClient()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s Store) Init() error {
	return nil
}

func (s Store) GetState(ref string) (state map[string]interface{}, encrypted bool, err error) {
	secret, err := s.c.read(ref + statePath)

	if err != nil {
		return nil, false, err
	}
	var stDoc types.StateDocument

	if secret == nil {
		return nil, false, store.ErrNotFound
	}
	err = mapstructure.Decode(secret.Data["data"], &stDoc)

	if err != nil {
		return nil, false, err
	}

	return stDoc.State, stDoc.Encrypted, nil
}

func (s Store) PutState(ref string, state, metadata map[string]interface{}, encrypted bool) error {
	var stDoc types.StateDocument

	stDoc.Ref = ref
	stDoc.State = state
	stDoc.Encrypted = encrypted
	stDoc.Metadata = metadata

	st := make(map[string]interface{})

	err := mapstructure.Decode(stDoc, &st)

	if err != nil {
		return err
	}
	_, err = s.c.write(ref+statePath, map[string]interface{}{
		"data": st,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteState(ref string) error {
	_, err := s.c.delete(ref + statePath)

	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetLock(ref string) (lock *types.Lock, err error) {
	secret, err := s.c.read(ref + LockPath)

	if err != nil {
		return nil, err
	}
	var l types.Lock

	if secret == nil {
		return nil, store.ErrNotFound
	}
	err = mapstructure.Decode(secret.Data["data"], &l)

	if err != nil {
		return nil, err
	}
	if l.ID == "" {
		return nil, store.ErrNotFound
	}

	return &l, nil
}

func (s Store) PutLock(ref string, lock types.Lock) error {
	lockData := make(map[string]interface{})
	err := mapstructure.Decode(lock, &lockData)

	if err != nil {
		return err
	}
	_, err = s.c.write(ref+LockPath,
		map[string]interface{}{"data": lockData})

	if err != nil {
		return err
	}

	return nil
}

func (s Store) DeleteLock(ref string) error {
	_, err := s.c.delete(ref + LockPath)

	if err != nil {
		return err
	}

	return nil
}
