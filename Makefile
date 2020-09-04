
all: clean build 

build: 
	glide create --non-interactive 
	glide install 
	go mod init github.com/wfchiang/davic-micro-core 

local: 
	go run main.go 

clean: 
	-rm -rf vendor
	-rm rm go.mod 
	-rm rm go.sum
	-rm glide.lock 
	-rm glide.yaml  