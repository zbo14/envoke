package state

import (
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	merk "github.com/tendermint/merkleeyes/client"
	tmsp "github.com/tendermint/tmsp/types"
	"github.com/zballs/go_resonate/types"
)

type State struct {
	cache *types.Cache
	chain string
	tmp   types.MemStore
	store types.Store
}

func NewState(store types.Store) *State {
	s := &State{
		tmp:   types.NewMemStore(),
		store: store,
	}
	return s
}

func (s *State) SetChain(chain string) {
	s.chain = chain
}

func (s *State) GetChain() string {
	if len(s.chain) == 0 {
		PanicSanity("Expected to have set chain")
	}
	return s.chain
}

func (s *State) Get(key []byte) []byte {
	if s.tmp != nil {
		if value := s.tmp.Get(key); value != nil {
			return value
		}
	}
	return s.store.Get(key)
}

func (s *State) Set(key, value []byte) {
	if s.tmp != nil {
		s.tmp[string(key)] = value
	}
	s.store.Set(key, value)
}

func (s *State) GetAccount(addr []byte) *types.Account {
	return GetAccount(s.store, addr)
}

func (s *State) SetAccount(addr []byte, acc *types.Account) {
	SetAccount(s.store, addr, acc)
}

func (s *State) CacheWrap() *State {
	cache := types.NewCache(s.store)
	snew := &State{
		cache: cache,
		chain: s.chain,
		tmp:   nil,
		store: cache,
	}
	return snew
}

func (s *State) CacheSync() {
	s.cache.Sync()
}

func (s *State) Commit() tmsp.Result {
	s.tmp = types.NewMemStore()
	return s.store.(*merk.Client).CommitSync()
}

func AccountKey(addr []byte) []byte {
	return append([]byte("base/a/"), addr...)
}

func GetAccount(store types.Store, addr []byte) (acc *types.Account) {
	data := store.Get(AccountKey(addr))
	if len(data) == 0 {
		return nil
	}
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic(Fmt("Error reading account %X error: %s", data, err.Error()))
	}
	return acc
}

func SetAccount(store types.Store, addr []byte, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	store.Set(AccountKey(addr), accBytes)
}
