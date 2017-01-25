package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type byteLogger struct {
	*bytes.Buffer
}

func (l *byteLogger) Close() error {
	return nil
}

func (l *byteLogger) Sync() error {
	return nil
}

func withByteBuffer(buf *bytes.Buffer) meterOption {
	return func(m *ProgressMeter) {
		m.logger.writeData = true
		m.logger.log = &byteLogger{Buffer: buf}
		m.out = buf
	}
}

func TestMeter(t *testing.T) {
	b := &bytes.Buffer{}
	m := NewMeter(withByteBuffer(b))

	time.Sleep(time.Millisecond * 201)

	assert.Equal(t, "", b.String())

	m.Add(int64(50))
	m.Add(int64(150))
	m.Skip(int64(50))
	assert.True(t, m.Start())
	assert.False(t, m.Start())
	time.Sleep(time.Millisecond * 201)
	assert.True(t, m.Pause())
	assert.False(t, m.Pause())

	expected := []string{
		"Git LFS: (0 of 1 files, 1 skipped) 0 B / 150 B, 50 B skipped",
		"Git LFS: (0 of 1 files, 1 skipped) 0 B / 150 B, 50 B skipped",
	}
	assert.Equal(t, expected, messages(b))

	assert.True(t, m.Start())
	assert.False(t, m.Start())
	m.Add(int64(10))
	time.Sleep(time.Millisecond * 201)
	assert.True(t, m.Pause())
	assert.False(t, m.Pause())

	msgs := messages(b)
	assert.Equal(t, "Git LFS: (0 of 2 files, 1 skipped) 0 B / 160 B, 50 B skipped",
		msgs[len(msgs)-1])

	assert.True(t, m.Finish())
	assert.False(t, m.Finish())
}

func messages(b *bytes.Buffer) []string {
	msgs := strings.Split(strings.TrimSpace(b.String()), "\r")
	out := make([]string, len(msgs))
	for i, m := range msgs {
		out[i] = strings.TrimSpace(m)
	}
	return out
}
