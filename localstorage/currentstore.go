package localstorage

import (
	"io/ioutil"
	"os"
	"path/filepath"

<<<<<<< HEAD
	"github.com/git-lfs/git-lfs/config"
	"github.com/git-lfs/git-lfs/errors"
<<<<<<< HEAD
=======
	"github.com/git-lfs/git-lfs/tools/longpathos"
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
=======
	"github.com/github/git-lfs/config"
	"github.com/github/git-lfs/errors"
>>>>>>> refs/remotes/git-lfs/register-commands-v2
)

const (
	tempDirPerms       = 0755
	localMediaDirPerms = 0755
	localLogDirPerms   = 0755
)

var (
	objects        *LocalStorage
	notInRepoErr   = errors.New("not in a repository")
	TempDir        = filepath.Join(os.TempDir(), "git-lfs")
	checkedTempDir string
)

func Objects() *LocalStorage {
	return objects
}

<<<<<<< HEAD
func InitStorage() error {
	if len(config.LocalGitStorageDir) == 0 || len(config.LocalGitDir) == 0 {
		return notInRepoErr
	}

=======
func ResolveDirs() error {
	config.ResolveGitBasicDirs()
>>>>>>> refs/remotes/git-lfs/register-commands-v2
	TempDir = filepath.Join(config.LocalGitDir, "lfs", "tmp") // temp files per worktree
	objs, err := NewStorage(
		filepath.Join(config.LocalGitStorageDir, "lfs", "objects"),
		filepath.Join(TempDir, "objects"),
	)

	if err != nil {
<<<<<<< HEAD
		return errors.Wrap(err, "init LocalStorage")
=======
		return errors.Wrap(err, "localstorage")
>>>>>>> refs/remotes/git-lfs/register-commands-v2
	}

	objects = objs
	config.LocalLogDir = filepath.Join(objs.RootDir, "logs")
<<<<<<< HEAD
	if err := os.MkdirAll(config.LocalLogDir, localLogDirPerms); err != nil {
<<<<<<< HEAD
=======
	if err := longpathos.MkdirAll(config.LocalLogDir, localLogDirPerms); err != nil {
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
		return errors.Wrap(err, "create log dir")
	}

	return nil
}

func InitStorageOrFail() {
	if err := InitStorage(); err != nil {
		if err == notInRepoErr {
			return
		}

		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
=======
		return errors.Wrap(err, "localstorage")
>>>>>>> refs/remotes/git-lfs/register-commands-v2
	}
	return nil
}

func ResolveDirs() {
	config.ResolveGitBasicDirs()
	InitStorageOrFail()
}

func TempFile(prefix string) (*os.File, error) {
	if checkedTempDir != TempDir {
		if err := longpathos.MkdirAll(TempDir, tempDirPerms); err != nil {
			return nil, err
		}
		checkedTempDir = TempDir
	}

	return ioutil.TempFile(TempDir, prefix)
}

func ResetTempDir() error {
	checkedTempDir = ""
	return longpathos.RemoveAll(TempDir)
}
