package db

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"

	badger "github.com/dgraph-io/badger/v4"
)

type NS []byte

// create a default badger db config
func DefaultKvDBConfig() *KvDBConfig {
	return &KvDBConfig{
		DataDir:        "cache/db",
		GcDiscardRatio: 0.5,
		GcInterval:     10 * time.Minute,
		Logger:         zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

type (
	// BadgerDB config
	KvDBConfig struct {
		// DataDir defines the directory where the BadgerDB database will be stored.
		DataDir        string
		GcDiscardRatio float64
		GcInterval     time.Duration
		Logger         zerolog.Logger // Embed a logger
	}

	// DB defines an embedded key/value store database interface.
	// we could think of the namespace as a table name or a collection name
	DB interface {
		Get(namespace, key []byte) (value []byte, err error)
		Set(namespace, key, value []byte) error
		Has(namespace, key []byte) (bool, error)
		Close() error
	}

	// BadgerDB is a wrapper around a BadgerDB backend database that implements
	// the DB interface.
	BadgerDB struct {
		db         *badger.DB
		ctx        context.Context
		cancelFunc context.CancelFunc
		cfg        *KvDBConfig
	}
)

// NewBadgerDB returns a new initialized BadgerDB database implementing the DB
// interface. If the database cannot be initialized, an error will be returned.
func NewBadgerDB(cfg KvDBConfig) (DB, error) {
	if err := os.MkdirAll(cfg.DataDir, 0774); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(cfg.DataDir)

	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	bdb := &BadgerDB{
		db:  badgerDB,
		cfg: &cfg,
	}
	bdb.ctx, bdb.cancelFunc = context.WithCancel(context.Background())

	go bdb.runGC()
	return bdb, nil
}

// Get implements the DB interface. It attempts to get a value for a given key
// and namespace. If the key does not exist in the provided namespace, an error
// is returned, otherwise the retrieved value.
func (bdb *BadgerDB) Get(namespace, key []byte) (value []byte, err error) {
	cKey := composeNSKey(namespace, key)

	err = bdb.db.View(func(txn *badger.Txn) error {
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

// Set implements the DB interface. It attempts to store a value for a given key
// and namespace. If the key/value pair cannot be saved, an error is returned.
func (bdb *BadgerDB) Set(namespace, key, value []byte) error {
	cKey := composeNSKey(namespace, key)
	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(cKey, value)
	})

	if err != nil {
		bdb.cfg.Logger.Debug().Msgf("failed to set key %s for namespace %s: %v", key, namespace, err)
		return err
	}

	return nil
}

// Has implements the DB interface. It returns a boolean reflecting if the
// database has a given key for a namespace or not. An error is only returned if
// an error to Get would be returned that is not of type badger.ErrKeyNotFound.
func (bdb *BadgerDB) Has(namespace, key []byte) (ok bool, err error) {
	_, err = bdb.Get(namespace, key)
	switch err {
	case badger.ErrKeyNotFound:
		ok, err = false, nil
	case nil:
		ok, err = true, nil
	}

	return
}

// Close implements the DB interface. It closes the connection to the underlying
// BadgerDB database as well as invoking the context's cancel function.
func (bdb *BadgerDB) Close() error {
	bdb.cancelFunc()
	return bdb.db.Close()
}

// runGC triggers the garbage collection for the BadgerDB backend database. It
// should be run in a goroutine.
func (bdb *BadgerDB) runGC() {
	ticker := time.NewTicker(bdb.cfg.GcInterval)
	for {
		select {
		case <-ticker.C:
			err := bdb.db.RunValueLogGC(bdb.cfg.GcDiscardRatio)
			if err != nil {
				// don't report error when GC didn't result in any cleanup
				if err == badger.ErrNoRewrite {
					bdb.cfg.Logger.Printf("no BadgerDB GC occurred: %v", err)
				} else {
					bdb.cfg.Logger.Error().Msgf("failed to GC BadgerDB: %v", err)
				}
			}

		case <-bdb.ctx.Done():
			return
		}
	}
}

func CreateNamespace(namespace string) NS {
	return []byte(namespace + "/")
}

// composeNSKey returns a composite key used for lookup and storage for a
// given namespace and key.
func composeNSKey(namespace NS, key []byte) []byte {
	return append(namespace, key...)
}
