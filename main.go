package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

		b, err := buf.Peek(12)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(b))

		val, err := r.Read()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(val)

		conn.Write(fmt.Appendf([]byte(val.str), "+%s", val.str))
	}
}
