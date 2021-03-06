package locking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

<<<<<<< HEAD
	"github.com/git-lfs/git-lfs/errors"
	"github.com/git-lfs/git-lfs/filepathfilter"
	"github.com/git-lfs/git-lfs/git"
	"github.com/git-lfs/git-lfs/tools"
=======
	"github.com/github/git-lfs/errors"
	"github.com/github/git-lfs/tools"

	"github.com/github/git-lfs/config"
)

var (
	// lockable patterns from .gitattributes
	cachedLockablePatterns []string
	cachedLockableMutex    sync.Mutex
>>>>>>> refs/remotes/git-lfs/locking-workflow
)

// GetLockablePatterns returns a list of patterns in .gitattributes which are
// marked as lockable
<<<<<<< HEAD
func (c *Client) GetLockablePatterns() []string {
	c.ensureLockablesLoaded()
	return c.lockablePatterns
}

// getLockableFilter returns the internal filter used to check if a file is lockable
func (c *Client) getLockableFilter() *filepathfilter.Filter {
	c.ensureLockablesLoaded()
	return c.lockableFilter
}

func (c *Client) ensureLockablesLoaded() {
	c.lockableMutex.Lock()
	defer c.lockableMutex.Unlock()

	// Only load once
	if c.lockablePatterns == nil {
		c.refreshLockablePatterns()
	}
}

// Internal function to repopulate lockable patterns
// You must have locked the c.lockableMutex in the caller
func (c *Client) refreshLockablePatterns() {

	paths := git.GetAttributePaths(c.LocalWorkingDir, c.LocalGitDir)
	// Always make non-nil even if empty
	c.lockablePatterns = make([]string, 0, len(paths))
	for _, p := range paths {
		if p.Lockable {
			c.lockablePatterns = append(c.lockablePatterns, p.Path)
		}
	}
	c.lockableFilter = filepathfilter.New(c.lockablePatterns, nil)
=======
func GetLockablePatterns() []string {
	cachedLockableMutex.Lock()
	defer cachedLockableMutex.Unlock()

	// Only load once
	if cachedLockablePatterns == nil {
		// Always make non-nil even if empty
		cachedLockablePatterns = make([]string, 0, 10)

		paths := config.GetAttributePaths()
		for _, p := range paths {
			if p.Lockable {
				cachedLockablePatterns = append(cachedLockablePatterns, p.Path)
			}
		}
	}

	return cachedLockablePatterns

}

// RefreshLockablePatterns causes us to re-read the .gitattributes and caches the result
func RefreshLockablePatterns() {
	cachedLockableMutex.Lock()
	defer cachedLockableMutex.Unlock()
	cachedLockablePatterns = nil
>>>>>>> refs/remotes/git-lfs/locking-workflow
}

// IsFileLockable returns whether a specific file path is marked as Lockable,
// ie has the 'lockable' attribute in .gitattributes
// Lockable patterns are cached once for performance, unless you call RefreshLockablePatterns
// path should be relative to repository root
<<<<<<< HEAD
func (c *Client) IsFileLockable(path string) bool {
	return c.getLockableFilter().Allows(path)
=======
func IsFileLockable(path string) bool {
	return tools.PathMatchesWildcardPatterns(path, GetLockablePatterns())
>>>>>>> refs/remotes/git-lfs/locking-workflow
}

// FixAllLockableFileWriteFlags recursively scans the repo looking for files which
// are lockable, and makes sure their write flags are set correctly based on
// whether they are currently locked or unlocked.
// Files which are unlocked are made read-only, files which are locked are made
// writeable.
// This function can be used after a clone or checkout to ensure that file
// state correctly reflects the locking state
<<<<<<< HEAD
func (c *Client) FixAllLockableFileWriteFlags() error {
	return c.fixFileWriteFlags(c.LocalWorkingDir, c.LocalWorkingDir, c.getLockableFilter(), nil)
=======
func FixAllLockableFileWriteFlags() error {
	return FixFileWriteFlagsInDir("", GetLockablePatterns(), nil, true)
>>>>>>> refs/remotes/git-lfs/locking-workflow
}

// FixFileWriteFlagsInDir scans dir (which can either be a relative dir
// from the root of the repo, or an absolute dir within the repo) looking for
// files to change permissions for.
// If lockablePatterns is non-nil, then any file matching those patterns will be
// checked to see if it is currently locked by the current committer, and if so
// it will be writeable, and if not locked it will be read-only.
// If unlockablePatterns is non-nil, then any file matching those patterns will
// be made writeable if it is not already. This can be used to reset files to
// writeable when their 'lockable' attribute is turned off.
<<<<<<< HEAD
func (c *Client) FixFileWriteFlagsInDir(dir string, lockablePatterns, unlockablePatterns []string) error {
=======
func FixFileWriteFlagsInDir(dir string, lockablePatterns, unlockablePatterns []string, recursive bool) error {
>>>>>>> refs/remotes/git-lfs/locking-workflow

	// early-out if no patterns
	if len(lockablePatterns) == 0 && len(unlockablePatterns) == 0 {
		return nil
	}

	absPath := dir
	if !filepath.IsAbs(dir) {
<<<<<<< HEAD
		absPath = filepath.Join(c.LocalWorkingDir, dir)
=======
		absPath = filepath.Join(config.LocalWorkingDir, dir)
>>>>>>> refs/remotes/git-lfs/locking-workflow
	}
	stat, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%q is not a valid directory", dir)
	}

<<<<<<< HEAD
	var lockableFilter *filepathfilter.Filter
	var unlockableFilter *filepathfilter.Filter
	if lockablePatterns != nil {
		lockableFilter = filepathfilter.New(lockablePatterns, nil)
	}
	if unlockablePatterns != nil {
		unlockableFilter = filepathfilter.New(unlockablePatterns, nil)
	}

	return c.fixFileWriteFlags(absPath, c.LocalWorkingDir, lockableFilter, unlockableFilter)
}

// Internal implementation of fixing file write flags with precompiled filters
func (c *Client) fixFileWriteFlags(absPath, workingDir string, lockable, unlockable *filepathfilter.Filter) error {

	var errs []error
	var errMux sync.Mutex

	addErr := func(err error) {
		errMux.Lock()
		defer errMux.Unlock()

		errs = append(errs, err)
	}

	tools.FastWalkGitRepo(absPath, func(parentDir string, fi os.FileInfo, err error) {
		if err != nil {
			addErr(err)
			return
		}
		// Skip dirs, we only need to check files
		if fi.IsDir() {
			return
		}
		abschild := filepath.Join(parentDir, fi.Name())

		// This is a file, get relative to repo root
		relpath, err := filepath.Rel(workingDir, abschild)
		if err != nil {
			addErr(err)
			return
		}

		err = c.fixSingleFileWriteFlags(relpath, lockable, unlockable)
		if err != nil {
			addErr(err)
		}

	})
=======
	// For simplicity, don't use goroutines to parallelise recursive scan
	// This routine is almost certainly disk-limited anyway
	// We don't need sorting so don't use ioutil.Readdir or filepath.Walk
	d, err := os.Open(absPath)
	if err != nil {
		return err
	}

	contents, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	var errs []error
	for _, fi := range contents {
		abschild := filepath.Join(absPath, fi.Name())
		if fi.IsDir() {
			if recursive {
				err = FixFileWriteFlagsInDir(abschild, lockablePatterns, unlockablePatterns, recursive)
			}
			continue
		}

		// This is a file, get relative to repo root
		relpath, err := filepath.Rel(config.LocalWorkingDir, abschild)
		if err != nil {
			return err
		}

		err = fixSingleFileWriteFlags(relpath, lockablePatterns, unlockablePatterns)
		if err != nil {
			errs = append(errs, err)
		}

	}
>>>>>>> refs/remotes/git-lfs/locking-workflow
	return errors.Combine(errs)
}

// FixLockableFileWriteFlags checks each file in the provided list, and for
// those which are lockable, makes sure their write flags are set correctly
// based on whether they are currently locked or unlocked. Files which are
// unlocked are made read-only, files which are locked are made writeable.
// Files which are not lockable are ignored.
// This function can be used after a clone or checkout to ensure that file
// state correctly reflects the locking state, and is more efficient than
// FixAllLockableFileWriteFlags when you know which files changed
<<<<<<< HEAD
func (c *Client) FixLockableFileWriteFlags(files []string) error {
	// early-out if no lockable patterns
	if len(c.GetLockablePatterns()) == 0 {
=======
func FixLockableFileWriteFlags(files []string) error {
	lockablePatterns := GetLockablePatterns()

	// early-out if no lockable patterns
	if len(lockablePatterns) == 0 {
>>>>>>> refs/remotes/git-lfs/locking-workflow
		return nil
	}

	var errs []error
	for _, f := range files {
<<<<<<< HEAD
		err := c.fixSingleFileWriteFlags(f, c.getLockableFilter(), nil)
=======
		err := fixSingleFileWriteFlags(f, lockablePatterns, nil)
>>>>>>> refs/remotes/git-lfs/locking-workflow
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Combine(errs)
}

// fixSingleFileWriteFlags fixes write flags on a single file
// If lockablePatterns is non-nil, then any file matching those patterns will be
// checked to see if it is currently locked by the current committer, and if so
// it will be writeable, and if not locked it will be read-only.
// If unlockablePatterns is non-nil, then any file matching those patterns will
// be made writeable if it is not already. This can be used to reset files to
// writeable when their 'lockable' attribute is turned off.
<<<<<<< HEAD
func (c *Client) fixSingleFileWriteFlags(file string, lockable, unlockable *filepathfilter.Filter) error {
=======
func fixSingleFileWriteFlags(file string, lockablePatterns, unlockablePatterns []string) error {
>>>>>>> refs/remotes/git-lfs/locking-workflow
	// Convert to git-style forward slash separators if necessary
	// Necessary to match attributes
	if filepath.Separator == '\\' {
		file = strings.Replace(file, "\\", "/", -1)
	}
<<<<<<< HEAD
	if lockable != nil && lockable.Allows(file) {
		// Lockable files are writeable only if they're currently locked
		err := tools.SetFileWriteFlag(file, c.IsFileLockedByCurrentCommitter(file))
=======
	if tools.PathMatchesWildcardPatterns(file, lockablePatterns) {
		// Lockable files are writeable only if they're currently locked
		err := tools.SetFileWriteFlag(file, IsFileLockedByCurrentCommitter(file))
>>>>>>> refs/remotes/git-lfs/locking-workflow
		// Ignore not exist errors
		if err != nil && !os.IsNotExist(err) {
			return err
		}
<<<<<<< HEAD
	} else if unlockable != nil && unlockable.Allows(file) {
=======
	} else if tools.PathMatchesWildcardPatterns(file, unlockablePatterns) {
>>>>>>> refs/remotes/git-lfs/locking-workflow
		// Unlockable files are always writeable
		// We only check files which match the incoming patterns to avoid
		// checking every file in the system all the time, and only do it
		// when a file has had its lockable attribute removed
		err := tools.SetFileWriteFlag(file, true)
<<<<<<< HEAD
		if err != nil && !os.IsNotExist(err) {
=======
		if err != nil {
>>>>>>> refs/remotes/git-lfs/locking-workflow
			return err
		}
	}
	return nil
}
