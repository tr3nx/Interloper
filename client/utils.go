package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

func read_int(data []byte) (ret int32) {
	b := bytes.NewBuffer(data)
	binary.Read(b, binary.LittleEndian, &ret)
	return
}

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

func packMessage(data []byte) []byte {
	dataLength := len(data)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(dataLength))
	return append(bs, data...)
}

func unpackMessage(c net.Conn) []byte {
	buf := bufio.NewReader(c)
	peek, err := buf.Peek(4)
	if err != nil && err != io.EOF {
		return []byte("")
	}
	dataLength := read_int(peek)

	_, err = buf.Discard(4)
	if err != nil {
		return []byte("")
	}

	newbuf := make([]byte, dataLength)
	readLength, err := buf.Read(newbuf)
	if err != nil && err != io.EOF {
		return []byte("")
	}

	if int32(readLength) != dataLength {
		return []byte("")
	}

	return newbuf
}
