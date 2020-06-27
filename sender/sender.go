package sender

import (
	"fmt"

	"github.com/shamhub/threadapp/data"
)

func LaunchSender(objCh chan *data.Object, closeCh chan bool) {

	// for random sequence numbers
	sendersRandomSequence := data.NewRandomSequence()

	// business related data
	id := "whatever"
	dataByte := []byte{'w', 'h', 'a', 't', 'e', 'v', 'e', 'r'}

	go func() { // Launch a sender on go-routine
		for {
			objCh <- data.CreateObject(id, dataByte, sendersRandomSequence)
			stop := <-closeCh
			if stop == true {
				fmt.Println("Received close signal")
				break
			}
		}
	}()
}
