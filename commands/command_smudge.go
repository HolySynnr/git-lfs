package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/git-lfs/git-lfs/errors"
	"github.com/git-lfs/git-lfs/filepathfilter"
	"github.com/git-lfs/git-lfs/lfs"
<<<<<<< HEAD
=======
	"github.com/git-lfs/git-lfs/tools/longpathos"
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	"github.com/spf13/cobra"
)

var (
<<<<<<< HEAD
	// smudgeSkip is a command-line flag belonging to the "git-lfs smudge"
=======
	// smudgeInfo is a command-line flag belonging to the "git-lfs smudge"
	// command specifying whether to skip the smudge process and simply
	// print out the info of the files being smudged.
	//
	// As of v1.5.0, it is deprecated.
	smudgeInfo = false
	// smudgeInfo is a command-line flag belonging to the "git-lfs smudge"
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	// command specifying whether to skip the smudge process.
	smudgeSkip = false
)

// smudge smudges the given `*lfs.Pointer`, "ptr", and writes its objects
// contents to the `io.Writer`, "to".
//
// If the smudged object did not "pass" the include and exclude filterset, it
// will not be downloaded, and the object will remain a pointer on disk, as if
// the smudge filter had not been applied at all.
//
// Any errors encountered along the way will be returned immediately if they
// were non-fatal, otherwise execution will halt and the process will be
// terminated by using the `commands.Panic()` func.
<<<<<<< HEAD
func smudge(to io.Writer, from io.Reader, filename string, skip bool, filter *filepathfilter.Filter) error {
	ptr, pbuf, perr := lfs.DecodeFrom(from)
=======
func smudge(to io.Writer, ptr *lfs.Pointer, filename string, skip bool) error {
	cb, file, err := lfs.CopyCallbackFile("smudge", filename, 1, 1)
	if err != nil {
		return err
	}

	download := tools.FilenamePassesIncludeExcludeFilter(filename, cfg.FetchIncludePaths(), cfg.FetchExcludePaths())

	if skip || cfg.Os.Bool("GIT_LFS_SKIP_SMUDGE", false) {
		download = false
	}

	err = ptr.Smudge(to, filename, download, TransferManifest(), cb)
	if file != nil {
		file.Close()
	}

	if err != nil {
		ptr.Encode(to)
		// Download declined error is ok to skip if we weren't requesting download
		if !(errors.IsDownloadDeclinedError(err) && !download) {
			LoggedError(err, "Error downloading object: %s (%s)", filename, ptr.Oid)
			if !cfg.SkipDownloadErrors() {
				os.Exit(2)
			}
		}
	}

	return nil
}

func smudgeCommand(cmd *cobra.Command, args []string) {
	requireStdin("This command should be run by the Git 'smudge' filter")
	lfs.InstallHooks(false)

	// keeps the initial buffer from lfs.DecodePointer
	b := &bytes.Buffer{}
	r := io.TeeReader(os.Stdin, b)

	ptr, perr := lfs.DecodePointer(r)
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	if perr != nil {
		if _, err := io.Copy(to, pbuf); err != nil {
			return errors.Wrap(err, perr.Error())
		}

		return errors.NewNotAPointerError(errors.Errorf(
			"Unable to parse pointer at: %q", filename,
		))
	}

	lfs.LinkOrCopyFromReference(ptr.Oid, ptr.Size)
<<<<<<< HEAD
=======
	if smudgeInfo {
		// only invoked from `filter.lfs.smudge`, not `filter.lfs.process`
		// NOTE: this is deprecated behavior and will be removed in v2.0.0

		fmt.Fprintln(os.Stderr, "WARNING: 'smudge --info' is deprecated and will be removed in v2.0")
		fmt.Fprintln(os.Stderr, "USE INSTEAD:")
		fmt.Fprintln(os.Stderr, "  $ git lfs pointer --file=path/to/file")
		fmt.Fprintln(os.Stderr, "  $ git lfs ls-files")
		fmt.Fprintln(os.Stderr, "")

		localPath, err := lfs.LocalMediaPath(ptr.Oid)
		if err != nil {
			Exit(err.Error())
		}

<<<<<<< HEAD
		if stat, err := longpathos.Stat(localPath); err != nil {
=======
		if stat, err := os.Stat(localPath); err != nil {
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
			Print("%d --", ptr.Size)
		} else {
			Print("%d %s", stat.Size(), localPath)
		}

<<<<<<< HEAD
		return nil
	}

>>>>>>> refs/remotes/origin/release-1.5
	cb, file, err := lfs.CopyCallbackFile("smudge", filename, 1, 1)
	if err != nil {
		return err
	}

<<<<<<< HEAD
<<<<<<< HEAD
	download := !skip
	if download {
		download = filter.Allows(filename)
=======
	filter := filepathfilter.New(cfg.FetchIncludePaths(), cfg.FetchExcludePaths())
=======
>>>>>>> refs/remotes/origin/release-1.5
	download := filter.Allows(filename)
	if skip || cfg.Os.Bool("GIT_LFS_SKIP_SMUDGE", false) {
		download = false
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
	}

	err = ptr.Smudge(to, filename, download, getTransferManifest(), cb)
	if file != nil {
		file.Close()
	}

	if err != nil {
		ptr.Encode(to)
		// Download declined error is ok to skip if we weren't requesting download
		if !(errors.IsDownloadDeclinedError(err) && !download) {
			LoggedError(err, "Error downloading object: %s (%s)", filename, ptr.Oid)
			if !cfg.SkipDownloadErrors() {
				os.Exit(2)
			}
		}
	}

	return nil
=======
		return
	}

	if err := smudge(os.Stdout, ptr, smudgeFilename(args, perr), smudgeSkip); err != nil {
		Error(err.Error())
	}
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
}

func smudgeCommand(cmd *cobra.Command, args []string) {
	requireStdin("This command should be run by the Git 'smudge' filter")
	lfs.InstallHooks(false)

	if !smudgeSkip && cfg.Os.Bool("GIT_LFS_SKIP_SMUDGE", false) {
		smudgeSkip = true
	}
	filter := filepathfilter.New(cfg.FetchIncludePaths(), cfg.FetchExcludePaths())

<<<<<<< HEAD
<<<<<<< HEAD
	if err := smudge(os.Stdout, os.Stdin, smudgeFilename(args), smudgeSkip, filter); err != nil {
		if errors.IsNotAPointerError(err) {
			fmt.Fprintln(os.Stderr, err.Error())
=======
	lfs.LinkOrCopyFromReference(ptr.Oid, ptr.Size)

	if smudgeInfo {
		fmt.Fprintln(os.Stderr, "WARNING: 'smudge --info' is deprecated and will be removed in v2.0")
		fmt.Fprintln(os.Stderr, "USE INSTEAD:")
		fmt.Fprintln(os.Stderr, "  $ git lfs pointer --file=path/to/file")
		fmt.Fprintln(os.Stderr, "  $ git lfs ls-files")
		fmt.Fprintln(os.Stderr, "")

		localPath, err := lfs.LocalMediaPath(ptr.Oid)
		if err != nil {
			Exit(err.Error())
		}

		if stat, err := longpathos.Stat(localPath); err != nil {
			Print("%d --", ptr.Size)
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
=======
	if err := smudge(os.Stdout, os.Stdin, smudgeFilename(args), smudgeSkip, filter); err != nil {
		if errors.IsNotAPointerError(err) {
			fmt.Fprintln(os.Stderr, err.Error())
>>>>>>> refs/remotes/origin/release-1.5
		} else {
			Error(err.Error())
		}
	}
}

func smudgeFilename(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "<unknown file>"
}

func init() {
	RegisterCommand("smudge", smudgeCommand, func(cmd *cobra.Command) {
		cmd.Flags().BoolVarP(&smudgeSkip, "skip", "s", false, "")
	})
}
