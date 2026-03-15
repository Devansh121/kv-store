package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/Devansh121/kv-store/config"
)

func readCommand(conn net.Conn) (string, error) {
	// TODO: Max read in one shot is 512 bytes.
	// To allow input > 512 bytes, repeated reads are needed until
	// EOF or the designated delimiter is reached.
	buffer := make([]byte, 512)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func respond(response string, conn net.Conn) error {
	if _, err := conn.Write([]byte(response)); err != nil {
		return err
	}
	return nil
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

			log.Println("Command:", command)
			if err = respond(command, conn); err != nil {
				log.Print("Write error:", err)
			}
		}
	}
}
