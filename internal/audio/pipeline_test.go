package audio

import (
	"context"
	"errors"
	"io"
	"testing"
)

// fakeStreamer 是測試用的 Streamer 實作。
type fakeStreamer struct {
	reader io.ReadCloser
	err    error
}

func (s *fakeStreamer) Stream(ctx context.Context, url string) (io.ReadCloser, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.reader, nil
}

// fakeVoiceConnection 是測試用的 VoiceConnection 實作。
type fakeVoiceConnection struct {
	speakingCalled bool
	speakingState  bool
	opusChan       chan []byte
	speakingErr    error
}

func newFakeVoiceConnection() *fakeVoiceConnection {
	return &fakeVoiceConnection{
		opusChan: make(chan []byte, 10),
	}
}

func (vc *fakeVoiceConnection) Speaking(speaking bool) error {
	vc.speakingCalled = true
	vc.speakingState = speaking
	return vc.speakingErr
}

func (vc *fakeVoiceConnection) OpusSend() chan<- []byte {
	return vc.opusChan
}

func (vc *fakeVoiceConnection) Disconnect() error {
	return nil
}

// fakeReader 是測試用的 io.ReadCloser。
type fakeReader struct {
	data   []byte
	pos    int
	closed bool
}

func newFakeReader(data []byte) *fakeReader {
	return &fakeReader{data: data}
}

func (r *fakeReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *fakeReader) Close() error {
	r.closed = true
	return nil
}

// TestPipeline_Play_SetsAndUnsetsSpeaking 測試 Play 會設定和取消說話狀態。
func TestPipeline_Play_SetsAndUnsetsSpeaking(t *testing.T) {
	fakeData := []byte{0x01, 0x02, 0x03, 0x04}
	reader := newFakeReader(fakeData)
	streamer := &fakeStreamer{reader: reader}
	vc := newFakeVoiceConnection()

	pipeline := NewPipeline(streamer)

	err := pipeline.Play(context.Background(), "test-url", vc)
	if err != nil {
		t.Fatalf("Play 不應回傳錯誤: %v", err)
	}

	if !vc.speakingCalled {
		t.Error("Play 應該呼叫 Speaking()")
	}
}

// TestPipeline_Play_ReturnsStreamerError 測試 streamer 失敗時回傳錯誤。
func TestPipeline_Play_ReturnsStreamerError(t *testing.T) {
	expectedErr := errors.New("stream failed")
	streamer := &fakeStreamer{err: expectedErr}
	vc := newFakeVoiceConnection()

	pipeline := NewPipeline(streamer)

	err := pipeline.Play(context.Background(), "test-url", vc)
	if err == nil {
		t.Fatal("streamer 失敗時應該回傳錯誤")
	}
}

// TestPipeline_Play_RespectsContext 測試 Play 會尊重 context 取消。
func TestPipeline_Play_RespectsContext(t *testing.T) {
	// 無限 reader（不會回傳 EOF）
	infiniteReader := &infiniteReader{}
	streamer := &fakeStreamer{reader: infiniteReader}
	vc := newFakeVoiceConnection()

	pipeline := NewPipeline(streamer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	err := pipeline.Play(ctx, "test-url", vc)
	if err == nil {
		t.Fatal("context 取消時應該回傳錯誤")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("錯誤應為 context.Canceled，實際為: %v", err)
	}
}

// infiniteReader 永遠不會回傳 EOF。
type infiniteReader struct {
	closed bool
}

func (r *infiniteReader) Read(p []byte) (int, error) {
	// 填滿 buffer 但不結束
	for i := range p {
		p[i] = 0xFF
	}
	return len(p), nil
}

func (r *infiniteReader) Close() error {
	r.closed = true
	return nil
}

// TestPipeline_Play_ClosesStream 測試 Play 完成後會關閉 stream。
func TestPipeline_Play_ClosesStream(t *testing.T) {
	fakeData := []byte{0x01}
	reader := newFakeReader(fakeData)
	streamer := &fakeStreamer{reader: reader}
	vc := newFakeVoiceConnection()

	pipeline := NewPipeline(streamer)

	_ = pipeline.Play(context.Background(), "test-url", vc)

	if !reader.closed {
		t.Error("Play 完成後應該關閉 stream")
	}
}
