package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shamhub/threadapp/device"
	"github.com/shamhub/threadapp/sender"

	"github.com/shamhub/threadapp/config"
	"github.com/shamhub/threadapp/data"
	"github.com/shamhub/threadapp/receiver"
)

const windowSize uint64 = 10

// main command package:
//    1) Launch sender
//    2) Behaves as receiver with abstractions provided by "receiver" package
//    3) Handles interaction with user

func main() {

	objectCh := make(chan *data.Object) // Signalling data
	closeCh := make(chan bool)          // manager signalling employee to stop sending objects

	sender, fileErr := sender.NewSender()
	if fileErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open file: ", fileErr.Error())
		os.Exit(1)
	}
	sender.LaunchSender(objectCh, closeCh)

	// Output sequence
	receiver, fileErr := receiver.NewReceiver(windowSize)
	if fileErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open file:  ", fileErr.Error())
		os.Exit(1)
	}

	outputFile, fileErr := device.NewDataFileDevice(config.GetOutputFileName())
	if fileErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open file:  ", fileErr.Error())
		os.Exit(1)
	}

	defer func() {
		if file, ok := sender.Log.Writer().(*os.File); ok {
			file.Sync()
			file.Close()
		} else if handler, ok := sender.Log.Writer().(io.Closer); ok {
			handler.Close()
		}

		if file, ok := receiver.Log.Writer().(*os.File); ok {
			file.Sync()
			file.Close()
		} else if handler, ok := receiver.Log.Writer().(io.Closer); ok {
			handler.Close()
		}

		if file, ok := outputFile.(*device.DataFileDevice); ok {
			file.Sync()
			file.Close()
		}
		fmt.Printf("*****Main exit\n")
	}()

	// Process each object
	debugCount := uint64(0) // for debugging
	for {
		object := <-objectCh // receive object
		receiver.Log.Printf("Received object: %v", object)
		seqNumber := object.GetSequenceNumber() // read seq num of an object
		debugCount++
		fileErr, continu := receiver.Print(seqNumber, config.GetBatchSize(), outputFile) // print sequence
		if fileErr != nil || continu == false {
			receiver.WaitForLastObject(objectCh, closeCh, fileErr, debugCount)
			break
		}
		// if isOkToContinue(debugCount) {
		// 	closeCh <- false
		// } else {
		// 	closeCh <- true
		// 	break
		// }
		closeCh <- false
	}
}

func isOkToContinue(debugCount uint64) bool {

	if debugCount%10 == 0 {

		fmt.Printf("batch size: %d\n", config.GetBatchSize())
		fmt.Printf("maxobjects to print: %d\n", config.GetMaxPrintSize())

		reader := bufio.NewReader(os.Stdin)
		return readInput(reader)
	}
	return true
}

func readInput(reader *bufio.Reader) bool {
	var text string = ""
	for {
		fmt.Print("To continue, say (Yes/No):")
		text, _ = reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1) // for windows CRLF to LF
		if strings.Compare("Yes", text) == 0 ||
			strings.Compare("yes", text) == 0 ||
			strings.Compare("No", text) == 0 ||
			strings.Compare("no", text) == 0 {
			break
		} else {
			fmt.Println("Invalid input")
		}
		fmt.Println(text)
	}
	if strings.Compare("Yes", text) == 0 || strings.Compare("yes", text) == 0 {
		return true
	}
	return false
}
