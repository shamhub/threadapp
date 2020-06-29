package sender

import (
	"log"

	"github.com/shamhub/threadapp/config"

	"github.com/shamhub/threadapp/data"
	"github.com/shamhub/threadapp/device"
)

type Sender struct {
	Log *log.Logger
}

func NewSender() (*Sender, error) {
	loggingDevice, fileErr := device.NewLogFileDevice(config.GetSenderLogFileName())
	if fileErr != nil {
		return nil, fileErr
	}

	l := log.New(loggingDevice, "sender: ", log.LstdFlags)
	return &Sender{l}, nil
}

func (sender *Sender) LaunchSender(objCh chan *data.Object, closeCh chan bool) {

	// for random sequence numbers
	sendersRandomSequence := data.NewRandomSequence()

	// business related data
	id := "whatever"
	dataByte := []byte{'w', 'h', 'a', 't', 'e', 'v', 'e', 'r'}

	go func() { // Launch a sender on go-routine
		for {
			object := data.CreateObject(id, dataByte, sendersRandomSequence)
			objCh <- object
			sender.Log.Println("Sent object: ", object)
			stop := <-closeCh
			if stop == true {
				sender.Log.Println("Received close signal")
				break
			}
		}
	}()
}
