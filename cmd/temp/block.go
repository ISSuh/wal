package main

import (
	"bytes"
)

const (
	HeaderSize = 16
)

type BlockHeader struct {
	id       uint64
	blockLen uint32
	crc      uint32
}

type Block struct {
	header BlockHeader
	data   bytes.Buffer
}
