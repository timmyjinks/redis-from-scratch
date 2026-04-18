package main

import (
	"bufio"
	"io"
	"strconv"
	"sync"
)

const (
	STRING  = '+'
	PANIC   = '-'
	INTEGER = ':'
	ARRAY   = '*'
	BULK    = '$'
)

var Data = map[string]string{}
var DataMutex = sync.RWMutex{}

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

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
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
	v.typ = "bulk"

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

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalArray() []byte {
	var bytes = []byte{}
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.array)))...)
	bytes = append(bytes, '\r', '\n')

	for i := range len(v.array) {
		fmt.Println("bulkkk", v.array[i])
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}
func (v Value) marshalBulk() []byte {
	var bytes = []byte{}
	bytes = append(bytes, BULK)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.bulk)))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalString() []byte {
	var bytes = []byte{}
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalError() []byte {
	var bytes = []byte{}
	bytes = append(bytes, PANIC)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
