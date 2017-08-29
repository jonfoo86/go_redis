package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"runtime"
)

func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

func responseHandler(reschan chan *res) {
	for {
		response := <-reschan
		response.conn.Write(response.buf)
		DelRes(response)
	}

}

func requestHandler(conn net.Conn, requestchan chan *req) {

	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())

	recvpack := new(socketPack)
	cachepack := new(socketPack)
	for {
		lenght, err := conn.Read(recvpack.buf[:])
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght <= 0 {
			continue
		}
		recvpack.length = lenght
		//fmt.Println(recvpack.length, " ", string(recvpack.buf[0:recvpack.length]))

		result, cmdlist := cmdParse(cachepack, recvpack)
		switch result {
		case RS_Ok:
			//fmt.Println("cmdlist:", cmdlist)
			request := new(req)
			request.cmdlist = cmdlist
			request.conn = conn
			requestchan <- request

			continue
		case RS_Fail:
			continue
		case RS_Error:
			conn.Close()
			break
		}
	}

}

func cmdHandler(responsechan chan *res, requestchan chan *req) {
	for {
		request := <-requestchan
		responsechan <- cmdProcess(request)
	}
}

func acceptHandler() {
	ln, err := net.Listen("tcp", ":7000")
	if err != nil {
		// handle error
	}

	responsechan := make(chan *res, 10000)
	requestchan := make(chan *req, 10000)
	go responseHandler(responsechan)
	go cmdHandler(responsechan, requestchan)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go requestHandler(conn, requestchan)
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	runtime.GOMAXPROCS(4) // 最多使用2个核
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go acceptHandler()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Println("pid: ", os.Getpid())
	fmt.Println("Wartint signal !")
	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s, " stop ")
}
