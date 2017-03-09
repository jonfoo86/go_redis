package main

import (
	"fmt"
	"testing"
)

func TestParseCmd(t *testing.T) {

	packa := new(socketPack)
	packb := new(socketPack)
	result, _ := ParseCmd(packa, packb)
	if result == RS_Fail {
		//t.Fatal("parse fail!")
	}

	packa.buf[0] = 1
	packa.buf[1] = 2
	packa.lenght = 2
	packb.buf[0] = 19
	packb.buf[1] = 20
	packb.lenght = 2
	//fmt.Println(packa.buf)
	//fmt.Println(packb.buf)
	result, _ = ParseCmd(packa, packb)
	if result == RS_Fail {
		//t.Fatal("parse fail!")
	}

	//fmt.Print(packa.buf)
}

func TestGetNum(t *testing.T) {
	packb := new(socketPack)

	command := "3\r\n"
	copy(packb.buf[:], command)
	packb.lenght = len(command)
	paramcount, ok, len := GetNum(packb.buf[:], packb.lenght)

	fmt.Println("paramcount, ok, len: ", paramcount, ",", ok, ", ", len)
}

func TestGetStr(t *testing.T) {
	packb := new(socketPack)

	command := "abcd\r\n"
	copy(packb.buf[:], command)
	packb.lenght = len(command)
	str, ok, len := GetStr(packb.buf[:], packb.lenght)

	fmt.Println("str, ok, len: ", str, ",", ok, ", ", len)
}

func TestParseFirstCmd(t *testing.T) {
	packa := new(socketPack)
	packb := new(socketPack)
	result, _ := ParseCmd(packa, packb)
	if result == RS_Fail {
		//t.Fatal("parse fail!")
	}

	packb.buf[0] = 19
	packb.buf[1] = 20
	command := "*3\r\n$3\r\nSET\r\n$5\r\nHENRY\r\n$8\r\nHENRYFAN\r\n"
	copy(packb.buf[:], command)
	packb.lenght = len(command)
	//fmt.Println(packa.buf)
	//fmt.Println(packb.buf)
	result, _ = ParseCmd(packa, packb)
	if result == RS_Fail {
		//t.Fatal("parse fail!")
	}
	fmt.Print(packa.buf)
}

func TestParseConglutinationCmd(t *testing.T) {
	packa := new(socketPack)
	packb := new(socketPack)

	command := "*3\r\n$3\r\nSET\r\n$5\r\nHENRY\r\n$8\r\nHENRYFAN\r\n*2\r\n$3\r\nGet\r\n$5\r\nHENRY\r\n$8\r\n"
	copy(packb.buf[:], command)
	packb.lenght = len(command)
	//fmt.Println(packa.buf)
	//fmt.Println(packb.buf)
	for {
		result, cmdlist := ParseCmd(packa, packb)
		fmt.Println("result1:", result)
		if result != RS_Ok {
			break
		} else {
			fmt.Println("cmdlist:", cmdlist)
		}

		//fmt.Print("after2:", packa)
	}

	command2 := "*3\r\n$3\r\nSET\r\n$5\r\nHENRY\r\n$8\r\nHENRYFAN\r\n*2\r\n$3\r\nGet\r\n$5\r\nHENRY\r"
	copy(packb.buf[:], command2)
	packb.lenght = len(command2)
	packa.lenght = 0
	//fmt.Println(packa.buf)
	//fmt.Println(packb.buf)
	for {
		result, cmdlist := ParseCmd(packa, packb)
		fmt.Println("result2:", result)
		if result != RS_Ok {
			break
		} else {
			fmt.Println("cmdlist:", cmdlist)
		}

	}

}
