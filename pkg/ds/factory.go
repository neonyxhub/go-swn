package ds

import (
	"go.neonyx.io/go-swn/pkg/ds/drivers"
	"go.neonyx.io/go-swn/pkg/ds/drivers/leveldb"
)

// New creates a new instance of DataStore
func New(cfg *drivers.DataStoreCfg) (drivers.DataStore, error) {
	return leveldb.New(cfg)
}
