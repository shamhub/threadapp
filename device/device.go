package device

import (
	"fmt"
	"io"
	"os"

	"github.com/shamhub/threadapp/config"
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

type InteractiveDevice struct {
	input   *os.File
	output  *os.File
	yesOrNo []byte
}

func NewInteractiveDevice() (io.ReadWriteCloser, error) {

	return &InteractiveDevice{
		input:   os.Stdin,
		output:  os.Stdout,
		yesOrNo: make([]byte, 4),
	}, nil
}

func (device *InteractiveDevice) Read(p []byte) (n int, err error) {
	return device.input.Read(p)
}

func (device *InteractiveDevice) Write(p []byte) (n int, err error) {

	device.output.Write(p)
	return len(p), nil
}

func (device *InteractiveDevice) Close() error {
	err := device.output.Close()
	if err != nil {
		return err
	}
	return device.input.Close()
}

func (device *InteractiveDevice) IsOkToContinue(debugCount uint64) bool {

	if debugCount%10 == 0 {

		device.output.Write([]byte(fmt.Sprintf("batch size: %d\n", config.GetBatchSize())))
		device.output.Write([]byte(fmt.Sprintf("maxobjects to print: %d\n", config.GetMaxPrintSize())))

		return device.readInput()
	}
	return true
}

func (device *InteractiveDevice) readInput() bool {
	// var text string = ""
	// for {
	// 	device.output.Write([]byte(fmt.Sprintf("To continue, say (Yes/No):")))
	// 	device.input.Read(device.yesOrNo)
	// 	text, _ = reader.ReadString('\n')
	// 	text = strings.Replace(text, "\n", "", -1) // for windows CRLF to LF
	// 	if strings.Compare("Yes", text) == 0 ||
	// 		strings.Compare("yes", text) == 0 ||
	// 		strings.Compare("No", text) == 0 ||
	// 		strings.Compare("no", text) == 0 {
	// 		break
	// 	} else {
	// 		fmt.Println("Invalid input")
	// 	}
	// 	fmt.Println(text)
	// }
	// if strings.Compare("Yes", text) == 0 || strings.Compare("yes", text) == 0 {
	// 	return true
	// }
	return false
}
