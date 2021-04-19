package geerpc

import (
	"net"
	"testing"
)

func TestStartServer(t *testing.T) {
	lis, _ := net.Listen("tcp", ":9999")
	defer lis.Close()
	Accept(lis)
}

