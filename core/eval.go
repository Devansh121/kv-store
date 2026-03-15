package core

import (
	"errors"
	"net"
)

func evalPING(args []string, conn net.Conn) error {
	var response []byte

	if len(args) > 1 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		response = Encode("PONG", true)
	} else {
		response = Encode(args[0], false)
	}

	_, err := conn.Write(response)
	return err
}

func EvalAndRespond(cmd *RedisCMD, conn net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, conn)
	default:
		return errors.New("ERR unknown command '" + cmd.Cmd + "'")
	}
}
