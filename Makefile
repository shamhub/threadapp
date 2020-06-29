.PHONY: install

install:
	go install -race github.com/shamhub/threadapp/device
	go install -race github.com/shamhub/threadapp/config
	go install -race github.com/shamhub/threadapp/data/shuffle
	go install -race github.com/shamhub/threadapp/data
	go install -race github.com/shamhub/threadapp/sender
	go install -race github.com/shamhub/threadapp/receiver
	go install -race github.com/shamhub/threadapp
