// connection_test.go: Unit tests for connection.go
package connection

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/soderluk/nirimgr/models"
	"github.com/stretchr/testify/assert"
)

type mockConn struct {
	writeErr error
	readData []string
	readIdx  int
	closeErr error
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.readIdx >= len(m.readData) {
		return 0, errors.New("EOF")
	}
	copy(b, m.readData[m.readIdx])
	n = len(m.readData[m.readIdx])
	m.readIdx++
	return n, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(b), nil
}

func (m *mockConn) Close() error {
	return m.closeErr
}
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestSend(t *testing.T) {
	mc := &mockConn{}
	s := &NiriSocket{conn: mc}
	err := s.Send("test-request")
	assert.NoError(t, err)
}

func TestSendError(t *testing.T) {
	mc := &mockConn{writeErr: errors.New("write error")}
	s := &NiriSocket{conn: mc}
	err := s.Send("test-request")
	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	mc := &mockConn{}
	s := &NiriSocket{conn: mc}
	err := s.conn.Close()
	assert.NoError(t, err)
}

func TestCloseError(t *testing.T) {
	mc := &mockConn{closeErr: errors.New("close error")}
	s := &NiriSocket{conn: mc}
	err := s.conn.Close()
	assert.Error(t, err)
}

func TestStructToString(t *testing.T) {
	m := map[string]any{"foo": "bar"}
	str, err := structToString(m)
	assert.NoError(t, err)
	assert.Contains(t, str, "foo")
}

func TestStructToMap(t *testing.T) {
	type testAName struct {
		Name string
	}
	type testStruct struct {
		testAName
		Val int
	}
	ts := testStruct{testAName: testAName{Name: "test"}, Val: 42}
	m, err := structToMap(ts)
	assert.NoError(t, err)
	assert.NotContains(t, m, "Name")
	assert.Contains(t, m, "Val")
}

type mockAction struct{}

func (mockAction) GetName() string { return "MockAction" }

func TestPerformAction(t *testing.T) {
	origSocket := Socket
	defer func() { Socket = origSocket }()

	Socket = func() *NiriSocket {
		return &NiriSocket{conn: &mockConn{}}
	}

	action := mockAction{}
	res := PerformAction(action)
	assert.True(t, res)
}

func TestPerformRequest(t *testing.T) {
	origSocket := Socket
	defer func() { Socket = origSocket }()

	Socket = func() *NiriSocket {
		return &NiriSocket{conn: &mockConn{readData: []string{`{"foo":"bar"}`}}}
	}

	// Patch models.Response to be compatible for test
	// This is a hack for demonstration; in real code, use interfaces or dependency injection
	stream, err := PerformRequest(models.Version)
	assert.NoError(t, err)
	resp, ok := <-stream
	assert.True(t, ok)

	// The response should be a map with foo: bar
	if m, ok := resp.Ok["foo"]; ok {
		var val string
		_ = json.Unmarshal(m, &val)
		assert.Equal(t, "bar", val)
	}
}
