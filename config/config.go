package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicholasjackson/env"
)

var outputFileMutex sync.RWMutex
var outputFile string

var receiverFileMutex sync.RWMutex
var receiverLogFile string

var senderFileMutex sync.RWMutex
var senderLogFile string

var batchMutex sync.RWMutex
var batchSize uint64

var printSizeMutex sync.RWMutex
var maxObjectsToPrint uint64

func init() {
	sizeOfSequenceToPrint := env.Int("MAX_OBJECTS_TO_PRINT", false, 50000, "Outputting the first 50000 objects for any input object")
	sizeOfBatch := env.Int("BATCH_SIZE", false, 100, "Batch size to print in")

	err := env.Parse()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *sizeOfBatch <= 0 {
		fmt.Println("Invalid batch size, input positive value")
		os.Exit(1)
	}

	if *sizeOfSequenceToPrint <= 0 || *sizeOfSequenceToPrint > 50000 {
		fmt.Println("Invalid print size, size range is 1-50000 \n")
		os.Exit(1)
	}

	if *sizeOfBatch > *sizeOfSequenceToPrint {
		fmt.Println("Invalid batch size, batchSize cannot be greater than print size")
		os.Exit(1)
	}
	batchSize = uint64(*sizeOfBatch)
	maxObjectsToPrint = uint64(*sizeOfSequenceToPrint)

	appRoot, _ := os.Getwd()
	logFile := filepath.Join(appRoot, "receiver_log")
	receiverLogFile = logFile

	logFile = filepath.Join(appRoot, "sender_log")
	senderLogFile = logFile

	dataFile := filepath.Join(appRoot, "receiver_output")
	outputFile = dataFile

}

func GetMaxPrintSize() uint64 {
	printSizeMutex.Lock()
	defer printSizeMutex.Unlock()
	return maxObjectsToPrint
}

func GetBatchSize() uint64 {
	batchMutex.Lock()
	defer batchMutex.Unlock()
	return batchSize
}

func GetReceiverLogFileName() string {
	receiverFileMutex.Lock()
	defer receiverFileMutex.Unlock()
	return receiverLogFile
}

func GetSenderLogFileName() string {
	senderFileMutex.Lock()
	defer senderFileMutex.Unlock()
	return senderLogFile
}

func GetOutputFileName() string {
	outputFileMutex.Lock()
	defer outputFileMutex.Unlock()
	return outputFile
}
