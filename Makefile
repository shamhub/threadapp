.PHONY: install

install:
	go install github.com/shamhub/threadapp/config
	go install github.com/shamhub/threadapp/data/shuffle
	go install github.com/shamhub/threadapp/data
	go install github.com/shamhub/threadapp/sender
	go install github.com/shamhub/threadapp/receiver
	go install github.com/shamhub/threadapp
