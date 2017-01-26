package commands

import (
<<<<<<< HEAD
=======
	"bytes"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/git-lfs/git-lfs/errors"
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	"github.com/git-lfs/git-lfs/filepathfilter"
	"github.com/git-lfs/git-lfs/git"
	"github.com/git-lfs/git-lfs/lfs"
	"github.com/git-lfs/git-lfs/progress"
<<<<<<< HEAD
=======
	"github.com/rubyist/tracerx"
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	"github.com/spf13/cobra"
)

func checkoutCommand(cmd *cobra.Command, args []string) {
	requireInRepo()
<<<<<<< HEAD
=======

	// Parameters are filters
	// firstly convert any pathspecs to the root of the repo, in case this is being executed in a sub-folder
	var rootedpaths []string

	inchan := make(chan string, 1)
	outchan, err := lfs.ConvertCwdFilesRelativeToRepo(inchan)
	if err != nil {
		Panic(err, "Could not checkout")
	}
	for _, arg := range args {
		inchan <- arg
		rootedpaths = append(rootedpaths, <-outchan)
	}
	close(inchan)
	checkoutWithIncludeExclude(filepathfilter.New(rootedpaths, nil))
}

// Checkout from items reported from the fetch process (in parallel)
func checkoutAllFromFetchChan(c chan *lfs.WrappedPointer) {
	tracerx.Printf("starting fetch/parallel checkout")
	checkoutFromFetchChan(nil, c)
}

func checkoutFromFetchChan(filter *filepathfilter.Filter, in chan *lfs.WrappedPointer) {
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	ref, err := git.CurrentRef()
	if err != nil {
		Panic(err, "Could not checkout")
	}

<<<<<<< HEAD
	var totalBytes int64
	meter := progress.NewMeter(progress.WithOSEnv(cfg.Os))
	singleCheckout := newSingleCheckout()
	chgitscanner := lfs.NewGitScanner(func(p *lfs.WrappedPointer, err error) {
		if err != nil {
			LoggedError(err, "Scanner error")
			return
=======
	// Map oid to multiple pointers
	mapping := make(map[string][]*lfs.WrappedPointer)
	for _, pointer := range pointers {
		if filter.Allows(pointer.Name) {
			mapping[pointer.Oid] = append(mapping[pointer.Oid], pointer)
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
		}

		totalBytes += p.Size
		meter.Add(p.Size)
		meter.StartTransfer(p.Name)

		singleCheckout.Run(p)

		// not strictly correct (parallel) but we don't have a callback & it's just local
		// plus only 1 slot in channel so it'll block & be close
		meter.TransferBytes("checkout", p.Name, p.Size, totalBytes, int(p.Size))
		meter.FinishTransfer(p.Name)
	})

	chgitscanner.Filter = filepathfilter.New(rootedPaths(args), nil)

<<<<<<< HEAD
	if err := chgitscanner.ScanTree(ref.Sha); err != nil {
		ExitWithError(err)
	}

	meter.Start()
	chgitscanner.Close()
	meter.Finish()
	singleCheckout.Close()
}

// Parameters are filters
// firstly convert any pathspecs to the root of the repo, in case this is being
// executed in a sub-folder
func rootedPaths(args []string) []string {
	pathConverter, err := lfs.NewCurrentToRepoPathConverter()
=======
func checkoutWithIncludeExclude(filter *filepathfilter.Filter) {
	ref, err := git.CurrentRef()
	if err != nil {
		Panic(err, "Could not checkout")
	}

	pointers, err := lfs.ScanTree(ref.Sha)
	if err != nil {
		Panic(err, "Could not scan for Git LFS files")
	}

	var wait sync.WaitGroup
	wait.Add(1)

	c := make(chan *lfs.WrappedPointer, 1)

	go func() {
		checkoutWithChan(c)
		wait.Done()
	}()

	// Count bytes for progress
	var totalBytes int64
	for _, pointer := range pointers {
		totalBytes += pointer.Size
	}

	logPath, _ := cfg.Os.Get("GIT_LFS_PROGRESS")
	progress := progress.NewProgressMeter(len(pointers), totalBytes, false, logPath)
	progress.Start()
	totalBytes = 0
	for _, pointer := range pointers {
		totalBytes += pointer.Size
		if filter.Allows(pointer.Name) {
			progress.Add(pointer.Name)
			c <- pointer
			// not strictly correct (parallel) but we don't have a callback & it's just local
			// plus only 1 slot in channel so it'll block & be close
			progress.TransferBytes("checkout", pointer.Name, pointer.Size, totalBytes, int(pointer.Size))
			progress.FinishTransfer(pointer.Name)
		} else {
			progress.Skip(pointer.Size)
		}
	}
	close(c)
	wait.Wait()
	progress.Finish()

}

// Populate the working copy with the real content of objects where the file is
// either missing, or contains a matching pointer placeholder, from a list of pointers.
// If the file exists but has other content it is left alone
// Callers of this function MUST NOT Panic or otherwise exit the process
// without waiting for this function to shut down.  If the process exits while
// update-index is in the middle of processing a file the git index can be left
// in a locked state.
func checkoutWithChan(in <-chan *lfs.WrappedPointer) {
	// Get a converter from repo-relative to cwd-relative
	// Since writing data & calling git update-index must be relative to cwd
	repopathchan := make(chan string, 1)
	cwdpathchan, err := lfs.ConvertRepoFilesRelativeToCwd(repopathchan)
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	if err != nil {
		Panic(err, "Could not checkout")
	}

	rootedpaths := make([]string, 0, len(args))
	for _, arg := range args {
		rootedpaths = append(rootedpaths, pathConverter.Convert(arg))
	}
	return rootedpaths
}

func init() {
	RegisterCommand("checkout", checkoutCommand, nil)
}
