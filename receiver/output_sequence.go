package receiver

import (
	"fmt"
	"io"
	"log"

	"github.com/shamhub/threadapp/config"
	"github.com/shamhub/threadapp/data"
	"github.com/shamhub/threadapp/device"
)

type OutputSequence struct {
	sequence       []bool
	printIndexFlag bool
	length         uint64
}

type Receiver struct {
	outputSequence *OutputSequence
	Log            *log.Logger
}

func newReceiverSequence(outputSeqSize uint64) *OutputSequence {
	return &OutputSequence{
		sequence:       make([]bool, outputSeqSize),
		printIndexFlag: false,
		length:         outputSeqSize,
	}
}

func NewReceiver(outputSeqSize uint64) (*Receiver, error) {

	loggingDevice, fileErr := device.NewLogFileDevice(config.GetReceiverLogFileName())
	if fileErr != nil {
		return nil, fileErr
	}

	l := log.New(loggingDevice, "receiver: ", log.LstdFlags)
	return &Receiver{
		outputSequence: newReceiverSequence(outputSeqSize),
		Log:            l,
	}, nil
}

func (outputSequence *OutputSequence) getSize() uint64 {
	return outputSequence.length
}

func (receiver *Receiver) Print(seqNumber uint64, batchSize uint64, outputFile io.Writer) (error, bool) {

	fmt.Fprintf(outputFile, "[ ")
	if seqNumber >= receiver.outputSequence.length {
		receiver.outputSequence.bufferSizeIncrease(seqNumber)
	}
	receiver.outputSequence.sequence[seqNumber] = true

	printedCount := uint64(0) // check for MAX_OBJECTS_TO_PRINT
	var nthBatchStartingIndex uint64
	MaxObjectsToPrint := config.GetMaxPrintSize()
Loop:
	for nthBatchStartingIndex < receiver.outputSequence.length { // check unbroken sequence
		var assessIndex = nthBatchStartingIndex
		for j := assessIndex; j < nthBatchStartingIndex+batchSize; j++ { // Assess nth batch
			if j >= receiver.outputSequence.length { //index out of range - edge case
				break Loop
			}
			if receiver.outputSequence.sequence[j] == false {
				break Loop
			}
		}

		count, printThresholdReached := receiver.printAssessedBatchIndexes(assessIndex, printedCount, batchSize, MaxObjectsToPrint, outputFile)
		if printThresholdReached { // print sequence threshold reached MAX_OBJECTS_TO_PRINT
			fmt.Fprintf(outputFile, " ]  ")
			fmt.Fprintf(outputFile, " ----for input value %d\n", seqNumber)
			return nil, false
		}
		printedCount += count
		if printedCount >= MaxObjectsToPrint { // print sequence threshold reached MAX_OBJECTS_TO_PRINT
			fmt.Fprintf(outputFile, " ]  ")
			fmt.Fprintf(outputFile, " ----for input value %d\n", seqNumber)
			receiver.Log.Printf("****MaxObjectsToPrint threshold(%d) reached \n", MaxObjectsToPrint)
			return nil, false
		}
		nthBatchStartingIndex = assessIndex + batchSize // next batch
	}
	fmt.Fprintf(outputFile, " ]  ")
	fmt.Fprintf(outputFile, " ----for input value %d\n", seqNumber)
	return nil, true
}

func (receiver *Receiver) printAssessedBatchIndexes(startingIndex uint64,
	printedCount uint64,
	batchSize uint64,
	maxObjectsToPrint uint64,
	outputFile io.Writer) (uint64, bool) {
	if printedCount+batchSize > maxObjectsToPrint { // check print size amidst before printing batch of sequence
		receiver.Log.Printf("****MaxObjectsToPrint threshold(%d) reached \n", maxObjectsToPrint)
		return 0, true
	}
	receiver.outputSequence.printIndexFlag = true
	for i := startingIndex; i < startingIndex+batchSize; i++ {
		fmt.Fprintf(outputFile, "%d,", i)
	}
	return batchSize, false
}

func (outputSequence *OutputSequence) bufferSizeIncrease(seqNumber uint64) {
	newBufferSize := 2 * seqNumber // this can be improvised
	tempBuffer := newReceiverSequence(newBufferSize)
	copy(tempBuffer.sequence, outputSequence.sequence)
	*outputSequence = *tempBuffer
	return
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
