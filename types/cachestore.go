package types

type Cache struct {
	*Mapping
	store Store
}

func NewCache(store Store) *Cache {
	c := &Cache{
		store: store,
	}
	c.Reset()
	return c
}

func (c *Cache) Reset() {
	c.Mapping = NewMapping()
}

func (c *Cache) Set(key []byte, value []byte) {
	c.Mapping.Set(key, value)
}

func (c *Cache) Get(key []byte) (value []byte) {
	value = c.Mapping.Get(key)
	if value != nil {
		return
	}
	value = c.store.Get(key)
	c.Set(key, value)
	return
}

func (c *Cache) Sync() {
	for e := c.List.head; e != nil; e = e.next {
		c.store.Set(e.key, e.value)
	}
	c.Reset()
}

type Store interface {
	Set(key, value []byte)
	Get(key []byte) (value []byte)
}

type MemStore map[string][]byte

func NewMemStore() MemStore {
	mstore := make(map[string][]byte)
	return mstore
}

func (mstore MemStore) Set(key, value []byte) {
	mstore[string(key)] = value
}

func (mstore MemStore) Get(key []byte) []byte {
	return mstore[string(key)]
}
