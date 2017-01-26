package commands

import (
<<<<<<< HEAD
=======
	"bytes"
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	"fmt"
	"io"
	"os"

<<<<<<< HEAD
	"github.com/git-lfs/git-lfs/errors"
	"github.com/git-lfs/git-lfs/filepathfilter"
	"github.com/git-lfs/git-lfs/git"
	"github.com/git-lfs/git-lfs/lfs"
=======
	"github.com/github/git-lfs/git"
	"github.com/github/git-lfs/lfs"
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	"github.com/spf13/cobra"
)

const (
	// cleanFilterBufferCapacity is the desired capacity of the
	// `*git.PacketWriter`'s internal buffer when the filter protocol
	// dictates the "clean" command. 512 bytes is (in most cases) enough to
	// hold an entire LFS pointer in memory.
	cleanFilterBufferCapacity = 512

	// smudgeFilterBufferCapacity is the desired capacity of the
	// `*git.PacketWriter`'s internal buffer when the filter protocol
	// dictates the "smudge" command.
	smudgeFilterBufferCapacity = git.MaxPacketLength
)

// filterSmudgeSkip is a command-line flag owned by the `filter-process` command
// dictating whether or not to skip the smudging process, leaving pointers as-is
// in the working tree.
var filterSmudgeSkip bool

<<<<<<< HEAD
=======
// filterSmudge is a gateway to the `smudge()` function and serves to bail out
// immediately if the pointer decoded from "from" has no data (i.e., is empty).
// This function, unlike the implementation found in the legacy smudge command,
// only combines the `io.Reader`s when necessary, since the implementation
// found in `*git.PacketReader` blocks while waiting for the following packet.
func filterSmudge(from io.Reader, to io.Writer, filename string) error {
	var pbuf bytes.Buffer
	from = io.TeeReader(from, &pbuf)

	ptr, err := lfs.DecodePointer(from)
	if err != nil {
		// If we tried to decode a pointer out of the data given to us,
		// and the file was _empty_, write out an empty file in
		// response. This occurs because when the clean filter
		// encounters an empty file, and writes out an empty file,
		// instead of a pointer.
		//
		// TODO(taylor): figure out if there is more data on the reader,
		// and buffer that as well.
		if len(pbuf.Bytes()) == 0 {
			if _, cerr := io.Copy(to, &pbuf); cerr != nil {
				Panic(cerr, "Error writing data to stdout:")
			}
			return nil
		}

		return err
	}

	lfs.LinkOrCopyFromReference(ptr.Oid, ptr.Size)

	return smudge(to, ptr, filename, filterSmudgeSkip)
}

>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
func filterCommand(cmd *cobra.Command, args []string) {
	requireStdin("This command should be run by the Git filter process")
	lfs.InstallHooks(false)

	s := git.NewFilterProcessScanner(os.Stdin, os.Stdout)

	if err := s.Init(); err != nil {
		ExitWithError(err)
	}
	if err := s.NegotiateCapabilities(); err != nil {
		ExitWithError(err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> refs/remotes/origin/release-1.5
	skip := filterSmudgeSkip || cfg.Os.Bool("GIT_LFS_SKIP_SMUDGE", false)
	filter := filepathfilter.New(cfg.FetchIncludePaths(), cfg.FetchExcludePaths())

	var malformed []string

<<<<<<< HEAD
=======
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
=======
>>>>>>> refs/remotes/origin/release-1.5
=======
Scan:
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	for s.Scan() {
		var err error
		var w *git.PktlineWriter

		req := s.Request()

<<<<<<< HEAD
		s.ForgetStatus()
		s.WriteStatus(statusFromErr(nil))
=======
		s.WriteStatus("success")
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased

		switch req.Header["command"] {
		case "clean":
			w = git.NewPktlineWriter(os.Stdout, cleanFilterBufferCapacity)
<<<<<<< HEAD
			err = clean(w, req.Payload, req.Header["pathname"])
		case "smudge":
			w = git.NewPktlineWriter(os.Stdout, smudgeFilterBufferCapacity)
			err = smudge(w, req.Payload, req.Header["pathname"], skip, filter)
		default:
			ExitWithError(fmt.Errorf("Unknown command %q", req.Header["command"]))
<<<<<<< HEAD
		}

		if errors.IsNotAPointerError(err) {
			malformed = append(malformed, req.Header["pathname"])
			err = nil
=======
>>>>>>> refs/remotes/git-lfs/1.5/filepathfilter
		}

		if errors.IsNotAPointerError(err) {
			malformed = append(malformed, req.Header["pathname"])
			err = nil
=======
			err = clean(req.Payload, w, req.Header["pathname"])
		case "smudge":
			w = git.NewPktlineWriter(os.Stdout, smudgeFilterBufferCapacity)
			err = filterSmudge(req.Payload, w, req.Header["pathname"])
		default:
			fmt.Errorf("Unknown command %s", cmd)
			break Scan
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
		}

		var status string
		if ferr := w.Flush(); ferr != nil {
<<<<<<< HEAD
			status = statusFromErr(ferr)
=======
			status = "error"
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
		} else {
			status = statusFromErr(err)
		}

		s.WriteStatus(status)
	}

<<<<<<< HEAD
<<<<<<< HEAD
	if len(malformed) > 0 {
		fmt.Fprintf(os.Stderr, "Encountered %d file(s) that should have been pointers, but weren't:\n", len(malformed))
		for _, m := range malformed {
			fmt.Fprintf(os.Stderr, "\t%s\n", m)
		}
	}
=======
	// TODO: Detect an EOF after a successful filter-process request (EOF at
	// any other point in the protocol would be an error) and wait for
	// downloaded files to finish. Afterwards copy all downloaded files to
	// their final location in the work tree.
	lfs.WaitForDownloads(TransferManifest())
>>>>>>> refs/remotes/git-lfs/promised-downloads

=======
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	if err := s.Err(); err != nil && err != io.EOF {
		ExitWithError(err)
	}
}

// statusFromErr returns the status code that should be sent over the filter
// protocol based on a given error, "err".
func statusFromErr(err error) string {
	if err != nil && err != io.EOF {
		return "error"
	}
	return "success"
}

func init() {
	RegisterCommand("filter-process", filterCommand, func(cmd *cobra.Command) {
		cmd.Flags().BoolVarP(&filterSmudgeSkip, "skip", "s", false, "")
	})
}
