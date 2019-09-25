package aria2shell

import "testing"

func TestNewRPC(t *testing.T) {
	NewRPC("http://aria2rpc.y11e.com", "5000", "jsonrpc")
}
