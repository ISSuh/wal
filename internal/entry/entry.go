/*
MIT License

Copyright (c) 2024 ISSuh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package entry

import "encoding/binary"

type Log struct {
	CRC      uint32
	Size     uint16
	Sequence uint32
	Type     uint8
	PayLoad  []byte
}

func Encode(log *Log) []byte {
	buf := []byte{}
	buf = binary.LittleEndian.AppendUint32(buf, log.CRC)
	buf = binary.LittleEndian.AppendUint16(buf, log.Size)
	buf = binary.LittleEndian.AppendUint32(buf, log.Sequence)
	buf = append(buf, log.Type)
	buf = append(buf, log.PayLoad...)
	return buf
}

func Decode(buf []byte) (*Log, error) {
	l := &Log{}

	index := 0
	l.CRC = binary.LittleEndian.Uint32(buf)
	index += 4

	l.Size = binary.LittleEndian.Uint16(buf[index:])
	index += 2

	l.Sequence = binary.LittleEndian.Uint32(buf[index:])
	index += 4

	l.Type = buf[index]
	index += 1

	copy(l.PayLoad, buf[index:])
	return l, nil
}
