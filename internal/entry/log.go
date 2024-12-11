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

const (
	logHeaderByteSize = 8
)

type Log struct {
	Sequence int
	PayLoad  []byte
}

func NewLog(size int, sequence int, payload []byte) Log {
	return Log{
		Sequence: sequence,
		PayLoad:  payload,
	}
}

func EncodeLog(log Log) []byte {
	buf := make([]byte, logHeaderByteSize)
	binary.BigEndian.PutUint32(buf, uint32(log.Sequence))
	buf = append(buf, log.PayLoad...)
	return buf
}

func DecodeLog(data []byte) (Log, error) {
	if len(data) < logHeaderByteSize {
		return Log{}, nil
	}

	sequence := binary.BigEndian.Uint32(data[0:4])
	payload := data[8:]

	return Log{
		Sequence: int(sequence),
		PayLoad:  payload,
	}, nil
}
