package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	sender.LaunchSender(objectCh, closeCh)

	// Output sequence
	outputBuffer := receiver.NewReceiverSequence(windowSize)
	var err error

	debugCount := uint64(0) // for debugging

	// Process each object
	for {
		object := <-objectCh                    // receive object
		seqNumber := object.GetSequenceNumber() // read seq num of an object
		debugCount++
		err = outputBuffer.Print(seqNumber, config.GetBatchSize()) // print sequence
		if err != nil {
			waitForLastObject(objectCh, closeCh, err, debugCount)
			break
		}
		if isOkToContinue(debugCount) {
			closeCh <- false
		} else {
			closeCh <- true
			break
		}
	}
	fmt.Printf("*****Main exit\n")
}

func waitForLastObject(objectCh chan *data.Object, closeCh chan bool, err error, debugCount uint64) {
	select {
	case <-objectCh:
		fmt.Printf("****main() - Received one last object(if any)\n")
	default:
		fmt.Printf("****main() - Received all %d objects from sender, without loss\n", debugCount)
	}
	closeCh <- true
	fmt.Fprintf(os.Stderr, "%v\n", err.Error())
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
