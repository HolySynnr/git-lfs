package git

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterProcessScannerInitializesWithCorrectSupportedValues(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	if err := pl.writePacketText("git-filter-client"); err != nil {
		t.Fatalf("expected... %v", err.Error())
	}

	require.Nil(t, pl.writePacketText("git-filter-client"))
	require.Nil(t, pl.writePacketList([]string{"version=2"}))

	fps := NewFilterProcessScanner(&from, &to)
	err := fps.Init()

	assert.Nil(t, err)

	out, err := newPktline(&to, nil).readPacketList()
	assert.Nil(t, err)
	assert.Equal(t, []string{"git-filter-server", "version=2"}, out)
}

func TestFilterProcessScannerRejectsUnrecognizedInitializationMessages(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	require.Nil(t, pl.writePacketText("git-filter-client-unknown"))
<<<<<<< HEAD
	require.Nil(t, pl.writeFlush())
=======
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased

	fps := NewFilterProcessScanner(&from, &to)
	err := fps.Init()

	require.NotNil(t, err)
<<<<<<< HEAD
	assert.Equal(t, "invalid filter-process pkt-line welcome message: git-filter-client-unknown", err.Error())
=======
	assert.Equal(t, "invalid filter pkt-line welcome message: git-filter-client-unknown", err.Error())
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	assert.Empty(t, to.Bytes())
}

func TestFilterProcessScannerRejectsUnsupportedFilters(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	require.Nil(t, pl.writePacketText("git-filter-client"))
	// Write an unsupported version
	require.Nil(t, pl.writePacketList([]string{"version=0"}))

	fps := NewFilterProcessScanner(&from, &to)
	err := fps.Init()

	require.NotNil(t, err)
	assert.Equal(t, "filter 'version=2' not supported (your Git supports: [version=0])", err.Error())
	assert.Empty(t, to.Bytes())
}

func TestFilterProcessScannerNegotitatesSupportedCapabilities(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	require.Nil(t, pl.writePacketList([]string{
		"capability=clean", "capability=smudge", "capability=not-invented-yet",
	}))

	fps := NewFilterProcessScanner(&from, &to)
	err := fps.NegotiateCapabilities()

	assert.Nil(t, err)

	out, err := newPktline(&to, nil).readPacketList()
	assert.Nil(t, err)
	assert.Equal(t, []string{"capability=clean", "capability=smudge"}, out)
}

func TestFilterProcessScannerDoesNotNegotitatesUnsupportedCapabilities(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	// Write an unsupported capability
	require.Nil(t, pl.writePacketList([]string{
		"capability=unsupported",
	}))

	fps := NewFilterProcessScanner(&from, &to)
	err := fps.NegotiateCapabilities()

	require.NotNil(t, err)
	assert.Equal(t, "filter 'capability=clean' not supported (your Git supports: [capability=unsupported])", err.Error())
	assert.Empty(t, to.Bytes())
}

func TestFilterProcessScannerReadsRequestHeadersAndPayload(t *testing.T) {
	var from, to bytes.Buffer

	pl := newPktline(nil, &from)
	// Headers
	require.Nil(t, pl.writePacketList([]string{
<<<<<<< HEAD
		"foo=bar", "other=woot", "crazy='sq',\\$x=.bin",
=======
		"foo=bar", "other=woot",
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
	}))
	// Multi-line packet
	require.Nil(t, pl.writePacketText("first"))
	require.Nil(t, pl.writePacketText("second"))
<<<<<<< HEAD
	require.Nil(t, pl.writeFlush())
=======
	_, err := from.Write([]byte{0x30, 0x30, 0x30, 0x30}) // flush packet
	assert.Nil(t, err)
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased

	req, err := readRequest(NewFilterProcessScanner(&from, &to))

	assert.Nil(t, err)
	assert.Equal(t, req.Header["foo"], "bar")
	assert.Equal(t, req.Header["other"], "woot")
<<<<<<< HEAD
	assert.Equal(t, req.Header["crazy"], "'sq',\\$x=.bin")
=======
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased

	payload, err := ioutil.ReadAll(req.Payload)
	assert.Nil(t, err)
	assert.Equal(t, []byte("first\nsecond\n"), payload)
}

func TestFilterProcessScannerRejectsInvalidHeaderPackets(t *testing.T) {
<<<<<<< HEAD
	from := bytes.NewBuffer([]byte{
		0x30, 0x30, 0x30, 0x34, // 0004 (invalid packet length)
	})

	req, err := readRequest(NewFilterProcessScanner(from, nil))
=======
	var from bytes.Buffer

	pl := newPktline(nil, &from)
	// (Invalid) headers
	require.Nil(t, pl.writePacket([]byte{}))

	req, err := readRequest(NewFilterProcessScanner(&from, nil))
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased

	require.NotNil(t, err)
	assert.Equal(t, "Invalid packet length.", err.Error())

	assert.Nil(t, req)
}

<<<<<<< HEAD
func TestFilterProcessScannerWritesStatusPackets(t *testing.T) {
	var buf bytes.Buffer

	require.Nil(t, NewFilterProcessScanner(nil, &buf).WriteStatus("success"))

	list, err := newPktline(&buf, nil).readPacketList()

	assert.Nil(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "status=success", list[0])
}

func TestFilterProcessScannerAbbreviatesUnchangedStatuses(t *testing.T) {
	var buf bytes.Buffer

	s := NewFilterProcessScanner(nil, &buf)

	for _, status := range []string{
		"success", "success", "error",
	} {
		require.Nil(t, s.WriteStatus(status))
	}

	pl := newPktline(&buf, nil)

	// Read the first status=success entirely
	assertPacketRead(t, pl, []byte("status=success\n"))
	assertPacketRead(t, pl, nil)

	// Second status is the same, so expect only a flush packet
	assertPacketRead(t, pl, nil)

	// Third status is different than the previous, so expect the full
	// packet
	assertPacketRead(t, pl, []byte("status=error\n"))
	assertPacketRead(t, pl, nil)
}

// readRequest performs a single scan operation on the given
=======
// readRequest preforms a single scan operation on the given
>>>>>>> refs/remotes/git-lfs/filter-stream-rebased
// `*FilterProcessScanner`, "s", and returns: an error if there was one, or a
// request if there was one.  If neither, it returns (nil, nil).
func readRequest(s *FilterProcessScanner) (*Request, error) {
	s.Scan()

	if err := s.Err(); err != nil {
		return nil, err
	}

	return s.Request(), nil
}
