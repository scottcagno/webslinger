package sessions

import (
	"time"

	"github.com/scottcagno/webslinger/pkg/random"
)

const defaultInterval = 5 * time.Minute

type MemoryStore struct {
	ds *random.TimeoutMap[string, []byte]
}

func NewMemoryStore() *MemoryStore {
	return NewMemoryStoreWithInterval(defaultInterval)
}

func NewMemoryStoreWithInterval(interval time.Duration) *MemoryStore {
	return &MemoryStore{
		ds: random.NewTimeoutMap[string, []byte](interval),
	}
}

func (m *MemoryStore) Find(token string) ([]byte, error) {
	b, found := m.ds.Get(token)
	if !found {
		return nil, ErrSessionNotFound
	}
	return b, nil
}

func (m *MemoryStore) Save(token string, b []byte, expiry time.Time) error {
	m.ds.Put(token, b, time.Until(expiry))
	return nil
}

func (m *MemoryStore) Delete(token string) error {
	m.ds.Del(token)
	return nil
}
