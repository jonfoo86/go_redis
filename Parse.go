package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type socketPack struct {
	lenght int
	buf    [1024]byte
}

type ResultType int

const (
	RS_Ok ResultType = iota
	RS_Fail
	RS_Error
)

func GetNum(buf []byte, length int) (int, ResultType, int) {
	//fmt.Println(length, "xxxxxx", buf)
	pos := 0
	for ; unicode.IsDigit(rune(buf[pos])) && pos < length; pos++ {
		//fmt.Println(buf[pos], "is digit")
	}
	if pos == 0 {
		return 0, RS_Error, 0
	}
	//len check
	if pos+2 > length {
		return 0, RS_Fail, 0
	}
	if buf[pos] != '\r' || buf[pos+1] != '\n' {
		return 0, RS_Error, 0
	}

	paramnum, err := strconv.Atoi(string(buf[0:pos]))
	if paramnum == 0 {
		fmt.Print(err.Error())
		return 0, RS_Error, 0
	}

	return paramnum, RS_Ok, pos + 2
}

func GetStr(buf []byte, length int) (string, ResultType, int) {
	pos := 0
	for ; unicode.IsPrint(rune(buf[pos])) && pos < length; pos++ {
	}

	if pos == 0 {
		return string(""), RS_Error, 0
	}
	//len check
	if pos+2 > length {
		return string(""), RS_Fail, 0
	}

	if buf[pos] != '\r' || buf[pos+1] != '\n' {
		return string(""), RS_Fail, 0
	}
	return string(buf[0:pos]), RS_Ok, pos + 2
}

func ParseCmd(oldpack *socketPack, newpack *socketPack) (ResultType, []string) {
	var strarray []string
	var pack *socketPack
	if oldpack.lenght > 0 {
		pack = oldpack
		copy(oldpack.buf[oldpack.lenght:], newpack.buf[0:newpack.lenght])
		pack.lenght += newpack.lenght
	} else if newpack.lenght > 0 {
		pack = newpack
	} else {
		return RS_Fail, strarray
	}

	if pack.buf[0] != '*' {
		return RS_Error, strarray
	}
	pos := 1
	paramcount, ok, len := GetNum(pack.buf[pos:], pack.lenght-1)
	if ok != RS_Ok {
		fmt.Println("get paramcount  fail")
		return ok, strarray
	}
	//fmt.Println("paramcount:", paramcount)
	pos += len
	strarray = make([]string, paramcount)
	for i := 0; i < paramcount && pos < pack.lenght; i++ {
		if pack.buf[pos] == '$' {
			pos += 1
		} else {
			return RS_Error, strarray
		}
		if pack.lenght-pos < 3 {
			return RS_Fail, strarray
		}
		count, ok1, len1 := GetNum(pack.buf[pos:], pack.lenght-pos)
		if ok1 != RS_Ok {
			return ok1, strarray
		}
		//fmt.Println(count)
		pos += len1

		if pack.lenght-pos < (2 + count) {
			return RS_Fail, strarray
		}

		str, ok2, len2 := GetStr(pack.buf[pos:], pack.lenght-pos)
		if ok2 != RS_Ok {
			return ok2, strarray
		}
		//fmt.Println("strlen:", count, " parse str:", str)
		strarray[i] = str
		pos += len2
		//fmt.Println("lose buf:" , string( pack.buf[pos:pack.lenght]), "  pos: ",pos, " len:" ,  pack.lenght)
	}

	//fmt.Println("--------lose buf:" , string( pack.buf[pos:pack.lenght]), "  pos: ",pos, " len:" ,  pack.lenght)

	copy(oldpack.buf[0:], pack.buf[pos:pack.lenght])

	oldpack.lenght = pack.lenght - pos
	newpack.lenght = 0

	return RS_Ok, strarray
}
