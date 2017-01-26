package progress

func Noop() Meter {
	return &nonMeter{}
}

type nonMeter struct{}

func (m *nonMeter) Start() bool                                                          { return false }
func (m *nonMeter) Pause() bool                                                          { return false }
func (m *nonMeter) Add(size int64)                                                       {}
func (m *nonMeter) Skip(size int64)                                                      {}
func (m *nonMeter) StartTransfer(name string)                                            {}
func (m *nonMeter) TransferBytes(direction, name string, read, total int64, current int) {}
func (m *nonMeter) FinishTransfer(name string)                                           {}
func (m *nonMeter) Finish() bool                                                         { return false }
