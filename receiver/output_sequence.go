package receiver

import (
	"fmt"
	"io"
	"log"

	"github.com/shamhub/threadapp/config"
	"github.com/shamhub/threadapp/data"
	"github.com/shamhub/threadapp/device"
)

type Receiver struct {
	receivedObjects    map[uint64]*data.Object
	Log                *log.Logger
	nextSeqNumExpected uint64
	printedSequences   uint64
	maxObjectsToPrint  uint64
	batchSize          uint64
}

func NewReceiver() (*Receiver, error) {

	loggingDevice, fileErr := device.NewLogFileDevice(config.GetReceiverLogFileName())
	if fileErr != nil {
		return nil, fileErr
	}

	l := log.New(loggingDevice, "receiver: ", log.LstdFlags)
	return &Receiver{
		receivedObjects:    make(map[uint64]*data.Object),
		Log:                l,
		nextSeqNumExpected: 0,
		printedSequences:   0,
		batchSize:          config.GetBatchSize(),
		maxObjectsToPrint:  config.GetMaxPrintSize(),
	}, nil
}

func (receiver *Receiver) Print(nextSequenceIn uint64, object *data.Object, outputFile io.Writer) bool {
	receiver.receivedObjects[nextSequenceIn] = object
	if receiver.nextSeqNumExpected == nextSequenceIn { // set the next sequence expected
		key := nextSequenceIn + 1
		delete(receiver.receivedObjects, nextSequenceIn)
		for {
			if _, ok := receiver.receivedObjects[key]; !ok {
				receiver.nextSeqNumExpected = key
				break
			}
			key++
		}
	}
	fmt.Fprintf(outputFile, "[ ")
	continu := receiver.printBatch(outputFile)
	fmt.Fprintf(outputFile, "]")
	receiver.printedSequences = 0
	fmt.Fprintf(outputFile, "   ----------for input value %d\n", nextSequenceIn)
	return continu
}

func (receiver *Receiver) printBatch(outputFile io.Writer) bool {
	sequenceNumber := uint64(0)
	batchSize := config.GetBatchSize()
	for sequenceNumber+(batchSize-1) < receiver.nextSeqNumExpected { // received unbroken sequences are [0, receiver.nextSeqNumExpected-1]
		if receiver.printedSequences+batchSize > receiver.maxObjectsToPrint {
			receiver.Log.Printf("****Max objects(%d) to print is reached\n", receiver.maxObjectsToPrint)
			return false
		}
		for j := sequenceNumber; j < sequenceNumber+batchSize; j++ {
			fmt.Fprintf(outputFile, "%d, ", j)
			receiver.printedSequences++

		}
		sequenceNumber += batchSize
	}
	return true
}

func (receiver *Receiver) WaitForLastObject(objectCh chan *data.Object, closeCh chan bool, err error, debugCount uint64) {
	select {
	case object := <-objectCh:
		receiver.Log.Printf("****Received one last object: %v\n", object)
	default:
		receiver.Log.Printf("****Received all(%d) objects from sender, without loss\n", debugCount)
	}
	closeCh <- true
	if err != nil {
		receiver.Log.Printf("%v\n", err.Error())
	}

}
