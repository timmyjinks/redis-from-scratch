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
)

type Value struct {
	typ     string
	str     string
	integer int
	array   []Value
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
	}

	return Value{}, nil
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
	return line, 0, nil
}

func (r *Resp) readInteger() (int64, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, err
	}

	return i64, nil
}

func (r *Resp) readArray() (Value, error) {
	val := Value{}
	size, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	fmt.Println("size", size*2)
	r.reader.ReadByte()
	r.reader.ReadByte()

	for range size * 2 {
		line, _, err := r.readLine()
		fmt.Println("l", string(line))
		if err != nil {
			return Value{}, err
		}
		val.array = append(val.array, Value{str: string(line)})
	}
	fmt.Println("end")

	return Value{}, nil
}
