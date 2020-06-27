package receiver

import (
	"fmt"

	"github.com/shamhub/threadapp/config"
)

type OutputSequence struct {
	sequence       []bool
	printIndexFlag bool
	length         uint64
}

func NewReceiverSequence(outputSeqSize uint64) *OutputSequence {
	return &OutputSequence{
		sequence:       make([]bool, outputSeqSize),
		printIndexFlag: false,
		length:         outputSeqSize,
	}
}

func (outputSequence *OutputSequence) GetSize() uint64 {
	return outputSequence.length
}

func (outputSequence *OutputSequence) Print(seqNumber uint64, batchSize uint64) error {
	fmt.Printf("[ ")
	if seqNumber >= outputSequence.length {
		outputSequence.bufferSizeIncrease(seqNumber)
	}
	outputSequence.sequence[seqNumber] = true

	printedCount := uint64(0) // check for MAX_OBJECTS_TO_PRINT
	var nthBatchStartingIndex uint64
	MaxObjectsToPrint := config.GetMaxPrintSize()
Loop:
	for nthBatchStartingIndex < outputSequence.length { // check unbroken sequence
		var assessIndex = nthBatchStartingIndex
		for j := assessIndex; j < nthBatchStartingIndex+batchSize; j++ { // Assess nth batch
			if j >= outputSequence.length { //index out of range - edge case
				break Loop
			}
			if outputSequence.sequence[j] == false {
				break Loop
			}
		}

		count, err := outputSequence.printAssessedBatchIndexes(assessIndex, printedCount, batchSize, MaxObjectsToPrint)
		if err != nil { // print sequence threshold reached MAX_OBJECTS_TO_PRINT
			fmt.Printf(" ]  ")
			fmt.Printf(" ----for input value %d\n", seqNumber)
			return err
		}
		printedCount += count
		if printedCount >= MaxObjectsToPrint { // print sequence threshold reached MAX_OBJECTS_TO_PRINT
			fmt.Printf(" ]  ")
			fmt.Printf(" ----for input value %d\n", seqNumber)
			return fmt.Errorf("****output.Print() - MaxObjectsToPrint threshold(%d) reached \n", MaxObjectsToPrint)
		}
		nthBatchStartingIndex = assessIndex + batchSize // next batch
	}
	fmt.Printf(" ]  ")
	fmt.Printf(" ----for input value %d\n", seqNumber)
	return nil
}

func (outputSequence *OutputSequence) printAssessedBatchIndexes(startingIndex uint64,
	printedCount uint64,
	batchSize uint64,
	MaxObjectsToPrint uint64) (uint64, error) {
	if printedCount+batchSize > MaxObjectsToPrint { // check print size amidst before printing batch of sequence
		return 0, fmt.Errorf("****output.Print() - MaxObjectsToPrint threshold(%d) reached \n", MaxObjectsToPrint)
	}
	outputSequence.printIndexFlag = true
	for i := startingIndex; i < startingIndex+batchSize; i++ {
		fmt.Printf("%d,", i)
	}
	return batchSize, nil
}

func (outputSequence *OutputSequence) bufferSizeIncrease(seqNumber uint64) {
	newBufferSize := 2 * seqNumber // this can be improvised
	tempBuffer := NewReceiverSequence(newBufferSize)
	copy(tempBuffer.sequence, outputSequence.sequence)
	*outputSequence = *tempBuffer
	return
}
