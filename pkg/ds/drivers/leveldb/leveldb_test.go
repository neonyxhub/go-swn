package leveldb_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"go.neonyx.io/go-swn/pkg/ds"
	"go.neonyx.io/go-swn/pkg/ds/drivers"
)

func new(t *testing.T) drivers.DataStore {
	cfg := &drivers.DataStoreCfg{
		Path: "mock",
		Options: opt.Options{
			BlockCacheCapacity: -1,
		},
	}
	driver, err := ds.New(cfg)
	require.NoError(t, err)
	return driver
}

func destroy(t *testing.T, driver drivers.DataStore) {
	err := driver.Close()
	require.NoError(t, err)

	err = os.RemoveAll("mock")
	require.NoError(t, err)
}

func TestNewKey(t *testing.T) {
	ldb := new(t)
	defer destroy(t, ldb)
	key := ldb.NewKey("newkey")
	require.Equal(t, key, []byte("/newkey/"))
}

func TestPut(t *testing.T) {
	ldb := new(t)
	defer destroy(t, ldb)
	key := ldb.NewKey("newkey")
	err := ldb.Put(key, []byte{0xbe, 0xef}, nil)
	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	ldb := new(t)
	defer destroy(t, ldb)
	key := ldb.NewKey("newkey")
	err := ldb.Put(key, []byte{0xbe, 0xef}, nil)
	require.NoError(t, err)

	val, err := ldb.Get(key, nil)
	require.NoError(t, err)
	require.Equal(t, val, []byte{0xbe, 0xef})
}
