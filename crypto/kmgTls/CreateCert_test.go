package kmgTls

import (
	"bytes"
	"crypto/tls"
	"github.com/bronze1man/kmg/kmgTest"
	"net"
	"testing"
	"time"
)

type mockAddr struct {
}

func (*mockAddr) Network() string {
	return "MockAddr.Network"
}
func (*mockAddr) String() string {
	return "MockAddr.String"
}

type mockListener struct {
}

func (*mockListener) Accept() (c net.Conn, err error) {
	return &mockConn{}, nil
}

func (*mockListener) Close() error {
	return nil
}

func (*mockListener) Addr() net.Addr {
	return &mockAddr{}
}

type mockConn struct {
	bytes.Buffer
}

func (*mockConn) Close() error {
	return nil
}
func (*mockConn) LocalAddr() net.Addr {
	return &mockAddr{}
}

func (*mockConn) RemoteAddr() net.Addr {
	return &mockAddr{}
}

func (*mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (*mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (*mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
func TestCreateCert(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	config, err := CreateTlsConfig()
	t.Equal(err, nil)
	t.Ok(config != nil)

	_ = tls.NewListener(&mockListener{}, config)
}
