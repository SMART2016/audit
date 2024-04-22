package wal

import (
	"os"
	"sync"
)

type WAL struct {
	sync.Mutex
	file *os.File
}

func NewWAL(filename string) (*WAL, error) {
	//Open file in append only mode to mimic WAL behavioir.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: file}, nil
}

func (wal *WAL) Write(entry string) error {
	wal.Lock()
	defer wal.Unlock()

	if _, err := wal.file.WriteString(entry + "\n"); err != nil {
		return err
	}
	return wal.file.Sync() // Ensure data is flushed to disk
}

func (wal *WAL) Close() error {
	return wal.file.Close()
}

//func main() {
//	wal, err := NewWAL("wal.log")
//	if err != nil {
//		panic(err)
//	}
//	defer wal.Close()
//
//	if err := wal.Write("Log Entry 1"); err != nil {
//		fmt.Println("Error writing to WAL:", err)
//	}
//	if err := wal.Write("Log Entry 2"); err != nil {
//		fmt.Println("Error writing to WAL:", err)
//	}
//}
