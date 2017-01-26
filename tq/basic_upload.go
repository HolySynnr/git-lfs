package tq

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/git-lfs/git-lfs/errors"
	"github.com/git-lfs/git-lfs/lfsapi"
	"github.com/git-lfs/git-lfs/progress"
	"github.com/git-lfs/git-lfs/tools/longpathos"
)

const (
	BasicAdapterName = "basic"
)

// Adapter for basic uploads (non resumable)
type basicUploadAdapter struct {
	*adapterBase
}

func (a *basicUploadAdapter) ClearTempStorage() error {
	// Should be empty already but also remove dir
	return longpathos.RemoveAll(a.tempDir())
}

func (a *basicUploadAdapter) tempDir() string {
	// Must be dedicated to this adapter as deleted by ClearTempStorage
	d := filepath.Join(os.TempDir(), "git-lfs-basic-temp")
	if err := longpathos.MkdirAll(d, 0755); err != nil {
		return os.TempDir()
	}
	return d
}

func (a *basicUploadAdapter) WorkerStarting(workerNum int) (interface{}, error) {
	return nil, nil
}
func (a *basicUploadAdapter) WorkerEnding(workerNum int, ctx interface{}) {
}

func (a *basicUploadAdapter) DoTransfer(ctx interface{}, t *Transfer, cb ProgressCallback, authOkFunc func()) error {
	rel, err := t.Actions.Get("upload")
	if err != nil {
		return err
		// return fmt.Errorf("No upload action for this object.")
	}

	req, err := a.newHTTPRequest("PUT", rel)
	if err != nil {
		return err
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	if req.Header.Get("Transfer-Encoding") == "chunked" {
		req.TransferEncoding = []string{"chunked"}
	} else {
		req.Header.Set("Content-Length", strconv.FormatInt(t.Size, 10))
	}

	req.ContentLength = t.Size

	f, err := longpathos.OpenFile(t.Path, os.O_RDONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "basic upload")
	}
	defer f.Close()

	// Ensure progress callbacks made while uploading
	// Wrap callback to give name context
	ccb := func(totalSize int64, readSoFar int64, readSinceLast int) error {
		if cb != nil {
			return cb(t.Name, totalSize, readSoFar, readSinceLast)
		}
		return nil
	}
	var reader lfsapi.ReadSeekCloser = progress.NewBodyWithCallback(f, t.Size, ccb)

	// Signal auth was ok on first read; this frees up other workers to start
	if authOkFunc != nil {
		reader = newStartCallbackReader(reader, func() error {
			authOkFunc()
			return nil
		})
	}

	req.Body = reader

	res, err := a.doHTTP(t, req)
	if err != nil {
		return errors.NewRetriableError(err)
	}

	a.apiClient.LogResponse("lfs.data.upload", res)

	// A status code of 403 likely means that an authentication token for the
	// upload has expired. This can be safely retried.
	if res.StatusCode == 403 {
		err = errors.New("http: received status 403")
		return errors.NewRetriableError(err)
	}

	if res.StatusCode > 299 {
		return errors.Wrapf(nil, "Invalid status for %s %s: %d",
			req.Method,
			strings.SplitN(req.URL.String(), "?", 2)[0],
			res.StatusCode,
		)
	}

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	return verifyUpload(a.apiClient, t)
}

// startCallbackReader is a reader wrapper which calls a function as soon as the
// first Read() call is made. This callback is only made once
type startCallbackReader struct {
	cb     func() error
	cbDone bool
	lfsapi.ReadSeekCloser
}

func (s *startCallbackReader) Read(p []byte) (n int, err error) {
	if !s.cbDone && s.cb != nil {
		if err := s.cb(); err != nil {
			return 0, err
		}
		s.cbDone = true
	}
	return s.ReadSeekCloser.Read(p)
}
func newStartCallbackReader(r lfsapi.ReadSeekCloser, cb func() error) *startCallbackReader {
	return &startCallbackReader{
		ReadSeekCloser: r,
		cb:             cb,
	}
}

func configureBasicUploadAdapter(m *Manifest) {
	m.RegisterNewAdapterFunc(BasicAdapterName, Upload, func(name string, dir Direction) Adapter {
		switch dir {
		case Upload:
			bu := &basicUploadAdapter{newAdapterBase(name, dir, nil)}
			// self implements impl
			bu.transferImpl = bu
			return bu
		case Download:
			panic("Should never ask this func for basic download")
		}
		return nil
	})
}
