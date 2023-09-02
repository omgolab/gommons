package gcdb

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"

	badger "github.com/dgraph-io/badger/v4"
)

var defaultOptions []KvDBOption

type (
	KvDBOption func(*rootConfig) error

	// badgerDB rootConfig
	rootConfig struct {
		// dataDir defines the directory where the badgerDB database will be stored.
		dataDir        string
		gcDiscardRatio float64
		gcInterval     time.Duration
		logger         zerolog.Logger // Embed a logger
	}

	DB interface {
		CreateNsCollection(ns string) Collection
		Close() error
	}

	// Collection defines an embedded key/value store database interface.
	// we could think of the Collection as a SQL table name or a mongodb collection name
	Collection interface {
		Get(key []byte) (value []byte, err error)
		Set(key, value []byte) error
		Has(key []byte) (bool, error)
	}

	// db is a wrapper around a db backend database that implements
	// the DB interface.
	db struct {
		badgerDb   *badger.DB
		ctx        context.Context
		cancelFunc context.CancelFunc
		logger     zerolog.Logger
	}

	// collection is a wrapper around the backend database and a table/collection namespace
	collection struct {
		ns  []byte
		rdb *db
	}
)

// create a default badger db config
func getDefaultOptions() []KvDBOption {
	if defaultOptions != nil {
		return defaultOptions
	}

	defaultOptions = []KvDBOption{
		WithDataDir("cache/db"),
		WithGcDiscardRatio(0.5),
		WithGcInterval(10 * time.Minute),
		WithLogger(zerolog.New(os.Stderr).With().Timestamp().Logger()),
	}
	return defaultOptions
}

// NewBadgerDB returns a new initialized badgerDB database implementing the DB
// interface. If the database cannot be initialized, an error will be returned.
func NewBadgerDB(kvOpts ...KvDBOption) (DB, error) {
	var cfg = &rootConfig{}
	var err error
	for _, opt := range getDefaultOptions() {
		err = opt(cfg)
		if err != nil {
			return nil, err
		}
	}
	for _, opt := range kvOpts {
		err = opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(cfg.dataDir, 0774); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(cfg.dataDir)
	bDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	bdb := &db{
		badgerDb: bDB,
		logger:   cfg.logger,
	}
	bdb.ctx, bdb.cancelFunc = context.WithCancel(context.Background())

	go bdb.runGC(cfg.gcInterval, cfg.gcDiscardRatio)
	return bdb, nil
}

// Get implements the DB interface. It attempts to get a value for a given key.
// If the key does not exist in the provided collection, an error
// is returned, otherwise the retrieved value.
func (t *collection) Get(key []byte) (value []byte, err error) {
	err = t.rdb.badgerDb.View(func(txn *badger.Txn) error {
		cKey := append(t.ns, key...)
		item, err := txn.Get(cKey)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

// Set implements the DB interface. It attempts to store a value for a given key.
// If the key/value pair cannot be saved, an error is returned.
func (t *collection) Set(key, value []byte) error {
	err := t.rdb.badgerDb.Update(func(txn *badger.Txn) error {
		cKey := append(t.ns, key...)
		return txn.Set(cKey, value)
	})

	if err != nil {
		t.rdb.logger.Debug().Msgf("failed to set key %s for the collection %s: %v", key, t.ns, err)
		return err
	}

	return nil
}

// Has implements the DB interface. It returns a boolean reflecting if the
// database has a given key for a ns or not. An error is only returned if
// an error to Get would be returned that is not of type badger.ErrKeyNotFound.
func (t *collection) Has(key []byte) (ok bool, err error) {
	_, err = t.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		ok, err = false, nil
	case nil:
		ok, err = true, nil
	}

	return
}

// Close implements the DB interface. It closes the connection to the underlying
// badgerDB database as well as invoking the context's cancel function.
func (bdb *db) Close() error {
	bdb.cancelFunc()
	return bdb.badgerDb.Close()
}

// runGC triggers the garbage collection for the badgerDB backend database. It
// should be run in a goroutine.
func (bdb *db) runGC(gcInterval time.Duration, gcDiscardRatio float64) {
	ticker := time.NewTicker(gcInterval)
	for {
		select {
		case <-ticker.C:
			err := bdb.badgerDb.RunValueLogGC(gcDiscardRatio)
			if err != nil {
				// don't report error when GC didn't result in any cleanup
				if err == badger.ErrNoRewrite {
					bdb.logger.Printf("no badgerDB GC occurred: %v", err)
				} else {
					bdb.logger.Error().Msgf("failed to GC badgerDB: %v", err)
				}
			}

		case <-bdb.ctx.Done():
			return
		}
	}
}

// CreateNsCollection returns a namespace (similar to a SQL table or MongoDB collection)
// internally a byte slice that can be used as a prefix for all keys
func (db *db) CreateNsCollection(name string) Collection {
	return &collection{
		ns:  []byte(name + "/"),
		rdb: db,
	}
}

// WithLogger sets the logger for the badgerDB database.
func WithLogger(logger zerolog.Logger) KvDBOption {
	return func(cfg *rootConfig) error {
		cfg.logger = logger
		return nil
	}
}

// WithGcDiscardRatio sets the garbage collection discard ratio for the badgerDB database.
func WithGcDiscardRatio(ratio float64) KvDBOption {
	return func(cfg *rootConfig) error {
		cfg.gcDiscardRatio = ratio
		return nil
	}
}

// WithGcInterval sets the garbage collection interval for the badgerDB database.
func WithGcInterval(interval time.Duration) KvDBOption {
	return func(cfg *rootConfig) error {
		cfg.gcInterval = interval
		return nil
	}
}

// WithDataDir sets the data directory for the badgerDB database.
func WithDataDir(dir string) KvDBOption {
	return func(cfg *rootConfig) error {
		cfg.dataDir = dir
		return nil
	}
}
