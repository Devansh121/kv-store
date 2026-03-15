package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Devansh121/kv-store/config"
	"github.com/Devansh121/kv-store/core"
)

func readCommand(conn net.Conn) (*core.RedisCMD, error) {
	// TODO: Max read in one shot is 512 bytes.
	// To allow input > 512 bytes, repeated reads are needed until
	// EOF or the designated delimiter is reached.
	buffer := make([]byte, 512)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buffer[:n])
	if err != nil {
		return nil, err
	}
	return &core.RedisCMD{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCMD, conn net.Conn) {
	err := core.EvalAndRespond(cmd, conn)
	if err != nil {
		respondError(err, conn)
	}
}

func RunSyncTCPServer() {
	log.Println("Starting sync TCP server on", config.Host, config.Port)

	concurrentClients := 0

	// Listen on the configured host:port.
	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		// Blocking call waiting for a new client connection.
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		concurrentClients++
		log.Println("Client connected:", conn.RemoteAddr(), "concurrent clients", concurrentClients)

		for {
			command, err := readCommand(conn)
			if err != nil {
				conn.Close()
				concurrentClients--
				log.Println("Client disconnected:", conn.RemoteAddr(), "concurrent clients", concurrentClients)
				if err == io.EOF {
					break
				}
				log.Println("Read error:", err)
			}

			respond(command, conn)
		}
	}
}
