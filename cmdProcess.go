package main

import (
	"net"
)

type req struct {
	conn    net.Conn
	cmdlist []string
}

type res struct {
	conn net.Conn
	buf  []byte
}

func getProcess(cmdlist []string) []byte {
	return okResponse("get ok ")

}

func setProcess(cmdlist []string) []byte {
	return okResponse("set ok ")

}

func okResponse(str string) []byte {
	buf := make([]byte, len(str)+3)
	copy(buf[:], "+")
	copy(buf[1:], str)
	copy(buf[1+len(str):], "\r\n")
	return buf
}

func errorProcess(str string) []byte {
	buf := make([]byte, len(str)+3)
	//command := "*3\r\n$3\r\nSET\r\n$5\r\nHENRY\r\n$8\r\nHENRYFAN\r\n"
	copy(buf[:], "-")
	copy(buf[1:], str)
	copy(buf[1+len(str):], "\r\n")
	return buf
}

func cmdProcess(request *req) *res {
	response := new(res)
	response.conn = request.conn
	if len(request.cmdlist) == 0 {
		return response
	}

	switch request.cmdlist[0] {
	case "get":
		response.buf = getProcess(request.cmdlist)
		break
	case "set":
		response.buf = setProcess(request.cmdlist)
		break
	default:
		response.buf = errorProcess(request.cmdlist[0])
	}

	return response
}
