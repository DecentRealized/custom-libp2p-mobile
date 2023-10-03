package database

import (
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/ristretto/z"
	"os"
	"path/filepath"
	"time"
)

var databaseDirectory = "./.db"
var _db *badger.DB = nil
var timerCanceler = make(chan struct{})

func SetDatabaseDirectory(baseDirectory string) error {
	stat, err := os.Stat(baseDirectory)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return os.ErrInvalid
	}
	databaseDirectory = baseDirectory
	notifier.QueueInfo(fmt.Sprintf("Database directory set to: %v", databaseDirectory))
	return nil
}

func GetDatabaseDirectory() string {
	return databaseDirectory
}

func Init(node *models.Node) error {
	if _db != nil {
		err := Close()
		if err != nil {
			return err
		}
	}
	rawKey, err := node.Peerstore().PrivKey(node.ID()).Raw()
	if err != nil {
		return err
	}
	encryptionKey := rawKey[len(rawKey)-16:]
	nodeDbPath := filepath.Join(databaseDirectory, node.ID().String())
	options := badger.DefaultOptions(nodeDbPath).WithEncryptionKey(encryptionKey).WithSyncWrites(true).
		WithEncryptionKeyRotationDuration(time.Until(time.Unix(1<<63-1, 0))).WithValueLogFileSize(1 << 20)
	options.IndexCacheSize = 64 << 20 // 64 Mb
	open, err := badger.Open(options)
	if err != nil {
		return err
	}
	_db = open
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		select {
		case <-ticker.C:
			var err error = nil
			for err == nil {
				if _db != nil {
					err = _db.RunValueLogGC(0.7)
				} else {
					err = ErrDatabaseRunning
				}
			}
		case <-timerCanceler:
			break
		}
	}()
	return nil
}

func Close() error {
	if _db == nil {
		return ErrDatabaseStopped
	}
	err := _db.Close()
	if err != nil {
		return err
	}
	timerCanceler <- struct{}{}
	_db = nil
	return nil
}

func Store(key []byte, value []byte) error {
	if _db == nil {
		return ErrDatabaseStopped
	}
	return _db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func Load(key []byte) ([]byte, error) {
	if _db == nil {
		return nil, ErrDatabaseStopped
	}
	var _readValue []byte = nil
	err := _db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			_readValue = val
			return nil
		})
	})
	return _readValue, err
}

func Delete(key []byte) error {
	if _db == nil {
		return ErrDatabaseStopped
	}
	return _db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func Range(keyBase []byte, f func(k []byte, v []byte) bool) error {
	if _db == nil {
		return ErrDatabaseStopped
	}
	stream := _db.NewStream()
	stream.NumGo = 4
	stream.Prefix = keyBase
	stream.KeyToList = nil
	stream.Send = func(buf *z.Buffer) error {
		kvList, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range kvList.GetKv() {
			continueIter := f(kv.Key, kv.Value)
			if !continueIter {
				return nil
			}
		}
		return nil
	}
	if err := stream.Orchestrate(context.TODO()); err != nil {
		return err
	}
	return nil
}
