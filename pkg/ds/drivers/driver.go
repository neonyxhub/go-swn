package drivers

type DataStoreCfg struct {
	Path    string
	Options interface{}
}

type DataStore interface {
	Close() error
	NewKey(k ...string) []byte
	Get(key []byte, opts interface{}) ([]byte, error)
	Put(key []byte, value []byte, opts interface{}) error
}
