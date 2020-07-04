package main

import (
	"fmt"
	"io"
	"os"

	"github.com/shamhub/threadapp/device"
	"github.com/shamhub/threadapp/sender"

	"github.com/shamhub/threadapp/config"
	"github.com/shamhub/threadapp/data"
	"github.com/shamhub/threadapp/receiver"
)

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
	receiver, fileErr := receiver.NewReceiver()
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
		nextSequenceIn := object.GetSequenceNumber() // read seq num of an object
		debugCount++
		continu := receiver.Print(nextSequenceIn, object, config.GetBatchSize(), outputFile) // print sequence
		if continu == false {
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
