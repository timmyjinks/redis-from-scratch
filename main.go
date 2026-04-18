package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		go Read(conn)
	}

}

func Read(conn net.Conn) {
	for {
		buf := bufio.NewReader(conn)
		r := NewReader(buf)

		val, err := r.Read()
		if err != nil {
			log.Println(err)
			return
		}

		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		fmt.Println(command)
		fmt.Println(args)

		action := Handlers[command]

		v := action(args)

		res := v.Marshal()

		conn.Write(res)
	}
}
