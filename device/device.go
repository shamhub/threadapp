package device

import (
	"fmt"
	"io"
	"os"
)

type LogFileDevice struct {
	fileHandler *os.File
}

func NewLogFileDevice(fileName string) (io.Writer, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	return &LogFileDevice{
		fileHandler: file,
	}, nil
}

func (d *LogFileDevice) Close() error {
	return d.fileHandler.Close()
}

func (d *LogFileDevice) Sync() error {
	return d.fileHandler.Sync()
}

func (d *LogFileDevice) Write(p []byte) (n int, err error) {

	fmt.Fprintf(d.fileHandler, string(p))
	return len(p), nil
}

type DataFileDevice struct {
	fileHandler *os.File
}

func NewDataFileDevice(fileName string) (io.Writer, error) {
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	return &LogFileDevice{
		fileHandler: file,
	}, nil
}

func (d *DataFileDevice) Close() error {
	return d.fileHandler.Close()
}

func (d *DataFileDevice) Sync() error {
	return d.fileHandler.Sync()
}

func (d *DataFileDevice) Write(p []byte) (n int, err error) {

	fmt.Fprintf(d.fileHandler, string(p))
	return len(p), nil
}
