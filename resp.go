package main

import (
	"bufio"
	"fmt"
	"strconv"
)

const (
	STRING  = '+'
	INTEGER = ':'
	ARRAY   = '*'
	BULK    = '$'
)

var Data map[string][]Value

type Value struct {
	typ     string
	str     string
	integer int
	array   []Value
	bulk    string
}

type Resp struct {
	reader *bufio.Reader
}

func NewReader(r *bufio.Reader) *Resp {
	return &Resp{
		reader: r,
	}
}

func (r *Resp) Read() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	}

	return Value{}, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "array"

	size, err := r.readInteger()
	if err != nil {
		return Value{}, nil
	}

	bulk := make([]byte, size)
	r.reader.Read(bulk)

	fmt.Println("size:", size)
	fmt.Println("bulk:", string(bulk))

	v.bulk = string(bulk)

	r.readLine()

	return v, nil
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		if b == '\r' {
			_, err := r.reader.ReadByte()
			if err != nil {
				return nil, 0, err
			}
			break
		}

		n += 1
		line = append(line, b)
	}
	return line, n, nil
}

func (r *Resp) readInteger() (int64, error) {
	line, _, err := r.readLine()
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}

	return i64, nil
}

func (r *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	size, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	for range size {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		val.array = append(val.array, v)
	}

	return val, nil
}

func (r *Resp) Write() {

}
