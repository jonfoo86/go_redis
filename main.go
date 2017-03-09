package main

import (
	"fmt"
	"net"
)

func PrintMsg(messages chan string) {
	for {
		msg := <-messages
		fmt.Print(msg)
	}

}

func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

func sendRes(reschan chan *res)  {
	for{
		response := <- reschan
		response.conn.Write(response.buf)
	}

}

func Handler(conn net.Conn, messages chan string, responsechan chan *res) {

	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())

	recvpack := new(socketPack)
	remainpack := new(socketPack)
	for {
		lenght, err := conn.Read(recvpack.buf[:])
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght <= 0 {
			continue
		}
		recvpack.lenght = lenght
		//fmt.Println(recvpack.lenght, " ", string(recvpack.buf[0:recvpack.lenght]))

		result, cmdlist := ParseCmd(remainpack, recvpack)
		switch result {
		case RS_Ok:
			//fmt.Println("cmdlist:", cmdlist)
			request := new(req)
			request.cmdlist = cmdlist
			request.conn = conn
			responsechan <-ProcessReq(request)
			continue
		case RS_Fail:
			continue
		case RS_Error:
			conn.Close()
			break
		}
	}

}

func main() {
	ln, err := net.Listen("tcp", ":7000")
	if err != nil {
		// handle error
	}

	messages := make(chan string)
	responsechan := make(chan *res, 10000)
	go sendRes(responsechan)
	go PrintMsg(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go Handler(conn, messages, responsechan)
	}
}
