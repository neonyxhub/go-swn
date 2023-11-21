package ds_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"go.neonyx.io/go-swn/pkg/ds"
	"go.neonyx.io/go-swn/pkg/ds/drivers"
)

func TestNew(t *testing.T) {
	cfg := &drivers.DataStoreCfg{
		Path: "mock",
		Options: opt.Options{
			BlockCacheCapacity: -1,
		},
	}
	driver, err := ds.New(cfg)

	require.NoError(t, err)

	_, err = os.Stat(cfg.Path)
	require.NoError(t, err)

	err = driver.Close()
	require.NoError(t, err)

	err = os.RemoveAll("mock")
	require.NoError(t, err)
}
