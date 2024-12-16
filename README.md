# wal

The `wal` (Write-Ahead Logging) project aims to provide an efficient log management system.
This system ensures high reliability by using a write-ahead logging mechanism.

## Features
- Disk based
- Append only write

## Format

### index file

The `Index` struct is encoded into a 20-byte format as follows:

```
+-------------------+--------------------+-------------------+
| Index (8B)        | MetadataOffset (8B)| MetadataSize (4B) |
+-------------------+--------------------+-------------------+
```

### metadata file

The `metadata` struct is encoded into format as follows:

```
// Layout of Data struct:
+----------+-------------+----- ... ----+
| Size (4B) | Index (8B) |  LogMetadata |
+----------+-------------+----- ... ----+

// Layout of LogMetadata struct:
+----------------+-------------+----------------+------------+-------------+
| SegmentID (4B) |   Size (4B) |  Sequence (4B) |  CRC (4B)  | Offset (8B) |
+----------------+-------------+----------------+------------+-------------+
```

### segment file

The `metadata` struct is encoded into format as follows:

```
+----- ... ----+
|   Payload    |
+----- ... ----+
```

## Installation

To install the `wal` package, use the following command:

```sh
go get github.com/ISSuh/wal
```

## Usage

### Creating a New Storage

To create a new storage instance, use the `NewStorage` function:

```go
package main

import (
	"fmt"

	"github.com/ISSuh/wal"
)

func main() {
	options := wal.Options{
		Path:            "/path/to/log/storage",
		SegmentFileSize: 1024 * 1024, // default segment size is 1 GB
		SyncAfterWrite:  true,        // sync file when after wrtie
	}

	storage, err := wal.NewStorage(options)
	if err != nil {
		fmt.Printf("failed to create storage: %v", err)
		return
	}
	defer storage.Close()
}
```

### Writing Data

To write data to the storage, use the `Write` method:

```go
data := []byte("example log data")
index, err := storage.Write(data)
if err != nil {
	log.Fatalf("failed to write data: %v", err)
}
log.Printf("data written at index: %d", index)
```

### Reading Data

To read data from the storage, use the `Read` method:

```go
readData, err := storage.Read(index)
if err != nil {
	log.Fatalf("failed to read data: %v", err)
}
log.Printf("read data: %s", string(readData))
```

### Synchronizing Data

To ensure all data is flushed to disk, use the `Sync` method:

```go
if err := storage.Sync(); err != nil {
	log.Fatalf("failed to sync data: %v", err)
}
```

If **SyncAfterWrite** option is true, don't need to call `Sync` method

### Closing Storage

To properly close the storage and release resources, use the `Close` method:

```go
if err := storage.Close(); err != nil {
	log.Fatalf("failed to close storage: %v", err)
}
```

### Example

```go
package main

import (
    "github.com/ISSuh/wal"
)

func main() {
    options := wal.Options{
        Path: "/path/to/log",
    }

    log, err := wal.NewLog(options)
    if err != nil {
        panic(err)
    }

    defer log.Close()

    // Append a log entry
    entry := []byte("log entry")
    err = log.Append(entry)
    if err != nil {
        panic(err)
    }
}
```

## Benchmark

```sh
BenchmarkWrite-11                          10000           2501022 ns/op             144 B/op          5 allocs/op
BenchmarkWriteWithSyncAfterWrite-11        10000           7084083 ns/op             144 B/op          5 allocs/op
BenchmarkRead-11                           10000              2290 ns/op             160 B/op          6 allocs/op
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.