package commands

import (
<<<<<<< HEAD
	"encoding/json"
	"os"

=======
	"github.com/github/git-lfs/locking"
>>>>>>> refs/remotes/git-lfs/locking-workflow
	"github.com/spf13/cobra"
)

var (
	locksCmdFlags = new(locksFlags)
)

func locksCommand(cmd *cobra.Command, args []string) {
<<<<<<< HEAD
=======

>>>>>>> refs/remotes/git-lfs/locking-workflow
	filters, err := locksCmdFlags.Filters()
	if err != nil {
		Exit("Error building filters: %v", err)
	}

<<<<<<< HEAD
	lockClient := newLockClient(lockRemote)
	defer lockClient.Close()

	var lockCount int
	locks, err := lockClient.SearchLocks(filters, locksCmdFlags.Limit, locksCmdFlags.Local)
	// Print any we got before exiting

	if locksCmdFlags.JSON {
		if err := json.NewEncoder(os.Stdout).Encode(locks); err != nil {
			Error(err.Error())
		}
		return
=======
	var lockCount int
	locks := locking.SearchLocks(lockRemote, filters, locksCmdFlags.Limit, locksCmdFlags.Local)

	for lock := range locks.Results {
		Print("%s\t%s <%s>", lock.Path, lock.Committer.Name, lock.Committer.Email)
		lockCount++
>>>>>>> refs/remotes/git-lfs/locking-workflow
	}
	err = locks.Wait()

<<<<<<< HEAD
	for _, lock := range locks {
		Print("%s\t%s", lock.Path, lock.Committer)
		lockCount++
	}

	if err != nil {
		Exit("Error while retrieving locks: %v", err)
	}
=======
	if err != nil {
		Exit("Error while retrieving locks: %v", err)
	}

>>>>>>> refs/remotes/git-lfs/locking-workflow
	Print("\n%d lock(s) matched query.", lockCount)
}

// locksFlags wraps up and holds all of the flags that can be given to the
// `git lfs locks` command.
type locksFlags struct {
	// Path is an optional filter parameter to filter against the lock's
	// path
	Path string
	// Id is an optional filter parameter used to filtere against the lock's
	// ID.
	Id string
	// limit is an optional request parameter sent to the server used to
	// limit the
	Limit int
	// local limits the scope of lock reporting to the locally cached record
	// of locks for the current user & doesn't query the server
	Local bool
<<<<<<< HEAD
	// JSON is an optional parameter to output data in json format.
	JSON bool
=======
>>>>>>> refs/remotes/git-lfs/locking-workflow
}

// Filters produces a filter based on locksFlags instance.
func (l *locksFlags) Filters() (map[string]string, error) {
	filters := make(map[string]string)

	if l.Path != "" {
		path, err := lockPath(l.Path)
		if err != nil {
			return nil, err
		}

		filters["path"] = path
	}
	if l.Id != "" {
		filters["id"] = l.Id
	}

	return filters, nil
}

func init() {
	if !isCommandEnabled(cfg, "locks") {
		return
	}

	RegisterCommand("locks", locksCommand, func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&lockRemote, "remote", "r", cfg.CurrentRemote, lockRemoteHelp)
		cmd.Flags().StringVarP(&locksCmdFlags.Path, "path", "p", "", "filter locks results matching a particular path")
		cmd.Flags().StringVarP(&locksCmdFlags.Id, "id", "i", "", "filter locks results matching a particular ID")
		cmd.Flags().IntVarP(&locksCmdFlags.Limit, "limit", "l", 0, "optional limit for number of results to return")
		cmd.Flags().BoolVarP(&locksCmdFlags.Local, "local", "", false, "only list cached local record of own locks")
<<<<<<< HEAD
		cmd.Flags().BoolVarP(&locksCmdFlags.JSON, "json", "", false, "print output in json")
=======
>>>>>>> refs/remotes/git-lfs/locking-workflow
	})
}
