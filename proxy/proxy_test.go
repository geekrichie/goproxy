package proxy

import "testing"

func TestTcpProxy(t *testing.T) {
	TcpProxy(":9999")
}