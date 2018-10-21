package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) >= length {
			return str
		}
	}
}

func padLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str
		}
	}
}

func read_int(data []byte) (ret int32) {
	b := bytes.NewBuffer(data)
	binary.Read(b, binary.LittleEndian, &ret)
	return
}

func packMessage(data []byte) []byte {
	dataLength := len(data)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(dataLength))
	return append(bs, data...)
}

func unpackMessage(r *bufio.Reader) ([]byte, error) {
	peek, err := r.Peek(4)
	if err != nil && err != io.EOF {
		return []byte(""), err
	}
	dataLength := read_int(peek)

	_, err = r.Discard(4)
	if err != nil {
		return []byte(""), err
	}

	b := make([]byte, dataLength)
	readLength, err := r.Read(b)
	if err != nil && err != io.EOF {
		return []byte(""), err
	}

	if int32(readLength) != dataLength {
		return []byte(""), err
	}

	return b, nil
}
