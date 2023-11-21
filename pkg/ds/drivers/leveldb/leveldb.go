package leveldb

import (
	"errors"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"go.neonyx.io/go-swn/pkg/ds/drivers"
)

var (
	ErrInvalidOpts      = errors.New("invalid leveldb options")
	ErrInvalidReadOpts  = errors.New("invalid leveldb read options")
	ErrInvalidWriteOpts = errors.New("invalid leveldb write options")
)

type LevelDB struct {
	Path    string
	DB      *leveldb.DB
	Options *opt.Options
}

// New creates a new LevelDB instance with passed cfg options
func New(cfg *drivers.DataStoreCfg) (drivers.DataStore, error) {
	o, ok := cfg.Options.(opt.Options)
	if !ok {
		return nil, errors.New("should be leveldb/opt.Options in cfg.Options")
	}

	ldb := &LevelDB{
		Path:    cfg.Path,
		Options: &o,
	}

	db, err := leveldb.OpenFile(ldb.Path, ldb.Options)
	if err != nil {
		return nil, err
	}
	ldb.DB = db

	return ldb, nil
}

func (ldb *LevelDB) Close() error {
	return ldb.DB.Close()
}

// Make compose key from subkeys
func (ldb *LevelDB) NewKey(k ...string) []byte {
	var ans string = ""

	for _, key := range k {
		ans = ans + "/" + strings.Trim(key, "/")
	}

	if !strings.HasSuffix(ans, "/") {
		ans = ans + "/"
	}

	return []byte(ans)
}

func (ldb *LevelDB) Get(key []byte, opts interface{}) ([]byte, error) {
	if opts != nil {
		ro, ok := opts.(opt.ReadOptions)
		if !ok {
			return nil, ErrInvalidReadOpts
		}
		return ldb.DB.Get(key, &ro)
	}
	return ldb.DB.Get(key, nil)
}

func (ldb *LevelDB) Put(key []byte, value []byte, opts interface{}) error {
	if opts != nil {
		wo, ok := opts.(opt.WriteOptions)
		if !ok {
			return ErrInvalidWriteOpts
		}
		return ldb.DB.Put(key, value, &wo)
	}
	return ldb.DB.Put(key, value, nil)
}
