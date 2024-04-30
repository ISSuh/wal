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

package segment

import (
	"fmt"
	"os"
	"time"
)

const (
	FileName      string = "log"
	FilePermition        = 0644
	FileFlag             = os.O_CREATE | os.O_RDWR | os.O_APPEND
)

func makeFilePath(basePath string, id uint64) string {
	return fmt.Sprintf("%s/%s_%d_%s", basePath, FileName, id, time.Now().String())
}

func openFile(basePath string, id uint64) (*os.File, error) {
	path := makeFilePath(basePath, id)
	file, err := os.OpenFile(path, FileFlag, FilePermition)
	if err != nil {
		return nil, err
	}
	return file, nil
}
