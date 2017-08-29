package main

import (
	"container/list"
	//"fmt"
	"net"
	"regexp"
	"strconv"
	"sync"
)

type req struct {
	conn    net.Conn
	cmdlist []string
}

type res struct {
	conn net.Conn
	buf  []byte
}

var (
	cachemap = make(map[string]*string)
	pipe     = &sync.Pool{New: func() interface{} { return new(res) }}
)

func NewRes() *res {
	return pipe.Get().(*res)
}

func DelRes(ptr *res) {
	pipe.Put(ptr)
}

func getProcess(cmdlist []string) []byte {
	if v, ok := cachemap[cmdlist[1]]; ok {
		return okResponse(*v)
	}
	return okResponse("Key Not Found")
}

func setProcess(cmdlist []string) []byte {
	cachemap[cmdlist[1]] = &(cmdlist[2])
	return okResponse("set ok ")

}

func keysProcess(cmdlist []string) []byte {
	matchlist := list.New()
	var validID = regexp.MustCompile(cmdlist[1])
	for k, _ := range cachemap {
		if validID.MatchString(k) {
			matchlist.PushBack(k)
		}
		if matchlist.Len() > 100 {
			break
		}
	}
	return strlistResponse(matchlist)
}

func okResponse(str string) []byte {
	buf := make([]byte, len(str)+3)
	copy(buf[:], "+")
	copy(buf[1:], str)
	copy(buf[1+len(str):], "\r\n")
	return buf
}

func getByteNum(num int) int {
	var bytenum int = 1
	for {
		num = num / 10
		if num < 1 {
			break
		}
		bytenum++
	}
	return bytenum
}

func strlistResponse(l *list.List) []byte {
	headerlen := getByteNum(l.Len()) + 3
	for e := l.Front(); e != nil; e = e.Next() {
		headerlen += getByteNum(len(e.Value.(string))) + 3 + len(e.Value.(string)) + 2
	}

	ret := make([]byte, headerlen)
	copy(ret[:], "*")
	copy(ret[1:], strconv.Itoa(l.Len()))
	pos := 1 + getByteNum(l.Len())
	copy(ret[pos:], "\r\n")
	pos += 2

	for e := l.Front(); e != nil; e = e.Next() {
		copy(ret[pos:], "$")
		pos += 1
		copy(ret[pos:], strconv.Itoa(len(e.Value.(string))))
		pos += getByteNum(len(e.Value.(string)))
		copy(ret[pos:], "\r\n")
		pos += 2

		copy(ret[pos:], e.Value.(string))
		pos += len(e.Value.(string))
		copy(ret[pos:], "\r\n")
		pos += 2
	}

	return ret
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
	response := NewRes()
	response.conn = request.conn
	if len(request.cmdlist) == 0 {
		return response
	}

	switch request.cmdlist[0] {
	case "get", "GET":
		response.buf = getProcess(request.cmdlist)
	case "set", "SET":
		response.buf = setProcess(request.cmdlist)
	case "keys", "KEYS":
		response.buf = keysProcess(request.cmdlist)
	default:
		response.buf = errorProcess("Unknow Cmd :" + request.cmdlist[0])
	}

	return response
}
