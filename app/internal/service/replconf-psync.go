package service

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
	"github.com/codecrafters-io/redis-starter-go/app/internal/encoder"
	"github.com/google/uuid"
)

func (h *MainNode) handleReplConf(conn net.Conn, _ internal.Request, _ MainNodeOptions) {
	_, err := conn.Write([]byte(encoder.NewSimpleString("OK")))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}

func (h *MainNode) handlePsync(conn net.Conn, _ internal.Request, handlerOpts MainNodeOptions) {
	// generate a unique conn_id for the connection
	conn_id := encoder.Sha1Hash(uuid.New().String())

	full_resync_resp := encoder.NewSimpleString(
		fmt.Sprintf("FULLRESYNC %v 0", conn_id),
	)
	_, err := conn.Write([]byte(full_resync_resp))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())

		conn.Close()
		os.Exit(1)
	}

	rdb_hex := []byte("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	rdb_bin, err := decoder.DecodeHexToBinary(rdb_hex)
	if err != nil {
		log.Println("Error decoding hex to binary: ", err.Error())

		conn.Close()
		os.Exit(1)
	}

	// Send RDB file to replica
	_, err = conn.Write([]byte(encoder.NewBinaryString(rdb_bin)))
	if err != nil {
		log.Println("Error writing to connection: ", err.Error())

		conn.Close()
		os.Exit(1)
	}

	// if connection is not closed, append to connection pool
	handlerOpts.ShouldClose = false
	h.AddToConnPool(conn, conn_id)
}