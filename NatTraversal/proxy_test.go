package NatTraversal

import "testing"

func TestRemoteOperation(t *testing.T) {
	RemoteOperation()
}

func TestLocalOperation(t *testing.T) {
	go Serve()
	LocalOperation()
}
