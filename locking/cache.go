package locking

import (
<<<<<<< HEAD
	"strings"

	"github.com/git-lfs/git-lfs/tools/kv"
)

const (
	// We want to use a single cache file for integrity, but to make it easy to
	// list all locks, prefix the id->path map in a way we can identify (something
	// that won't be in a path)
	idKeyPrefix string = "*id*://"
)

type LockCache struct {
	kv *kv.Store
}

func NewLockCache(filepath string) (*LockCache, error) {
	kv, err := kv.NewStore(filepath)
	if err != nil {
		return nil, err
	}
	return &LockCache{kv}, nil
}

// Cache a successful lock for faster local lookup later
func (c *LockCache) Add(l Lock) error {
	// Store reference in both directions
	// Path -> Lock
	c.kv.Set(l.Path, &l)
	// EncodedId -> Lock (encoded so we can easily identify)
	c.kv.Set(c.encodeIdKey(l.Id), &l)
	return nil
}

// Remove a cached lock by path becuase it's been relinquished
func (c *LockCache) RemoveByPath(filePath string) error {
	ilock := c.kv.Get(filePath)
	if lock, ok := ilock.(*Lock); ok && lock != nil {
		c.kv.Remove(lock.Path)
		// Id as key is encoded
		c.kv.Remove(c.encodeIdKey(lock.Id))
=======
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/github/git-lfs/api"

	"github.com/boltdb/bolt"
	"github.com/github/git-lfs/config"
)

var (
	dbInit             sync.Once
	lockDb             *bolt.DB
	pathToIdBucketName = []byte("pathToId")
	idToPathBucketName = []byte("idToPath")
)

func initDb() error {
	// Open on demand - bolt will lock this file & other processes trying to access it
	// we'll wait a max 5 seconds
	// TODO: could have option to open read-only to take shared lock
	var initerr error
	dbInit.Do(func() {
		lockDir := filepath.Join(config.LocalGitStorageDir, "lfs")
		os.MkdirAll(lockDir, 0755)
		lockFile := filepath.Join(lockDir, "db.lock")
		lockDb, initerr = bolt.Open(lockFile, 0644, &bolt.Options{Timeout: 5 * time.Second})
	})
	if initerr != nil {
		// TODO maybe suggest re-initialising lock cache
		return initerr
	}

	// Very important that Cleanup() is called before shutdown

	return nil
}

func Cleanup() error {
	if lockDb != nil {
		return lockDb.Close()
>>>>>>> refs/remotes/git-lfs/locking-workflow
	}
	return nil
}

<<<<<<< HEAD
// Remove a cached lock by id because it's been relinquished
func (c *LockCache) RemoveById(id string) error {
	// Id as key is encoded
	idkey := c.encodeIdKey(id)
	ilock := c.kv.Get(idkey)
	if lock, ok := ilock.(*Lock); ok && lock != nil {
		c.kv.Remove(idkey)
		c.kv.Remove(lock.Path)
	}
	return nil
}

// Get the list of cached locked files
func (c *LockCache) Locks() []Lock {
	var locks []Lock
	c.kv.Visit(func(key string, val interface{}) bool {
		// Only report file->id entries not reverse
		if !c.isIdKey(key) {
			lock := val.(*Lock)
			locks = append(locks, *lock)
		}
		return true // continue
	})
	return locks
}

// Clear the cache
func (c *LockCache) Clear() {
	c.kv.RemoveAll()
}

// Save the cache
func (c *LockCache) Save() error {
	return c.kv.Save()
}

func (c *LockCache) encodeIdKey(id string) string {
	// Safety against accidents
	if !c.isIdKey(id) {
		return idKeyPrefix + id
	}
	return id
}

func (c *LockCache) decodeIdKey(key string) string {
	// Safety against accidents
	if c.isIdKey(key) {
		return key[len(idKeyPrefix):]
	}
	return key
}

func (c *LockCache) isIdKey(key string) bool {
	return strings.HasPrefix(key, idKeyPrefix)
=======
// Run a read-only lock database function
// Deals with initialisation and function is a transaction
// Can run in a goroutine if needed
func runLockDbReadOnlyFunc(f func(tx *bolt.Tx) error) error {

	if err := initDb(); err != nil {
		return err
	}

	return lockDb.View(f)
}

// Run a read-write lock database function
// Deals with initialisation and function is a transaction
// Can run in a goroutine if needed
func runLockDbFunc(f func(tx *bolt.Tx) error) error {

	if err := initDb(); err != nil {
		return err
	}

	// Use Batch() to improve write performance and goroutine friendly
	return lockDb.Batch(f)
}

// This file caches active locks locally so that we can more easily retrieve
// a list of locally locked files without consulting the server
// This only includes locks which the local committer has taken, not all locks

// Cache a successful lock for faster local lookup later
func cacheLock(filePath, id string) error {
	return runLockDbFunc(func(tx *bolt.Tx) error {
		path2id, err := tx.CreateBucketIfNotExists(pathToIdBucketName)
		if err != nil {
			return err
		}
		id2path, err := tx.CreateBucketIfNotExists(idToPathBucketName)
		if err != nil {
			return err
		}
		// Store path -> id and id -> path
		if err := path2id.Put([]byte(filePath), []byte(id)); err != nil {
			return err
		}
		return id2path.Put([]byte(id), []byte(filePath))
	})
}

// Remove a cached lock by path becuase it's been relinquished
func cacheUnlock(filePath string) error {
	return runLockDbFunc(func(tx *bolt.Tx) error {
		path2id := tx.Bucket(pathToIdBucketName)
		id2path := tx.Bucket(idToPathBucketName)
		if path2id == nil || id2path == nil {
			return nil
		}
		idbytes := path2id.Get([]byte(filePath))
		if idbytes != nil {
			if err := id2path.Delete(idbytes); err != nil {
				return err
			}
			return path2id.Delete([]byte(filePath))
		}
		return nil
	})
}

// Remove a cached lock by id becuase it's been relinquished
func cacheUnlockById(id string) error {
	return runLockDbFunc(func(tx *bolt.Tx) error {
		path2id := tx.Bucket(pathToIdBucketName)
		id2path := tx.Bucket(idToPathBucketName)
		if path2id == nil || id2path == nil {
			return nil
		}
		pathbytes := id2path.Get([]byte(id))
		if pathbytes != nil {
			if err := path2id.Delete(pathbytes); err != nil {
				return err
			}
			return id2path.Delete([]byte(id))
		}
		return nil
	})
}

type CachedLock struct {
	Path string
	Id   string
}

// Get the list of cached locked files
func cachedLocks() []CachedLock {
	var ret []CachedLock
	runLockDbReadOnlyFunc(func(tx *bolt.Tx) error {
		path2id := tx.Bucket(pathToIdBucketName)
		if path2id == nil {
			return nil
		}
		path2id.ForEach(func(k []byte, v []byte) error {
			ret = append(ret, CachedLock{string(k), string(v)})
			return nil
		})
		return nil
	})
	return ret
}

// Fetch locked files for the current committer and cache them locally
// This can be used to sync up locked files when moving machines
func fetchLocksToCache(remoteName string) error {

	// TODO: filters don't seem to currently define how to search for a
	// committer's email. Is it "committer.email"? For now, just iterate
	lockWrapper := SearchLocks(remoteName, nil, 0, false)
	var locks []CachedLock
	email := api.CurrentCommitter().Email
	for l := range lockWrapper.Results {
		if l.Committer.Email == email {
			locks = append(locks, CachedLock{l.Path, l.Id})
		}
	}
	err := lockWrapper.Wait()

	if err != nil {
		return err
	}

	// replace cached locks (only do this if search was OK)
	return runLockDbFunc(func(tx *bolt.Tx) error {
		// Ignore errors deleting buckets
		tx.DeleteBucket(pathToIdBucketName)
		tx.DeleteBucket(idToPathBucketName)
		path2id, err := tx.CreateBucket(pathToIdBucketName)
		if err != nil {
			return err
		}
		id2path, err := tx.CreateBucket(idToPathBucketName)
		if err != nil {
			return err
		}
		for _, l := range locks {
			path2id.Put([]byte(l.Path), []byte(l.Id))
			id2path.Put([]byte(l.Id), []byte(l.Path))
		}
		return nil
	})
}

// IsFileLockedByCurrentCommitter returns whether a file is locked by the
// current committer, as cached locally
func IsFileLockedByCurrentCommitter(path string) bool {
	locked := false
	runLockDbReadOnlyFunc(func(tx *bolt.Tx) error {
		path2id := tx.Bucket(pathToIdBucketName)
		if path2id == nil {
			return nil
		}
		locked = path2id.Get([]byte(path)) != nil
		return nil
	})

	return locked
>>>>>>> refs/remotes/git-lfs/locking-workflow
}
