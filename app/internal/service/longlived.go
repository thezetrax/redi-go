package service

import (
	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
)

// IsLongLived checks if the connection is long-lived.
// A long-lived connection is a connection used for replication.
func IsLongLived(buf []byte, read int) bool {
	req := internal.ParseRequest(buf[:read])

	switch req.CMD.CMD {
	case decoder.CMD_PSYNC:
		return true
	default:
		return false
	}
}

// IsDelegateReq checks if the connection is a delegate request.
// A delegate request should be delegated to other replica nodes.
func IsDelegateReq(cmd decoder.CMD) bool {
	switch cmd.CMD {
	case decoder.CMD_SET:
		return true
	default:
		return false
	}
}
